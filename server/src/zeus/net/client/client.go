// client 包是配合 server 所用的客户端.
/*
用 Dial() 创建一个会话，会话封装了连接。

	session, err := client.Dial("kcp", "1.2.3.4:80")


为了发送和接收消息，需要调用会话的 Start() 方法。Start() 之前还需要注册消息处理器。

	client.RegMsgProc(session, &ClientMsgProc{})
	session.Start()

ClientMsgProc 类是一个消息处理类。可以用代码生成工具生成该类的框架代码。
client.RegMsgProc() 是生成代码中的 client 包提供的。
首先定义消息，编写一个消息定义 toml 文件，如 room.toml:

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

*/
package client

import (
	"fmt"
	"net"

	kcp "github.com/xtaci/kcp-go"
)

// Dial 创建一个连接
func Dial(protocal string, addr string) (*Session, error) {
	var conn net.Conn
	var err error

	if protocal == "tcp" {
		if conn, err = net.Dial(protocal, addr); err != nil {
			return nil, err
		}
	} else if protocal == "kcp" {
		if conn, err = kcp.Dial(addr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unknown network protocol '%s'", protocal)
	}

	sess := NewSession(conn)

	// Todo: 添加心跳
	//	if viper.GetBool("Config.HeartBeat") {
	//		go func() {
	//			hbTickerInterval := time.Duration(viper.GetInt("Config.HBTicker")) * time.Second
	//			hbTicker := time.NewTicker(hbTickerInterval)
	//			for {
	//				select {
	//				case <-hbTicker.C:
	//					if sess.IsClosed() {
	//						hbTicker.Stop()
	//						return
	//					}
	//					sess.Send(&msgdef.HeartBeat{})
	//				}
	//			}
	//		}()
	//	}

	return sess, nil
}
