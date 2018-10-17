package client

import (
	"fmt"
	"strings"
	"zeus/net/msgen/misc"
)

type CltWrap struct {
	info string
}

func GetNewCltWrap() *CltWrap {
	return &CltWrap{info: TMPL_Wrap}
}

func (cw *CltWrap) ReplaceSrvName(srvName string) {
	cw.info = strings.Replace(cw.info, "SERVER_NAME", srvName, -1)
}

func (cw *CltWrap) ReplaceFunctions(srvName string, srvInfo, cliInfo map[string]string) error {
	// Deal interface
	if err := cw.dealInterface(cliInfo); err != nil {
		return err
	}

	// Deal RegMsg
	if err := cw.dealRegMsg(cliInfo); err != nil {
		return err
	}

	// Add deal function
	if err := cw.addFunctions(srvName, cliInfo); err != nil {
		return err
	}
	return nil
}

func (cw *CltWrap) GetContent() string {
	return cw.info
}

func (cw *CltWrap) dealInterface(cltInfo map[string]string) error {
	//sess.RegMsgProcFunc(2102, w.New_LoginResp_Lobby, w.MsgProc_LoginResp_Lobby)
	var err error
	var replaceStr string
	misc.OrderedForEach(cltInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("CltWrap:dealInterface: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\nMsgProc_%s(msg *%s)", ret[1], value)
		return true
	})

	cw.info = strings.Replace(cw.info, "INTERFACE_FUNCTIONS", replaceStr, -1)
	return err
}

func (cw *CltWrap) dealRegMsg(cltInfo map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(cltInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("CltWrap:dealRegMsg: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\nsess.RegMsgProcFunc(%s, w.New_%s, w.MsgProc_%s)", key, ret[1], ret[1])
		return true
	})
	cw.info = strings.Replace(cw.info, "REGMSG_FUNCTIONS", replaceStr, -1)
	return err
}

func (cw *CltWrap) addFunctions(srvName string, cltInfo map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(cltInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("CltWrap:addFunctions: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\n\nfunc (t *T%s_MsgProcWrapper) New_%s() client.IMsg {"+
			"\nreturn &%s{}"+
			"\n}"+
			"\n\nfunc (t *T%s_MsgProcWrapper) MsgProc_%s(msg client.IMsg) {"+
			"\nt.msgProc.MsgProc_%s(msg.(*%s))"+
			"\n}", srvName, ret[1], value, srvName, ret[1], ret[1], value)
		return true
	})
	cw.info += replaceStr
	return err
}
