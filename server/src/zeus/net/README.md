# net 库

net 库是从光荣使命的 zeus 库中提取出来的网络库。

net 库带服务器库和客户端库，都是 go 语言。
如果是Unity客户端, 则使用客户端框架 GameBox 中的客户端网络库。

客户端连接服务器，并在该连接上双向传输数据流。连接支持 TCP 和 KCP.

服务器和客户端之间以消息的方式通信。

net 库主要用于与Unity客户端的互连，服务器之间推荐使用 gRPC。
因为 Unity gRPC 支持还是实验阶段，所以不能统一使用 gRPC, 未来会统一。
如果服务器固定使用 go 语言，也可选用 go rpc. 
net 库侧重服务器方面的功能，客户端功能仅仅是基本可用。

## 数据流格式

流式结构如下：
```
消息头 | 消息 | 消息头 | 消息 ...
```
消息头，4字节，前三字节表示消息长度，最后一字节为标志位。

标志第0位表示是否压缩，第1位表示是否加密。

### 单条消息组成结构
```
msg_id 消息号 2 byte
msg_body 消息体
```

0 号为无效消息. 

消息体是消息序列化后的二进制串。序列化方式由应用指定。默认采用 protobuf 消息。

## 教程

演示如何用 net 库建立一个简单的服务器和客户端，客户端发送一条消息，服务器回应一条消息。
代码见 demo 子目录。

主要步骤为：
1. 用proto定义消息，然后生成go代码。也可用非pb消息.
1. 为每个消息分配一个ID号，配置成 toml 文件，然后生成消息注册代码和消息处理框架代码。
1. 编写主程序，创建客户端和服务器。
1. 以消息处理框架代码为基础，实现消息处理逻辑。

以下为分步操作.

### 消息定义

用 proto 定义消息，见 [proto/test/hello.proto](demo/proto/test/hello.proto)，
其中有消息 `SayRequest`，`SayResponse`.

```
message SayRequest {
    string data = 1;
}

message SayResponse {
    string data = 1;
    string result = 2;
}
```

虽然proto文件中定义了服务，但仅用到了消息定义，服务仅用于说明。

### 生成消息代码

需要 protoc, 来自 https://github.com/google/protobuf

再安装 protoc-gen-gogofaster

	go get github.com/gogo/protobuf/protoc-gen-gogofaster

生成代码

	protoc --gogofaster_out=. test/hello.proto

生成 `test/hello.pb.go`

### 配置消息号

每个消息要分配一个唯一消息号，以 toml 文件配置。

示例消息号配置于：[msgdef/test.toml](demo/msgdef/test.toml)
```
[ClientToServer]
11001 = "zeus.net.demo.proto.test.SayRequest"

[ServerToClient]
11002 = "zeus.net.demo.proto.test.SayResponse"
```

按消息方向将消息分成2类，然后每条消息配置一个唯一ID号。消息名需要带全部的包名。

消息ID号为 uint16, 0 号为无效消息.

可以按功能分成多个消息号配置文件。

这些配置文件用来自动生成消息注册代码和消息处理框架代码。使用 net/gen 工具生成代码。

### 生成消息注册代码

代码生成工具为 net/gen。

#### 分步执行

进入 `msgdef` 目录，先构建 gen 代码生成工具（须先正确设置GOPATH）：
```
set GOPATH=..\..\..\..\..
go build ../../gen
```

然后在 msgdef/ 目录下执行：
```
gen.exe *.toml
```

msggen 会分 gensvr 和 genclt 子目录创建代码文件。
其中 proc 目录下文件用作消息处理类的框架代码，对其更改会合并到下次生成代码。

#### 批处理执行

可以编写如下类似 bat，先构建 gen.exe, 
然后调用 gen.exe 为所有 toml 文件生成代码，
在当前目录生成代码：

```bat
set GOPATH=...
go build zeus/net/gen
gen.exe a.toml b.toml OTHER_DIR/other.toml 
gofmt.exe -w gensvr
gofmt.exe -w genclt
```

已创建了 [msgdef/generate.bat](demo/msgdef/generate.bat), 可直接运行，一键生成。

### 主程序

#### 服务器

创建 server/main.go, 主函数如下：

```go
	svr, err := gensvr.New("tcp", ":5678", 10000)
	if err != nil {
		panic(err)
	}
	svr.Run()
```

#### 客户端

创建 server/main.go, 主函数如下：
```go
	session, err := genclt.Dial("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}

	session.Start()

	// 发送请求并等待应答
	req := &test.SayRequest{
		Data: "this is a test",
	}
	session.Send(req)
```

### 实现消息处理器

直接更改生成的 proc/Proc_hello.go,
修改其中的 `panic("待实现")` 代码为自己的逻辑代码。
如服务器收到消息会应答，客户端收到消息则打印输出。

```
func (p *Proc_hello) MsgProc_SayRequest(msg *test.SayRequest) {
	fmt.Println("Got test request: ", msg)
}
```

## IMsg 消息类型

代码中用 IMsg 表示所有消息类型的接口：
```go
// IMsg 消息接口，所有的消息类都必须实现的接口
type IMsg interface {
	MarshalTo(data []byte) (n int, err error)
	Unmarshal(data []byte) error
	Size() (n int)
}
```

protobuf 消息自然符合 IMsg 接口。
建议使用protobuf 定义消息，然后生成相应的消息类型。
考虑到客户端使用 protobuf 性能消耗较大，
可使用自定义 MarshalTo() 和 Unmarshal() 的非 pb 消息.

数据流中的消息体就是 IMsg 序列化后的二进制数据串。

protobuf 生成消息见上面示例。

非 pb 消息示例如下，使用一个自定义的 ByteStream 序列化和反序列化消息。
```go
// ClientVertifyReq 验证消息
type ClientVertifyReq struct {
	// UID: 玩家UID或者服务器ID
	UID uint64 //data[1:9]
	// Token: 客户端登录时需要携带Token
	Token string //data[9:41]
}

// MarshalTo 序列化
func (msg *ClientVertifyReq) MarshalTo(data []byte) (n int, err error) {
	bw := common.NewByteStream(data)
	return msg.Size(), bw.Marshal(msg)
}

// Unmarshal 反序列化
func (msg *ClientVertifyReq) Unmarshal(data []byte) error {
	sw := common.NewByteStream(data)
	msg.UID, _ = sw.ReadUInt64()
	msg.Token, _ = sw.ReadStr()
	return nil
}

// Size 获取长度
func (msg *ClientVertifyReq) Size() (n int) {
	return 1 + 8 + len(msg.Token) + 2
}
```

## 服务器会话关闭事件

服务器端功能。

服务器可以接收会话创建和断开事件，进行相应的处理。
会话创建事件在 `NewProc_*()` 函数中处理。会话断开在 
`(*Proc_*).OnClosed()` 函数中处理。这2个函数在生成的proc代码中。

## 会话的用户数据

服务器端功能。

服务器会话提供了接口存取会话的用户数据，用户数据是一个任意的 interface{}.

客户端会话没有提供该功能。

## 设置会话的验证消息ID

服务器端功能。

`(*Server).SetVerifyMsgID(MsgID)` 设置会话的验证消息ID.
强制会话必须验证，会话的第1个消息将做为验证消息，必须是指定消息号。
应用的Proc处理器必须调用 session.SetVerified(), 不然连接将被强制关闭。
