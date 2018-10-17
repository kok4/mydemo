package server

import (
	"fmt"
	"strings"
	"zeus/net/msgen/misc"
)

type SrvWrap struct {
	info string
}

func GetNewSrvWrap() *SrvWrap {
	return &SrvWrap{info: TMPL_Wrap}
}

func (sw *SrvWrap) ReplaceSrvName(srvName string) {
	sw.info = strings.Replace(sw.info, "SERVER_NAME", srvName, -1)
}

func (sw *SrvWrap) ReplaceFunctions(srvName string, srvInfo, cliInfo map[string]string) error {
	// Deal interface
	if err := sw.dealInterface(srvInfo); err != nil {
		return err
	}

	// Deal register
	if err := sw.dealRegister(srvInfo); err != nil {
		return err
	}

	// Add deal function
	if err := sw.addFunctions(srvName, srvInfo); err != nil {
		return err
	}
	return nil
}

func (sw *SrvWrap) dealInterface(srvInfo map[string]string) error {
	var err error
	var replaceStr string

	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		//fmt.Println(key, value)
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("SrvWrap:dealInterface: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\nMsgProc_%s(msg *%s)", ret[1], value)
		return true
	})

	sw.info = strings.Replace(sw.info, "INTERFACE_FUNCTIONS", replaceStr, -1)
	return err
}

func (sw *SrvWrap) dealRegister(srvInfo map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("dealRegister: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\nsess.RegMsgProcFunc(%s, result.MsgProc_%s)", key, ret[1])
		return true
	})

	sw.info = strings.Replace(sw.info, "REGISTER_FUNCTIONS", replaceStr, -1)
	return err
}

func (sw *SrvWrap) addFunctions(srvName string, srvInfo map[string]string) error {
	var err error
	var replaceStr string

	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("SrvWrap:addFunctions: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\n\nfunc (t *t%s_MsgProcWrapper) MsgProc_%s(msg server.IMsg) {"+
			"\nt.msgProc.MsgProc_%s(msg.(*%s))"+
			"\n}", srvName, ret[1], ret[1], value)
		return true
	})

	sw.info += replaceStr
	return err
}

func (sw *SrvWrap) GetContent() string {
	return sw.info
}
