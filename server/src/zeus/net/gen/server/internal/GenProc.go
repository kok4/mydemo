package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/misc"
)

type SvrMsgProc struct {
	info string
}

func NewSvrMsgProc() *SvrMsgProc {
	return &SvrMsgProc{info: kTemplateMsgProc}
}

func (s *SvrMsgProc) ReplaceSvcName(svcName string) {
	s.info = strings.Replace(s.info, "${SERVICE_NAME}", svcName, -1)
}

func (s *SvrMsgProc) ReplaceFunctions(cfg *misc.Cfg) error {
	// Deal RegMsgCreator
	if err := s.dealMsgProc(cfg.Name, cfg.C2sMsgs); err != nil {
		return err
	}
	return nil
}

func (s *SvrMsgProc) GetContent() string {
	return s.info
}

func (s *SvrMsgProc) dealMsgProc(svcName string, C2sMsgs map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(C2sMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\n\nfunc (p *Proc_%s) MsgProc_%s(msg *%s) {"+
			"\npanic(\"待实现\")\n} ", svcName, msg, typeName)
		return true
	})

	s.info += replaceStr
	return err
}
