package server

/*
替换字符串
	SERVER_NAME:
	SERVER_IMPORT
末尾增加函数
*/
var TMPL_Example = `// Code generated by msgen.
// 本文件是对应 SERVER_NAMEMessage 的 MsgProc 实现文件示例。
// 本文件可作为 MsgProc 实现的框架代码。
// 不可以在 generated/server 中实现，因为这样会造成循环导入错误。

package main
import (
	gensvr "SERVER_IMPORT"
	"pb"
)

// SERVER_NAMEMessage_MsgProc 是消息处理类.
// 必须实现 ISERVER_NAMEMessage_MsgProc 接口。
// 名字任意，但建议有 MsgProc 后缀。
type SERVER_NAME_MsgProc struct {
	sess server.ISession // 一般都需要包含session对象

	// 可能还应该包含用户和房间对象
	// user *User
	// room *Room
}

func init() {
	// 设置MsgProc类，这样每个连接就可以Clone()一个MsgProc来处理消息.
	// 必须设置，不然 generated server.New() 会报错。
	gensvr.Set_SERVER_NAME_MsgProc(&SERVER_NAME_MsgProc{})
}

func (m *SERVER_NAME_MsgProc) Clone(sess server.ISession) gensvr.ISERVER_NAME_MsgProc {
	return &SERVER_NAME_MsgProc{
		sess: sess,
		// user, room 暂时为空，待创建
	}
}
`
