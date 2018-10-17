package server

import (
	"fmt"
	"strings"
	"zeus/net/msgen/misc"
)

type SrvRegMsg struct {
	info string
}

func GetNewSrvRegMsg() *SrvRegMsg {
	return &SrvRegMsg{info: TMPL_RegMsg}
}

func (srm *SrvRegMsg) ReplaceSrvName(srvName string) {
	srm.info = strings.Replace(srm.info, "SERVER_NAME", srvName, -1)
}

func (srm *SrvRegMsg) ReplaceFunctions(srvName string, srvInfo, cltInfo map[string]string) error {
	// Deal RegMsg2ID
	if err := srm.dealRegMsg(srvName, cltInfo); err != nil {
		return err
	}

	// Deal RegMsgCreator
	if err := srm.dealRegMsgCreator(srvInfo); err != nil {
		return err
	}
	return nil
}

func (srm *SrvRegMsg) GetContent() string {
	return srm.info
}

func (srm *SrvRegMsg) dealRegMsg(srvName string, cltInfo map[string]string) error {
	//msgreg.RegMsg2ID(&pb.LoginResp_Lobby{}, 2102)
	var err error
	var replaceStr string
	misc.OrderedForEach(cltInfo, func(key, value interface{}) bool {
		replaceStr += fmt.Sprintf("\nmsgreg.RegMsg2ID(&%s{}, %s)", value, key)
		return true
	})
	srm.info = strings.Replace(srm.info, "REG_MSG", replaceStr, -1)
	return err
}

func (srm *SrvRegMsg) dealRegMsgCreator(srvInfo map[string]string) error {
	//svr.RegMsgCreator(2101, func() info.IMsg { return &pb.LoginReq_Lobby{} })
	var err error
	var replaceStr string
	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		replaceStr += fmt.Sprintf("\nsvr.RegMsgCreator(%s, func() server.IMsg { return &%s{} })", key, value)
		return true
	})
	srm.info = strings.Replace(srm.info, "MSG_CREATOR", replaceStr, -1)

	return err
}
