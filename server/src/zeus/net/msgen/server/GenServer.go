package server

import (
	"fmt"
	"strings"
)

type GenServer struct {
	info string
}

func GetNewGenServer() *GenServer {
	return &GenServer{info: TMPL_Server}
}

func (gs *GenServer) ReplaceSrvName(srvName string) {
	//gs.info = strings.Replace(gs.info, "SERVER_NAME", srvName, -1)
}

func (gs *GenServer) ReplaceFunctions(srvNames []string, srvInfo, cltInfo map[string]string) error {
	return nil
}

func (gs *GenServer) GetContent() string {
	return gs.info
}

func (gs *GenServer) Do(srvNames []string) {
	gs.replaceCheckIns(srvNames)
	gs.replaceRegMsg(srvNames)
	gs.replaceMsgProc(srvNames)
}

func (gs *GenServer) replaceCheckIns(srvNames []string) {
	var replaceStr string
	for _, srvName := range srvNames {
		replaceStr += fmt.Sprintf("\nif ins_%s_MsgProc == nil {"+
			"\npanic(\"必须实现 I%s_MsgProc 并且调用 server.Set_%s_MsgProc().\")"+
			"\n}", srvName, srvName, srvName)
	}
	gs.info = strings.Replace(gs.info, "CHECK_INS", replaceStr, -1)
}

func (gs *GenServer) replaceRegMsg(srvNames []string) {
	var replaceStr string
	for _, srvName := range srvNames {
		replaceStr += fmt.Sprintf("\nRegMsgCreators_%s(svr)", srvName)
	}
	gs.info = strings.Replace(gs.info, "REG_MSG", replaceStr, -1)
}

func (gs *GenServer) replaceMsgProc(srvNames []string) {
	var replaceStr string
	for _, srvName := range srvNames {
		replaceStr += fmt.Sprintf("\nsvr.AddMsgProc(&t%s_MsgProcWrapper{})", srvName)
	}
	gs.info = strings.Replace(gs.info, "MSG_PROC", replaceStr, -1)
}
