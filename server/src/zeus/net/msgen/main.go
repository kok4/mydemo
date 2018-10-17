// msgen 为消息定义生成消息注册代码和消息处理框架代码.
// 可以输入多个消息定义 toml 文件：
// Usage: msgen MsgA.toml MsgB.toml d:\msg\C.toml
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
// 将在执行目录下生成 generated 目录，并分 server, client 子目录生成代码。

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"zeus/net/msgen/file"
	"zeus/net/msgen/folder"

	"zeus/net/msgen/client"
	"zeus/net/msgen/parsego"
	"zeus/net/msgen/server"

	"github.com/spf13/viper"
)

func main() {
	// 参数检查
	CheckArgment()
	confPaths := os.Args[1:]

	// 创建目录
	if err := CreateFolders(); err != nil {
		panic(err)
	}

	var srvNames []string
	for _, confPath := range confPaths {
		// 配置文件检查
		CheckConf(confPath)

		Do(confPath)
		srvNames = append(srvNames, file.GetFileMiniNameByPath(confPath))
	}

	// GenServer特殊处理
	genSrv := server.GetNewGenServer()
	genSrv.Do(srvNames)
	strRet := genSrv.GetContent()
	rootPath := "./generated"
	srvPath := path.Join(rootPath, "server")
	filePath := fmt.Sprintf("%s/server.go", srvPath)
	file.WriteFile(filePath, strRet)
}

func Do(confPath string) {
	var err error
	var str string
	var filePath string
	srvName := file.GetFileMiniNameByPath(confPath)
	rootPath := "./generated"
	srvPath := path.Join(rootPath, "server")
	cltPath := path.Join(rootPath, "client")

	srvObj := server.GetNewServerInfo(srvName)
	if str, err = srvObj.GenWrap(); err != nil {
		LogErr(err.Error())
	}
	//ServerName_MsgProcWrapper.go
	filePath = fmt.Sprintf("%s/%s_MsgProcWrapper.go", srvPath, srvName)
	file.WriteFile(filePath, str)

	if str, err = srvObj.GenRegMsg(); err != nil {
		LogErr(err.Error())
	}
	filePath = fmt.Sprintf("%s/%s_RegMsg.go", srvPath, srvName)
	file.WriteFile(filePath, str)

	if str, err = srvObj.GenExample(); err != nil {
		LogErr(err.Error())
	}
	filePath = fmt.Sprintf("%s/%s_MsgProc.go.example", srvPath, srvName)
	srvDesPath := viper.GetString("Path.server")
	srvPaths := strings.Split(srvDesPath, ";")
	srvPaths = append(srvPaths, "generated/svrproc")
	for _, value := range srvPaths {
		if len(value) == 0 {
			continue
		}

		genPath := fmt.Sprintf("%s/%s_MsgProc.go", value, srvName)
		parsego.MergeFunction(str, genPath)
	}

	//file.WriteFile(filePath, str)

	cltObj := client.GetNewClientInfo(srvName)
	if str, err = cltObj.GenWrap(); err != nil {
		LogErr(err.Error())
	}
	filePath = fmt.Sprintf("%s/%s_MsgProcWrapper.go", cltPath, srvName)
	file.WriteFile(filePath, str)

	if str, err = cltObj.GenRegMsg(); err != nil {
		LogErr(err.Error())
	}
	filePath = fmt.Sprintf("%s/%s_RegMsg.go", cltPath, srvName)
	file.WriteFile(filePath, str)

	if str, err = cltObj.GenExample(); err != nil {
		LogErr(err.Error())
	}
	cltDesPath := viper.GetString("Path.client")
	cltPaths := strings.Split(cltDesPath, ";")
	cltPaths = append(cltPaths, "generated/cltproc")
	for _, value := range cltPaths {
		if len(value) == 0 {
			continue
		}
		genPath := fmt.Sprintf("%s/%s_MsgProc.go", value, srvName)
		parsego.MergeFunction(str, genPath)
	}
}

func LogErr(format string, a ...interface{}) {
	if len(a) == 0 {
		fmt.Print("Err:" + format)
	} else {
		fmt.Printf("Err:"+format, a)
	}
	os.Exit(0)
}

func CheckArgment() {
	if len(os.Args) < 2 {
		LogErr("参数不正确")
	}
}

func CheckConf(conf string) {
	if ok, err := file.PathExists(conf); ok == false {
		LogErr(err.Error())
	}

	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		LogErr("加载配置文件失败")
	}
}

func CreateFolders() error {
	rootPath := "./generated"
	if err := folder.GenDir(rootPath); err != nil {
		return err
	}
	if err := folder.GenDir(path.Join(rootPath, "server")); err != nil {
		return err
	}
	if err := folder.GenDir(path.Join(rootPath, "client")); err != nil {
		return err
	}
	if err := folder.GenDir(path.Join(rootPath, "svrproc")); err != nil {
		return err
	}
	if err := folder.GenDir(path.Join(rootPath, "cltproc")); err != nil {
		return err
	}
	return nil
}
