package misc

import (
	"fmt"

	"github.com/spf13/viper"
)

type Cfg struct {
	Name string // 配置文件名

	C2sMsgs map[string]string
	S2cMsgs map[string]string

	C2sMsgImports string
	S2cMsgImports string
}

func NewCfg(name string) *Cfg {
	cfg := &Cfg{
		Name:    name,
		C2sMsgs: viper.GetStringMapString("ClientToServer"),
		S2cMsgs: viper.GetStringMapString("ServerToClient"),
	}
	mapKvMustNotBeEmpty(cfg.C2sMsgs)
	mapKvMustNotBeEmpty(cfg.S2cMsgs)
	cfg.C2sMsgImports = parseImports(cfg.C2sMsgs)
	cfg.S2cMsgImports = parseImports(cfg.S2cMsgs)
	return cfg
}

func parseImports(msgs map[string]string) string {
	// 记录消息的包导入路径，去除重复。
	m := make(map[string]bool)
	for _, msg := range msgs {
		importPath := getMsgImportPath(msg)
		if importPath != "" {
			m["\""+importPath+"\""] = true
		}
	}

	var result string
	for s, _ := range m {
		result += fmt.Sprintf("\t%s\n", s)
	}
	return result
}
