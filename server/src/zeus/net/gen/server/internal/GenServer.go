package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/folder"
	"zeus/net/gen/misc"
)

type GenServer struct {
	info string
}

func NewGenServer() *GenServer {
	return &GenServer{info: kTemplateServer}
}

func (gs *GenServer) ReplaceFunctions(cfg *misc.Cfg) error {
	return nil
}

func (gs *GenServer) GetContent() string {
	curImportDir := folder.GetCurImportDir()
	return strings.Replace(gs.info, "${CURRENT_IMPORT_DIR}", curImportDir, -1)
}

func (gs *GenServer) Do(srvNames []string) {
	gs.replaceRegMsg(srvNames)
	gs.replaceMsgProc(srvNames)
}

func (gs *GenServer) replaceRegMsg(srvNames []string) {
	var replaceStr string
	for _, svcName := range srvNames {
		replaceStr += fmt.Sprintf("\ninternal.RegMsgCreators_%s(svr)", svcName)
	}
	gs.info = strings.Replace(gs.info, "${REG_MSG_CREATORS}", replaceStr, -1)
}

func (gs *GenServer) replaceMsgProc(srvNames []string) {
	var replaceStr string
	for _, svcName := range srvNames {
		replaceStr += fmt.Sprintf("\nsvr.AddMsgProc(&internal.Wrapper_%s{})", svcName)
	}
	gs.info = strings.Replace(gs.info, "${ADD_MSG_PROC_FUNS}", replaceStr, -1)
}
