package server

import (
	"fmt"
	"zeus/net/gen/file"
	"zeus/net/gen/merge"
	"zeus/net/gen/misc"
	"zeus/net/gen/server/internal"
)

type _ServerInfo struct {
	cfg *misc.Cfg
}

func newServerInfo(cfg *misc.Cfg) *_ServerInfo {
	return &_ServerInfo{
		cfg: cfg,
	}
}

func (srv *_ServerInfo) GenWrap() {
	obj := internal.NewSrvWrap()
	obj.ReplaceSvcName(srv.cfg.Name)
	if err := obj.ReplaceFunctions(srv.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, srv.cfg)
	filePath := fmt.Sprintf("gensvr/internal/Wrapper_%s.go", srv.cfg.Name)
	file.WriteFile(filePath, str)
}

func (srv *_ServerInfo) GenRegMsg() {
	obj := internal.NewSrvRegMsg()
	obj.ReplaceSvcName(srv.cfg.Name)
	if err := obj.ReplaceFunctions(srv.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, srv.cfg)
	filePath := fmt.Sprintf("gensvr/internal/RegMsg_%s.go", srv.cfg.Name)
	file.WriteFile(filePath, str)
}

func (srv *_ServerInfo) GenMsgProc() {
	obj := internal.NewSvrMsgProc()
	obj.ReplaceSvcName(srv.cfg.Name)
	if err := obj.ReplaceFunctions(srv.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, srv.cfg)
	genPath := fmt.Sprintf("gensvr/proc/Proc_%s.go", srv.cfg.Name)
	merge.MergeFunction(str, genPath)
}

func GenServerFile(names []string) {
	genSrv := internal.NewGenServer()
	genSrv.Do(names)
	strRet := genSrv.GetContent()
	file.WriteFile("gensvr/generated_server.go", strRet)
	file.WriteFile("gensvr/internal/IRegMsgCreator.go", internal.IRegMsgCreatorGoContent)
}

func GenerateFiles(cfg *misc.Cfg) {
	srvObj := newServerInfo(cfg)
	srvObj.GenWrap()
	srvObj.GenRegMsg()
	srvObj.GenMsgProc()
}
