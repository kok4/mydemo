package session

import (
	"net"
	"zeus/net/internal"
	"zeus/net/internal/types"
	st "zeus/net/server/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
	"go.uber.org/atomic"
)

var sessionIDGen atomic.Uint64

// Session 封装客户端服务器通用会话，并提供服务器专用功能。
type Session struct {
	_IInternalSession

	id uint64

	// 可以保存任意用户数据
	userData atomic.Value

	// 会话事件接收器
	evtSink st.ISessEvtSink

	onClosedFuncs []func() // 关闭时依次调用
}

// _IInternalSession 代表一个客户端服务器通用会话.
type _IInternalSession interface {
	RegMsgProcFunc(msgID types.MsgID, procFunc func(types.IMsg))

	Send(types.IMsg)

	Start()
	Close()
	IsClosed() bool

	RemoteAddr() string

	ResetHb()

	SetVerifyMsgID(verifyMsgID types.MsgID)
	SetVerified()

	SetOnClosed(func())
}

func New(conn net.Conn, encryEnabled bool, sessEvtSink st.ISessEvtSink, msgCreator types.IMsgCreator) *Session {
	assert.True(msgCreator != nil, "msgCreator is nil")
	result := &Session{
		_IInternalSession: internal.NewSession(conn, encryEnabled, msgCreator),

		id:      sessionIDGen.Inc(),
		evtSink: sessEvtSink,
	}
	result.SetOnClosed(result.onClosed)

	if sessEvtSink != nil {
		sessEvtSink.OnConnected(result)
	}
	return result
}

func (s *Session) onClosed() {
	for _, f := range s.onClosedFuncs {
		f()
	}

	if s.evtSink == nil {
		return
	}
	s.evtSink.OnClosed(s)
}

// GetID 获取ID
func (s *Session) GetID() uint64 {
	return s.id
}

func (s *Session) GetUserData() interface{} {
	return s.userData.Load()
}

func (s *Session) SetUserData(data interface{}) {
	s.userData.Store(data)
}

func (s *Session) AddOnClosed(f func()) {
	s.onClosedFuncs = append(s.onClosedFuncs, f)
}
