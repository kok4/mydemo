// gen 为消息定义生成消息注册代码和消息处理框架代码.
// 可以输入多个消息定义 toml 文件：
// Usage: gen MsgA.toml MsgB.toml d:\msg\C.toml
// toml 文件示例：
/*
	[ClientToServer]
	# 服务器用来生成 MsgProc
	11000 = "pb.EnterReq"

	[ServerToClient]
	# 客户端用来生成 MsgProc
	11001 = "pb.EnterResp"
	20001 = "my.package.sub.Msg"
*/
// 将在执行目录下生成 gensvr, genclt 目录。
// gensvr/proc/, genclt/proc/ 目录下生成消息处理框架代码，
//   每次生成采用合并的方式进行, 所有更改会保留。
// 其他代码生成都是覆盖方式，手工更改会被清除。

package main

import (
	"fmt"
	"os"
	"strings"
	"zeus/net/gen/file"
	"zeus/net/gen/folder"
	"zeus/net/gen/misc"

	"zeus/net/gen/client"
	"zeus/net/gen/server"

	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// 创建目录
	createFolders()

	var names []string // 输入文件名，不带目录名和扩展名
	for _, confPath := range os.Args[1:] {
		readInConfig(confPath)
		name := getName(confPath)
		cfg := misc.NewCfg(name)
		server.GenerateFiles(cfg)
		client.GenerateFiles(cfg)
		names = append(names, getName(confPath))
	}

	// generated_server.go, generated_client.go 特殊处理
	server.GenServerFile(names)
	client.GenClientFile(names)
}

func usage() {
	fmt.Println(`Usage: gen "MsgA.toml" "MsgB.toml" "d:\msg\C.toml"`)
	os.Exit(1)
}

func readInConfig(conf string) {
	if !file.PathExists(conf) {
		panic(fmt.Sprintf("`%s` does not exist", conf))
	}

	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}
}

func getName(confPath string) string {
	fileName := file.GetFileNameByPath(confPath)
	if len(fileName) == 0 {
		panic(fmt.Sprintf("给定配置文件名为空: `%s`", confPath))
	}

	svcName := strings.Split(fileName, ".")
	if len(svcName) < 2 {
		panic(fmt.Sprintf("参数异常: `%s`", confPath))
	}
	return svcName[0]
}

func createFolders() {
	dirs := []string{
		"gensvr", "gensvr/proc", "gensvr/internal",
		"genclt", "genclt/proc", "genclt/internal",
	}
	for _, dir := range dirs {
		if err := folder.GenDir(dir); err != nil {
			panic(err)
		}
	}
}
