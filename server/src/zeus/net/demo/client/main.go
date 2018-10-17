package main

import (
	"zeus/net/demo/msgdef/genclt"
	"zeus/net/demo/proto/test"
)

func main() {
	session, err := genclt.Dial("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}

	proc := session.Proc_hello
	proc.Done = make(chan bool)

	session.Start()

	// 发送请求并等待应答
	req := &test.SayRequest{
		Data: "this is a test",
	}
	session.Send(req)

	<-proc.Done
}
