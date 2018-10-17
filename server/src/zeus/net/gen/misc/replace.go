package misc

import (
	"strings"
	"zeus/net/gen/folder"
)

func Replace(tmpl string, cfg *Cfg) string {
	curImportDir := folder.GetCurImportDir()
	tmpl = strings.Replace(tmpl, "${CURRENT_IMPORT_DIR}", curImportDir, -1)
	tmpl = strings.Replace(tmpl, "${S2C_MSG_IMPORTS}", cfg.S2cMsgImports, -1)
	tmpl = strings.Replace(tmpl, "${C2S_MSG_IMPORTS}", cfg.C2sMsgImports, -1)
	return tmpl
}
