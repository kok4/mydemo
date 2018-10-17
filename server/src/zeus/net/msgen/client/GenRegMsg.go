package client

import (
	"fmt"
	"strings"
	"zeus/net/msgen/misc"
)

type CltRegMsg struct {
	info string
}

func GetNewCltRegMsg() *CltRegMsg {
	return &CltRegMsg{info: TMPL_RegMsg}
}

func (cm *CltRegMsg) ReplaceSrvName(srvName string) {
	cm.info = strings.Replace(cm.info, "SERVER_NAME", srvName, -1)
}

func (cm *CltRegMsg) ReplaceFunctions(srvName string, srvInfo, cliInfo map[string]string) error {
	// Add deal function
	if err := cm.addFunctions(srvInfo); err != nil {
		return err
	}
	return nil
}

func (cm *CltRegMsg) GetContent() string {
	return cm.info
}

func (cm *CltRegMsg) addFunctions(srvInfo map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(srvInfo, func(key, value interface{}) bool {
		replaceStr += fmt.Sprintf("\nmsgreg.RegMsg2ID(&%s{}, %s)", value, key)
		return true
	})
	cm.info = strings.Replace(cm.info, "REG_MSG", replaceStr, -1)
	return err
}
