package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/folder"
	"zeus/net/gen/misc"
)

type SrvWrap struct {
	info string
}

func NewSrvWrap() *SrvWrap {
	return &SrvWrap{info: kTemplateWrap}
}

func (sw *SrvWrap) ReplaceSvcName(svcName string) {
	sw.info = strings.Replace(sw.info, "${SERVICE_NAME}", svcName, -1)
}

func (sw *SrvWrap) ReplaceFunctions(cfg *misc.Cfg) error {
	// Deal interface
	if err := sw.dealInterface(cfg.C2sMsgs); err != nil {
		return err
	}

	// Deal register
	if err := sw.dealRegister(cfg.C2sMsgs); err != nil {
		return err
	}

	// Add deal function
	if err := sw.addFunctions(cfg.Name, cfg.C2sMsgs); err != nil {
		return err
	}
	return nil
}

func (sw *SrvWrap) dealInterface(C2sMsgs map[string]string) error {
	var err error
	var replaceStr string

	misc.OrderedForEach(C2sMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nMsgProc_%s(msg *%s)", msg, typeName)
		return true
	})

	sw.info = strings.Replace(sw.info, "${INTERFACE_FUNCTIONS}", replaceStr, -1)
	return err
}

func (sw *SrvWrap) dealRegister(C2sMsgs map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(C2sMsgs, func(key, value string) bool {
		msg, _ := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\nsess.RegMsgProcFunc(%s, result.MsgProc_%s)", key, msg)
		return true
	})

	sw.info = strings.Replace(sw.info, "${REG_MSG_PROC_FUNCS}", replaceStr, -1)
	return err
}

func (sw *SrvWrap) addFunctions(svcName string, C2sMsgs map[string]string) error {
	var err error
	var replaceStr string

	misc.OrderedForEach(C2sMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\n\nfunc (w *Wrapper_%s) MsgProc_%s(msg server.IMsg) {"+
			"\nw.proc.MsgProc_%s(msg.(*%s))"+
			"\n}", svcName, msg, msg, typeName)
		return true
	})

	sw.info += replaceStr
	return err
}

func (sw *SrvWrap) GetContent() string {
	curImportDir := folder.GetCurImportDir()
	return strings.Replace(sw.info, "${CURRENT_IMPORT_DIR}", curImportDir, -1)
}
