package conn_handler

import (
	"net"
	"zeus/net/internal/types"
	"zeus/net/server/internal/msg_proc_set"
	"zeus/net/server/internal/session"
	st "zeus/net/server/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
)

type ConnHandler struct {
	msgCreator  types.IMsgCreator // 会话创建时需要
	msgProcSet  *msg_proc_set.MsgProcSet
	sessEvtSink st.ISessEvtSink

	isSessNeedVerify bool        // 是否会话需要验证
	sessVerifyMsgID  types.MsgID // 验证消息ID
}

func New(msgCreator types.IMsgCreator) *ConnHandler {
	assert.True(msgCreator != nil, "msgCreator is nil")
	return &ConnHandler{
		msgCreator: msgCreator,
		msgProcSet: msg_proc_set.New(),
	}
}

// 处理连接.
// 创建 Session, 并且开始会话协程。
func (h *ConnHandler) HandleConn(conn net.Conn) {
	var encrypt bool // Todo: enable encrypt
	sess := session.New(conn, encrypt, h.sessEvtSink, h.msgCreator)

	// 在会话上注册所有处理函数
	if h.msgProcSet != nil {
		h.msgProcSet.RegisterToSession(sess)
	}

	if h.isSessNeedVerify {
		sess.SetVerifyMsgID(h.sessVerifyMsgID)
	}

	sess.ResetHb() // Todo: Move it into Start()
	sess.Start()
}

func (h *ConnHandler) SetSessEvtSink(sink st.ISessEvtSink) {
	h.sessEvtSink = sink
}

func (h *ConnHandler) AddMsgProc(msgProc msg_proc_set.IMsgProc) {
	h.msgProcSet.AddMsgProc(msgProc)
}

func (h *ConnHandler) HasMsgProc() bool {
	return !h.msgProcSet.IsEmpty()
}

// SetVerifyMsg 设置会话的验证消息.
// 强制会话必须验证，会话的第1个消息将做为验证消息，消息类型必须为输入类型.
func (h *ConnHandler) SetVerifyMsgID(verifyMsgID types.MsgID) {
	h.isSessNeedVerify = true
	h.sessVerifyMsgID = verifyMsgID
}
