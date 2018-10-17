# server

应该使用代码生成后封装的 generated.server.New(), 而不是直接 server.New(),
因为生成的代码中有消息处理器的检查。

SessionMsgProc是消息处理器，其中定义了一些消息处理函数，带 "MsgProc_" 前缀，
MsgProc_ 消息处理函数将会自动注册为对应消息名的处理函数，
如 MsgProc_Test(content) 将会处理 Test 消息。
参数content是个interface{}, 不同消息实际类型不同。

SessionMsgProc 示例：

	type SessionMsgProc struct {
		session ISession
		room    *Room // 游戏房间
		...
	}
	func (s *SessionMsgProc) MsgProc_Test(content interface{}) {
		msg := content.(string)
		fmt.printf("Test msg '%s' from session %d", msg, session.GetID())

		room.ActionChan <- func() {
			room.Test(msg) // 在房间协程中运行
		}
	}

应用一般会要求在房间协程中处理消息，可以通过一个channel传递函数调用到房间协程。
例如Room中有个 ActionChan, 在 Run() 协程中执行动作：

	type Room struct {
		...
		ActionChan chan func() // = make(chan func(), 1024)
	}

	func (r *Room) Run() {
		...
		for {
			select {
			...
			case act := <-r.ActionChan:
				act()
		}
	}
