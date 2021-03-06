package client

/*
替换字符
	SERVER_NAME
	INTERFACE_FUNCTIONS
	REGMSG_FUNCTIONS
*/
var TMPL_Wrap = `// Code generated by msgen.
package client

import (
	"pb"
	"zeus/net/client"

	assert "github.com/aurelien-rainone/assertgo"
)

type ISERVER_NAME_MsgProc interface {
	INTERFACE_FUNCTIONS
}

// MsgProcWrapper 是 MsgProc 的封装。
// MsgProc 由用户实现，其 MsgProc_MyMsg(msg *pb.MyMsg) 中的参数是具体的消息类。
// 封装后为 MsgProc_MyMsg(msg IMsg), 这个接口才能注册到服务器。
type TSERVER_NAME_MsgProcWrapper struct {
	msgProc ISERVER_NAME_MsgProc
}

func RegMsgProc_SERVER_NAME(sess client.ISession, msgProc ISERVER_NAME_MsgProc) {
	assert.True(sess != nil, "session is nil")
	assert.True(msgProc != nil, "msg proc is nil")
	w := &TSERVER_NAME_MsgProcWrapper{
		msgProc: msgProc,
	}

	// [ServerToClient] 注册接收的消息。需要从ID创建消息。REGMSG_FUNCTIONS
}`
