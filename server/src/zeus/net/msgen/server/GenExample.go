package server

import (
	"fmt"
	"strings"
	"zeus/net/msgen/folder"
	"zeus/net/msgen/misc"
)

const (
	importPath = "/generated/server"
)

type SrvExample struct {
	info string
}

func GetNewSrvExample() *SrvExample {
	return &SrvExample{info: TMPL_Example}
}

func (se *SrvExample) ReplaceSrvName(srvName string) {
	se.info = strings.Replace(se.info, "SERVER_NAME", srvName, -1)
}

func (se *SrvExample) ReplaceFunctions(srvName string, srvInfo, cltInfo map[string]string) error {
	// Deal SERVER_IMPORT
	if err := se.dealImport(); err != nil {
		return err
	}

	// Deal RegMsgCreator
	if err := se.dealMsgProc(srvName, srvInfo); err != nil {
		return err
	}
	return nil
}

func (se *SrvExample) GetContent() string {
	return se.info
}

func (se *SrvExample) dealImport() error {
	ret := folder.GetRelativePathForGoPath()
	ret += importPath
	se.info = strings.Replace(se.info, "SERVER_IMPORT", ret, -1)
	return nil
}

func (se *SrvExample) dealMsgProc(srvName string, srvInfo map[string]string) error {
	//func (m *LobbyServer_MsgProc) MsgProc_LoginReq_Lobby(msg *pb.LoginReq_Lobby) {
	//	panic("待实现")
	//}

	var err error
	var replaceStr string
	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("SrvExample:addFunctions: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\n\nfunc (m *%s_MsgProc) MsgProc_%s(msg *%s) {"+
			"\npanic(\"待实现\")\n} ", srvName, ret[1], value)
		return true
	})

	se.info += replaceStr
	return err
}
