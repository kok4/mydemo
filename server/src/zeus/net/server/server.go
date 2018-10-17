/*
server 包用于创建一个服务器，监听某个端口，接受连接并处理消息.

消息包格式配合客户端 Unity GameBox 网络组件. 具体见: net/消息定义.md

Server 对象封装了 listen, accept 和 handleConnection 操作。
不对外暴露 Conn 和 Listener 接口，而是暴露了一个 ISession 接口。
ISession 是协程安全的。

ISession 接口可用来向对方发送消息。消息接收和处理是自动的。
消息接收和处理在会话的独立的协程中进行。

代码示例：

	svr, err := server.New("kcp", ":80", 10000)
	go svr.Run()

应该使用代码生成后封装的 generated/server.New(), 而不是直接 server.New(),
因为生成的代码中有消息处理器的检查，并且会自动调用消息注册。

编写一个消息定义 toml 文件，如 room.toml:

	[ClientToServer]
	# 服务器用来生成 MsgProc
	1000 = "pb.EnterReq"

	[ServerToClient]
	# 客户端用来生成 MsgProc
	1001 = "pb.EnterResp"

可以有多个 toml 文件，如下执行代码生成

	generate.exe room.toml test1.toml test2.toml

代码生成将创建 generated 目录, 下分 server, client 目录，其中生成多个 go 文件。
还会生成消息处理器的示例代码，如 Room_MsgProc.go.example, 可在示例代码基础上实现自己的消息处理器。
客户端和服务器应该使用同样的消息定义，生成同样的代码。

可以在 New 之前设置会话事件接收器，用来接收会话创建，关闭等事件。如：

	type SessEvtSink struct {
	}

	func (s *SessEvtSink) OnConnected(sess ISession) {
		fmt.Println("Connected")
	}

	...

	func (s *SessEvtSink) OnClosed(sess ISession) {
		fmt.Println("Closed")
	}

	svr.SetSessEvrSink(&SessEvtSink{})

*/
package server

import (
	"fmt"
	"zeus/net/internal/msgcrtr"
	ntypes "zeus/net/internal/types"
	"zeus/net/server/internal/conn_handler"
	"zeus/net/server/internal/listener"
	"zeus/net/server/internal/msg_proc_set"
	"zeus/net/server/internal/types"

	assert "github.com/aurelien-rainone/assertgo"
)

type Server struct {
	listener    listener.IListener
	connHandler *conn_handler.ConnHandler
	msgCreator  *msgcrtr.MsgCreator
}

// New 创建一个服务服.
// 自动开始监听。
// protocal 支持："kcp", "tcp", "tcp+kcp".
// "tcp+kcp"可接受kcp客户端，也可接受tcp客户端。
// addr 形如：":80", "1.2.3.4:80"
// maxConns 是最大连接数
func New(protocal string, addr string, maxConns int) (*Server, error) {
	msgCreator := msgcrtr.NewMsgCreator()
	srv := &Server{
		connHandler: conn_handler.New(msgCreator),
		msgCreator:  msgCreator,
	}

	var err error
	// 接受新连接时会 go sessMgr.HandleConn(), 运行 session.Start().
	srv.listener, err = listener.NewListener(protocal, addr, maxConns)
	return srv, err
}

// Run 运行服务器.
func (s *Server) Run() {
	if !s.connHandler.HasMsgProc() {
		// 防止直接创建, 忘记 AddMsgProc().
		// 如果是生成的代码，会自动 AddMsgProc().
		panic("No MsgProc! Please use the generated codes to new a server.")
	}
	if s.msgCreator.IsEmpty() {
		// 防止忘记 RegMsgCreator().
		// 如果是生成的代码，会自动 RegMsgCreator().
		panic("Please use the generated codes to RegMsgCreator().")
	}

	s.listener.Run(s.connHandler)
}

// End 结束监听, 结束所有消息处理.
func (s *Server) Close() {
	s.listener.Close()
}

// SetSessEvtSink() 设置一个会话事件接收器.
func (s *Server) SetSessEvtSink(sink types.ISessEvtSink) {
	s.connHandler.SetSessEvtSink(sink)
}

func (s *Server) AddMsgProc(msgProc msg_proc_set.IMsgProc) {
	s.connHandler.AddMsgProc(msgProc)
}

// SetVerifyMsgID 设置会话的验证消息ID.
// 强制会话必须验证，会话的第1个消息将做为验证消息，必须是指定消息号。
// 应用的MsgProc处理器必须调用 session.SetVerified(), 不然连接将被强制关闭。
func (s *Server) SetVerifyMsgID(verifyMsgID ntypes.MsgID) {
	assert.True(s.msgCreator != nil, "msgCreator is nil")
	msg := s.msgCreator.NewMsg(verifyMsgID)
	if msg == nil {
		panic(fmt.Sprintf("unknown verification request message ID %d in (*Server).SetVerifyMsgID(%d)", verifyMsgID, verifyMsgID))
		return
	}

	s.connHandler.SetVerifyMsgID(verifyMsgID)
}

func (s *Server) RegMsgCreator(msgID ntypes.MsgID, msgCreator func() IMsg) {
	s.msgCreator.RegMsgCreator(msgID, msgCreator)
}
