package client

import (
	"fmt"
	"strings"
	"zeus/net/msgen/folder"
	"zeus/net/msgen/misc"
)

const (
	importPath = "/generated/client"
)

type CltExample struct {
	info string
}

func GetNewCltExample() *CltExample {
	return &CltExample{info: TMPL_Example}
}

func (ce *CltExample) ReplaceSrvName(srvName string) {
	ce.info = strings.Replace(ce.info, "SERVER_NAME", srvName, -1)
}

func (ce *CltExample) ReplaceFunctions(srvName string, srvInfo, cliInfo map[string]string) error {
	// 这里暂时不处理
	// deal client import
	//if err := ce.dealImport(); err != nil {
	//	return err
	//}

	// Add deal function
	if err := ce.addFunctions(srvName, cliInfo); err != nil {
		return err
	}
	return nil
}

func (ce *CltExample) GetContent() string {
	return ce.info
}

func (ce *CltExample) addFunctions(srvName string, cltInfo map[string]string) error {
	var err error
	var replaceStr string
	misc.OrderedForEach(cltInfo, func(key, value interface{}) bool {
		ret := strings.Split(value.(string), ".")
		if len(ret) != 2 {
			err = fmt.Errorf("CltExample::addFunctions: Split string fail %s", value.(string))
			return false
		}
		replaceStr += fmt.Sprintf("\n\nfunc (m *%s_MsgProc) MsgProc_%s(msg *%s) {"+
			"\npanic(\"待实现\")"+
			"\n}", srvName, ret[1], value)
		return true
	})
	ce.info += replaceStr
	return err
}

func (ce *CltExample) dealImport() error {
	ret := folder.GetRelativePathForGoPath()
	ret += importPath
	ce.info = strings.Replace(ce.info, "CLIENT_IMPORT", ret, -1)
	return nil
}
