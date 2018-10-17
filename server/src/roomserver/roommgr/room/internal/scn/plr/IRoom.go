package plr

import (
	"roomserver/types"
	"zeus/net/server"
)

type IRoom interface {
	ID() types.RoomID
	BroadcastMsg(msg server.IMsg)
	BroadcastMsgExcept(msg server.IMsg, uid types.PlayerID)
	Frame() uint32
}
