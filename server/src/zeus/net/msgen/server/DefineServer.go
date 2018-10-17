package server

/*
替换字符串
	// SERVER_NAME:
	CHECK_INS
	REG_MSG
	MSG_PROC
末尾增加函数
*/

var TMPL_Server = `// Code generated by msgen.
package server

import (
	"zeus/net/server"
)

type Server = server.Server
type IMsg = server.IMsg
type ISession = server.ISession

// 使用 server.New() 创建服务器可以保证 server.init() 得到执行。
func New(protocal string, addr string, maxConns int) (*Server, error) {
	// 检查MsgProc是否都实现了.CHECK_INS

	svr, err := server.New(protocal, addr, maxConns)
	if err != nil {
		return nil, err
	}

	// 注册消息创建器，用于接收数据后创建消息。REG_MSG

	// 添加MsgProc, 这样新连接创建时会
	// CloneAndRegiterMsgProcFunctions() 注册所有处理函数。MSG_PROC

	return svr, nil
}
`