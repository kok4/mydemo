package msgenc

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"zeus/net/internal/internal/consts"
	"zeus/net/internal/internal/crypt"
	"zeus/net/internal/types"
	"zeus/net/msgreg"

	assert "github.com/aurelien-rainone/assertgo"
	"github.com/golang/snappy"
)

const minCompressSize = 100
const MsgHeadSize = consts.MsgHeadSize

// EncodeMsg 序列化消息.
// 返回的slice带头部长度和消息ID, 但是长度待设置.
func EncodeMsg(msg types.IMsg) ([]byte, error) {
	msgID := msgreg.GetMsgID(msg)
	if msgID == 0 {
		// 应该用 msgreg.RegMsg2ID()注册才行
		return nil, fmt.Errorf("message '%s' is not registered", reflect.TypeOf(msg))
	}

	buf := make([]byte, MsgHeadSize+msg.Size())
	binary.LittleEndian.PutUint16(buf[4:], uint16(msgID))
	n, err := msg.MarshalTo(buf[MsgHeadSize:])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: ", err)
	}
	assert.True(n <= msg.Size())
	return buf[:n+MsgHeadSize], nil
}

// CompressAndEncrypt 压缩和加密已序列化消息.
// 输入消息缓冲区无压缩和加密，带头部长度和消息ID。
func CompressAndEncrypt(buf []byte, forceNoCompress bool, encryptEnabled bool) ([]byte, error) {
	// 压缩后会返回新的buf, 如果没压缩就返回原buf
	msgBuf, err := compress(buf, forceNoCompress)
	if err != nil {
		return nil, err
	}

	if encryptEnabled {
		data := msgBuf[MsgHeadSize:]
		crypt.EncryptData(data)
		msgBuf[3] = msgBuf[3] | 0x2
	}

	// 设置头部长度
	setMsgBufLen(msgBuf)
	return msgBuf, nil
}

func compress(buf []byte, forceNoCompress bool) ([]byte, error) {
	msgSize := len(buf) - MsgHeadSize
	if forceNoCompress || msgSize < minCompressSize {
		return buf, nil
	}

	maxLen := snappy.MaxEncodedLen(msgSize) // 不压缩2字节的ID
	p := make([]byte, MsgHeadSize+maxLen)
	mbuff := snappy.Encode(p[MsgHeadSize:], buf[MsgHeadSize:])

	// p[0..2]长度暂不设置，仅设置p[3]标志位
	p[3] = buf[3] | 0x1 // 压缩标志
	// MsgID 2 字节
	p[4] = buf[4]
	p[5] = buf[5]

	return p[:len(mbuff)+MsgHeadSize], nil
}

func setMsgBufLen(msgBuf []byte) {
	bufLen := len(msgBuf)
	cmdSize := bufLen - 4 // 去除长度和标志共4字节
	msgBuf[0] = byte(cmdSize)
	msgBuf[1] = byte(cmdSize >> 8)
	msgBuf[2] = byte(cmdSize >> 16)
}
