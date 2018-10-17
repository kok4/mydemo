package merge

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"zeus/net/gen/file"
	"zeus/net/gen/merge/internal"
)

type tFuncInfoMap = internal.FuncInfoMap

// MergeFunction 将 srcContent 合并到文件 desPath.
func MergeFunction(srcContent, desPath string) {
	// 文件不存在直接拷贝
	if !file.PathExists(desPath) {
		file.WriteFile(desPath, srcContent)
		return
	}

	mergeNewFuntions(srcContent, desPath)
	deleteOldMsgProc(srcContent, desPath)
}

func parseFromStr(srcContent string) tFuncInfoMap {
	npg1 := internal.NewGoParserFromStr(srcContent)
	return npg1.GetFuncInfoMap()
}

func parseFromFile(desPath string) tFuncInfoMap {
	npg2 := internal.NewGoParserFromFile(desPath)
	return npg2.GetFuncInfoMap()
}

func mergeNewFuntions(srcContent string, desPath string) {
	srcInfo := parseFromStr(srcContent)

	// 按函数名排序
	var srcFuncs []string
	for funName, _ := range srcInfo {
		srcFuncs = append(srcFuncs, funName)
	}
	sort.Strings(srcFuncs)

	desInfo := parseFromFile(desPath)
	for _, funName := range srcFuncs {
		// Remove duplicate function names.
		if _, ok := desInfo[funName]; ok {
			continue
		}

		fmt.Println("Merge function: ", funName)
		info := srcInfo[funName]
		strFun := srcContent[info.FunBeginPos : info.FunEndPos+1]
		appendFun(desPath, strFun)
	}
}

func appendFun(destPath string, funStr string) {
	// 代码效率低, 代码待优化
	f, err := os.Open(destPath)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Failed to read `%s`: %s", destPath, err))
	}
	f.Close()

	writeStr := fmt.Sprintf("%s\n%s", content, funStr)
	file.WriteFile(destPath, writeStr)
}

func deleteOldMsgProc(srcContent string, desPath string) {
	srcInfo := parseFromStr(srcContent)
	desInfo := parseFromFile(desPath)
	for funName, _ := range desInfo {
		// 函数名中没有 MsgProc_ 打头的一律忽略
		if !strings.Contains(funName, "MsgProc_") {
			continue
		}
		if _, ok := srcInfo[funName]; ok {
			continue
		}

		fmt.Println("Delete function: ", funName)
		deleteFunction(desPath, funName)
	}
}

func deleteFunction(desPath string, funcName string) {
	desInfo := parseFromFile(desPath)
	info, ok := desInfo[funcName]
	if !ok {
		return
	}
	deleteContent(desPath, info.FunBeginPos, info.FunEndPos)
}

func deleteContent(path string, startPos, endPos int) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Failed to read `%s`: %s", path, err))
	}
	f.Close()

	preStr := content[:startPos]
	aftStr := content[endPos+1:]
	writeStr := fmt.Sprintf("%s%s", preStr, aftStr)
	file.WriteFile(path, writeStr)
}
