package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/misc"
)

type CltWrap struct {
	info string
}

func NewCltWrap() *CltWrap {
	return &CltWrap{info: kTemplateWrap}
}

func (cw *CltWrap) ReplaceSvcName(svcName string) {
	cw.info = strings.Replace(cw.info, "${SERVICE_NAME}", svcName, -1)
}

func (cw *CltWrap) ReplaceFunctions(cfg *misc.Cfg) error {
	// Deal interface
	if err := cw.dealInterface(cfg.S2cMsgs); err != nil {
		return err
	}

	// Deal RegMsg
	if err := cw.dealRegMsg(cfg.S2cMsgs); err != nil {
		return err
	}

	// Add deal function
	if err := cw.addFunctions(cfg.Name, cfg.S2cMsgs); err != nil {
		return err
	}
	return nil
}

func (cw *CltWrap) GetContent() string {
	return cw.info
}

func (cw *CltWrap) dealInterface(s2cMsgs map[string]string) error {
	//sess.RegMsgProcFunc(2102, w.New_LoginResp_Lobby, w.MsgProc_LoginResp_Lobby)
	var err error
	var replaceStr string
	misc.OrderedForEach(s2cMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nMsgProc_%s(msg *%s)", msg, typeName)
		return true
	})

	cw.info = strings.Replace(cw.info, "${INTERFACE_FUNCTIONS}", replaceStr, -1)
	return err
}

func (cw *CltWrap) dealRegMsg(s2cMsgs map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(s2cMsgs, func(key, value string) bool {
		msg, _ := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nsess.RegMsgProcFunc(%s, w.New_%s, w.MsgProc_%s)", key, msg, msg)
		return true
	})
	cw.info = strings.Replace(cw.info, "${REGMSG_FUNCTIONS}", replaceStr, -1)
	return err
}

func (cw *CltWrap) addFunctions(svcName string, s2cMsgs map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(s2cMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\n\nfunc (w *Wrapper_%s) New_%s() client.IMsg {"+
			"\nreturn &%s{}"+
			"\n}"+
			"\n\nfunc (w *Wrapper_%s) MsgProc_%s(msg client.IMsg) {"+
			"\nw.proc.MsgProc_%s(msg.(*%s))"+
			"\n}", svcName, msg, typeName, svcName, msg, msg, typeName)
		return true
	})
	cw.info += replaceStr
	return err
}
