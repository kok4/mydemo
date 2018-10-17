package types

// ISession 代表一个网络连接
type ISession interface {
	RegMsgProcFunc(msgID MsgID, procFunc func(IMsg))

	Send(IMsg)

	Start()
	Close()
	IsClosed() bool

	RemoteAddr() string

	ResetHb()

	SetVerifyMsgID(verifyMsgID MsgID)
	SetVerified()

	SetOnClosed(func())
}
