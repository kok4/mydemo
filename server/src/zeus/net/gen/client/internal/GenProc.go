package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/misc"
)

const (
	importPath = "/genclt"
)

type CltMsgProc struct {
	info string
}

func NewCltMsgProc() *CltMsgProc {
	return &CltMsgProc{info: kTemplateMsgProc}
}

func (c *CltMsgProc) ReplaceSvcName(svcName string) {
	c.info = strings.Replace(c.info, "${SERVICE_NAME}", svcName, -1)
}

func (c *CltMsgProc) ReplaceFunctions(cfg *misc.Cfg) error {
	// Add deal function
	if err := c.addFunctions(cfg.Name, cfg.S2cMsgs); err != nil {
		return err
	}
	return nil
}

func (c *CltMsgProc) GetContent() string {
	return c.info
}

func (c *CltMsgProc) addFunctions(svcName string, s2cMsgs map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(s2cMsgs, func(key, value string) bool {
		msg, typeName := misc.ParseMsgName(value)
		replaceStr += fmt.Sprintf("\n\nfunc (p *Proc_%s) MsgProc_%s(msg *%s) {"+
			"\npanic(\"待实现\")"+
			"\n}", svcName, msg, typeName)
		return true
	})
	c.info += replaceStr
	return err
}
