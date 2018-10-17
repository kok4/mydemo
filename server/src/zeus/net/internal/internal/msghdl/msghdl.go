// msghdl 包处理消息.
// 提供统一的消息处理接口：
//    HandleRawMsg(msgID uint16, rawMsgBuf []byte)
// 需要在创建时注册各个 MsgID -> MsgProc函数.
//
package msghdl

import (
	"fmt"
	"zeus/net/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
)

type tMsgID = types.MsgID

// 处理器函数，如 (*MsgProc).MsgProc_Test(msg types.IMsg)
type tMsgProcFunc func(types.IMsg)

// 非线程安全。仅被 session.recvLoop() 协程使用。
type MessageHandler struct {
	mapIDToFunc map[tMsgID]tMsgProcFunc
	msgCreator  types.IMsgCreator
}

func New(msgCreator types.IMsgCreator) *MessageHandler {
	assert.True(msgCreator != nil, "msgCreator is nil")
	return &MessageHandler{
		mapIDToFunc: make(map[tMsgID]tMsgProcFunc),
		msgCreator:  msgCreator,
	}
	// 需要由调用者来 RegMsgProcFunc() 注册所有处理函数。
}

// HandleRawMsg 处理原始消息.
// 输入为去头解密解压后的数据。
// 应该允许不同版本的客户端，所以忽略所有版本不同造成的错误。
func (m *MessageHandler) HandleRawMsg(msgID tMsgID, rawMsgBuf []byte) {
	assert.True(m.msgCreator != nil, "msgCreator is nil")
	var msg types.IMsg
	msg = m.msgCreator.NewMsg(msgID)
	if msg == nil {
		log.Debugf("unknown message ID %d", msgID)
		return
	}
	if err := msg.Unmarshal(rawMsgBuf); err != nil {
		log.Debugf("illegal message: %s", err)
		return
	}

	f, ok := m.mapIDToFunc[msgID]
	if !ok {
		log.Debugf("no handler for msg ID %d", msgID)
		return
	}
	f(msg)
}

// RegMsgProcFunc 注册一个消息处理函数.
// Session 创建时会注册所有的 MsgProc.
func (m *MessageHandler) RegMsgProcFunc(msgID tMsgID, msgProcFunc tMsgProcFunc) {
	assert.True(msgID != 0, "message ID 0 is illegal")
	assert.True(nil != msgProcFunc, "msgProc is nil")
	assert.False(m.isRegistered(msgID),
		fmt.Sprintf("message ID is already registered: %d", msgID))

	m.mapIDToFunc[msgID] = msgProcFunc
}

func (m *MessageHandler) isRegistered(msgID tMsgID) bool {
	_, ok := m.mapIDToFunc[msgID]
	return ok
}
