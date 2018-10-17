package internal

import "zeus/net/internal/types"
import "zeus/net/internal/internal/consts"

type MsgID = types.MsgID

// getMsgID 获取消息ID
func getMsgID(buf []byte) MsgID {
	if len(buf) < consts.MsgIDSize {
		return 0
	}

	return MsgID(buf[0]) | MsgID(buf[1])<<8
}
