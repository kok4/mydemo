package misc

import (
	"fmt"
	"strings"
)

// ParseMsgName 分析消息全名，返回消息名和类型名(包名.消息名)
// 消息全名中的包名可以用 '.' 或 '/' 分隔。
// 如："pb.MyMsg" -> "MyMsg", "pb.MyMsg"
//     "a/b/c.Msg2" -> "Msg2", "c.Msg2"
//     "a/b/c/Msg3" -> "Msg3", "c.Msg3"
//     "a.b.c.Msg4" -> "Msg4", "c.Msg4"
func ParseMsgName(fullName string) (msgName, typeName string) {
	replaced := strings.Replace(fullName, ".", "/", -1)
	segs := strings.Split(replaced, "/")
	if len(segs) == 0 {
		panic(fmt.Sprintf("illegal message: '%s'", fullName))
	}

	msgName = segs[len(segs)-1]
	if len(segs) == 1 {
		return msgName, msgName
	}

	pkgName := segs[len(segs)-2]
	typeName = fmt.Sprintf("%s.%s", pkgName, msgName)
	return msgName, typeName
}

// GetMsgImportPath 分析消息全名，返回消息导入路径名.
// 消息全名中的包名可以用 '.' 或 '/' 分隔。
// 如："pb.MyMsg" -> "pb"
//     "a/b/c.Msg2" -> "a/b/c"
//     "a.b.c.Msg4" -> "a/b/c"
func getMsgImportPath(fullName string) string {
	idx := strings.LastIndexAny(fullName, "./")
	if idx == -1 {
		return ""
	}
	path := fullName[:idx]
	return strings.Replace(path, ".", "/", -1)
}
