package client

import (
	"net"
	"zeus/net/internal"
	"zeus/net/internal/msgcrtr"
	"zeus/net/internal/types"
)

// Session 包装内部会话，并提供额外的客户端功能
type Session struct {
	sess       types.ISession
	msgCreator *msgcrtr.MsgCreator
}

func NewSession(conn net.Conn) *Session {
	msgCreator := msgcrtr.NewMsgCreator()
	return &Session{
		sess:       internal.NewSession(conn, false, msgCreator),
		msgCreator: msgCreator,
	}
}

func (s *Session) Send(msg IMsg) {
	s.sess.Send(msg)
}

func (s *Session) Start() {
	if s.msgCreator.IsEmpty() {
		panic("no message defined") // 防止忘记 RegMsgProcFunc()
	}
	s.sess.Start()
}

func (s *Session) Close() {
	s.sess.Close()
}

func (s *Session) IsClosed() bool {
	return s.sess.IsClosed()
}

func (s *Session) SetOnClosed(onClosed func()) {
	s.sess.SetOnClosed(onClosed)
}

// RegMsgProcFunc 注册消息处理函数.
// 3个参数为：消息ID, 消息创建函数，消息处理函数。
// 必须在Start()之前。
func (s *Session) RegMsgProcFunc(msgID types.MsgID, msgCreator func() types.IMsg, procFunc func(IMsg)) {
	s.msgCreator.RegMsgCreator(msgID, msgCreator)
	s.sess.RegMsgProcFunc(msgID, procFunc)
}
