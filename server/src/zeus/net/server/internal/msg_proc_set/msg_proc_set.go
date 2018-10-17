package msg_proc_set

import (
	st "zeus/net/server/internal/types"
)

// 生成的代码中，MsgProcWrapper 实现了 IMsgProc 接口.
// 实际的 MsgProc 用 Wrapper 封装之后就有了这个接口。
type IMsgProc interface {
	// 所有添加的 MsgProc 有个接口，用来在会话创建时，
	// 在会话上注册其所有的处理函数。
	CloneAndRegisterMsgProcFunctions(sess st.ISession)
}

// 消息处理器集合.
// 非线程安全,要求在服务器Run()之前添加所有MsgProc。
type MsgProcSet struct {
	// MsgProc集合：IMsgProc -> true
	msgProcSet map[IMsgProc]bool
}

func New() *MsgProcSet {
	return &MsgProcSet{
		msgProcSet: make(map[IMsgProc]bool),
	}
}

func (m *MsgProcSet) AddMsgProc(msgProc IMsgProc) {
	m.msgProcSet[msgProc] = true
}

func (m *MsgProcSet) RegisterToSession(sess st.ISession) {
	for msgProc, _ := range m.msgProcSet {
		msgProc.CloneAndRegisterMsgProcFunctions(sess)
	}
}

func (m *MsgProcSet) IsEmpty() bool {
	return len(m.msgProcSet) == 0
}
