package internal

import (
	"fmt"
	"strings"
	"zeus/net/gen/folder"
)

type GenClient struct {
	info  string
	names []string
}

func NewGenClient(names []string) *GenClient {
	result := &GenClient{
		info:  kTemplateClient,
		names: names,
	}
	result.replaceNames()
	return result
}

func (gs *GenClient) GetContent() string {
	curImportDir := folder.GetCurImportDir()
	return strings.Replace(gs.info, "${CURRENT_IMPORT_DIR}", curImportDir, -1)
}

func (gs *GenClient) replaceNames() {
	gs.replaceMsgProcMembers()
	gs.replaceMsgProcIntiMembers()
	gs.replaceRegMsgProcMethods()
}

func (gs *GenClient) replaceMsgProcMembers() {
	template := "\tProc_%s\t*proc.Proc_%s\n"
	var replaceStr string
	for _, name := range gs.names {
		replaceStr += fmt.Sprintf(template, name, name)
	}
	gs.info = strings.Replace(gs.info, "${MSG_PROC_MEMBERS}", replaceStr, -1)
}

func (gs *GenClient) replaceMsgProcIntiMembers() {
	template := "\t\tProc_%s:\t&proc.Proc_%s{},\n"
	var replaceStr string
	for _, name := range gs.names {
		replaceStr += fmt.Sprintf(template, name, name)
	}
	gs.info = strings.Replace(gs.info, "${MSG_PROC_INIT_MEMBERS}", replaceStr, -1)
}

func (gs *GenClient) replaceRegMsgProcMethods() {
	template := "\tinternal.RegProc_%s(s, s.Proc_%s)\n"

	var replaceStr string
	for _, name := range gs.names {
		replaceStr += fmt.Sprintf(template, name, name)
	}
	gs.info = strings.Replace(gs.info, "${REG_MSG_PROC_METHODS}", replaceStr, -1)
}
