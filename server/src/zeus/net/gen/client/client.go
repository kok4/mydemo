package client

import (
	"fmt"
	"zeus/net/gen/client/internal"
	"zeus/net/gen/file"
	"zeus/net/gen/merge"
	"zeus/net/gen/misc"
)

type _ClientInfo struct {
	cfg *misc.Cfg
}

func newClientInfo(cfg *misc.Cfg) *_ClientInfo {
	return &_ClientInfo{
		cfg: cfg,
	}
}

func (ci *_ClientInfo) GenWrap() {
	obj := internal.NewCltWrap()
	obj.ReplaceSvcName(ci.cfg.Name)
	if err := obj.ReplaceFunctions(ci.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, ci.cfg)
	filePath := fmt.Sprintf("genclt/internal/Wrapper_%s.go", ci.cfg.Name)
	file.WriteFile(filePath, str)
}

func (ci *_ClientInfo) GenRegMsg() {
	obj := internal.NewCltRegMsg()
	obj.ReplaceSvcName(ci.cfg.Name)
	if err := obj.ReplaceFunctions(ci.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, ci.cfg)
	filePath := fmt.Sprintf("genclt/internal/RegMsg_%s.go", ci.cfg.Name)
	file.WriteFile(filePath, str)
}

func (ci *_ClientInfo) GenMsgProc() {
	obj := internal.NewCltMsgProc()
	obj.ReplaceSvcName(ci.cfg.Name)
	if err := obj.ReplaceFunctions(ci.cfg); err != nil {
		panic(err)
	}
	str := obj.GetContent()
	str = misc.Replace(str, ci.cfg)
	genPath := fmt.Sprintf("genclt/proc/Proc_%s.go", ci.cfg.Name)
	merge.MergeFunction(str, genPath)
}

func GenerateFiles(cfg *misc.Cfg) {
	cltObj := newClientInfo(cfg)
	cltObj.GenWrap()
	cltObj.GenRegMsg()
	cltObj.GenMsgProc()
}

func GenClientFile(names []string) {
	str := internal.NewGenClient(names).GetContent()
	file.WriteFile("genclt/generated_client.go", str)
	file.WriteFile("genclt/internal/IRegMsgProcFunc.go", internal.IRegMsgProcFuncGoContent)
}
