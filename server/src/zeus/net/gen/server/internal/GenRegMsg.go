package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/misc"
)

type SrvRegMsg struct {
	info string
}

func NewSrvRegMsg() *SrvRegMsg {
	return &SrvRegMsg{info: kTemplateRegMsg}
}

func (srm *SrvRegMsg) ReplaceSvcName(svcName string) {
	srm.info = strings.Replace(srm.info, "${SERVICE_NAME}", svcName, -1)
}

func (srm *SrvRegMsg) ReplaceFunctions(cfg *misc.Cfg) error {
	// Deal RegMsg2ID
	if err := srm.dealRegMsg(cfg.Name, cfg.S2cMsgs); err != nil {
		return err
	}

	// Deal RegMsgCreator
	if err := srm.dealRegMsgCreator(cfg.C2sMsgs); err != nil {
		return err
	}
	return nil
}

func (srm *SrvRegMsg) GetContent() string {
	return srm.info
}

func (srm *SrvRegMsg) dealRegMsg(svcName string, s2cMsgs map[string]string) error {
	//msgreg.RegMsg2ID(&pb.LoginResp_Lobby{}, 2102)
	var err error
	var replaceStr string
	misc.OrderedForEach(s2cMsgs, func(key, value string) bool {
		_, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nmsgreg.RegMsg2ID(&%s{}, %s)", typeName, key)
		return true
	})
	srm.info = strings.Replace(srm.info, "${REG_MSG2ID_FUNS}", replaceStr, -1)
	return err
}

func (srm *SrvRegMsg) dealRegMsgCreator(C2sMsgs map[string]string) error {
	//svr.RegMsgCreator(2101, func() info.IMsg { return &pb.LoginReq_Lobby{} })
	var err error
	var replaceStr string
	misc.OrderedForEach(C2sMsgs, func(key, value string) bool {
		_, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nsvr.RegMsgCreator(%s, func() server.IMsg { return &%s{} })", key, typeName)
		return true
	})
	srm.info = strings.Replace(srm.info, "${MSG_CREATOR}", replaceStr, -1)

	return err
}
