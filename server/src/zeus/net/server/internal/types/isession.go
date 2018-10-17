package types

import (
	"zeus/net/internal/types"
)

type ISession interface {
	RegMsgProcFunc(msgID types.MsgID, procFunc func(types.IMsg))
	AddOnClosed(func()) // 注册断开时动作，可注册多个依次调用

	Send(types.IMsg)

	Close()
	IsClosed() bool

	RemoteAddr() string

	GetID() uint64

	GetUserData() interface{}
	SetUserData(data interface{})

	SetVerified()
}
