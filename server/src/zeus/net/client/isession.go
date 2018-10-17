package client

import "zeus/net/internal/types"

// XXX 可以删除了。只是有些用户还在用 client.ISession. 应该直接用 *client.Session
type ISession interface {
	RegMsgProcFunc(msgID types.MsgID, msgCreator func() types.IMsg, procFunc func(IMsg))

	Send(IMsg)

	Start()
	Close()
	IsClosed() bool
	SetOnClosed(func())
}
