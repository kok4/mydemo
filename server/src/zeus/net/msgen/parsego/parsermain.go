package parsego

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"zeus/net/msgen/file"
)

/*
设定copy文件的路径
eg xxxx.exe srcPath, desPath
*/

func MergeFunction(srcContent, desPath string) {
	// 文件不存在直接拷贝
	if ok, _ := PathExists(desPath); !ok {
		dst, err := os.OpenFile(desPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
			return
		}
		defer dst.Close()
		if _, ok := dst.Write([]byte(srcContent)); ok != nil {
			fmt.Println("Err:CopyFile io.copy fail", ok)
		}
		return
	}

	// Src file
	npg1 := NewPaseGoByStr(srcContent)
	srcInfo := make(map[string]*FunInfo)
	npg1.GetFunctions(func(retKey, retVal interface{}) bool {
		bi := new(FunInfo)
		rv := retVal.(*FunInfo)
		bi.Lbrace = rv.Lbrace
		bi.Rbrace = rv.Rbrace
		srcInfo[retKey.(string)] = bi
		return true
	})

	//Test code
	//for k, v := range srcInfo {
	//	println(k, v.Lbrace, v.Rbrace)
	//}

	// Des file
	npg2 := NewParseGo(desPath)
	desInfo := make(map[string]*FunInfo)
	npg2.GetFunctions(func(retKey, retVal interface{}) bool {
		bi := new(FunInfo)
		rv := retVal.(*FunInfo)
		bi.Lbrace = rv.Lbrace
		bi.Rbrace = rv.Rbrace
		desInfo[retKey.(string)] = bi
		return true
	})

	// Test code
	//for k, v := range desInfo {
	//	println(k, v.Lbrace, v.Rbrace)
	//}

	// 合并函数
	// sort item
	var sortMsg []int
	for _, v := range srcInfo {
		sortMsg = append(sortMsg, int(v.Lbrace))
	}
	sort.Ints(sortMsg)

	// src file opreatre
	//srcFile, err := os.Open(srcContent)
	//if err != nil {
	//	panic(err)
	//}
	//defer srcFile.Close()

	var lastPos int // 记录目标文件中上个函数的 "}" 的位置
	for _, valSort := range sortMsg {
		for key, val := range srcInfo {
			if valSort == int(val.Lbrace) {

				// Remove duplicate function names.
				if _, ok := desInfo[key]; ok {
					lastPos = desInfo[key].Rbrace + 1
					continue
				}
				fmt.Println("Merge=>fun name: ", key)
				strFun := srcContent[int(val.Lbrace) : int(val.Rbrace)+1]
				fmt.Println("Merge====:", strFun)
				MergeFun(desPath, fmt.Sprintf("\n\n%s", strFun), lastPos)
			}
		}
	}

	// 遍历map
	desInfo = make(map[string]*FunInfo)
	npg3 := NewParseGo(desPath)
	npg3.GetFunctions(func(retKey, retVal interface{}) bool {
		bi := new(FunInfo)
		rv := retVal.(*FunInfo)
		bi.Lbrace = rv.Lbrace
		bi.Rbrace = rv.Rbrace
		desInfo[retKey.(string)] = bi
		return true
	})
	// 删除多余的函数
	// 代码待优化.( 函数已 'MsgProc_' 默认为pb文件产生的)
	for key, val := range desInfo {

		// 函数名中没有 MsgProc_ 打头的一律忽略
		if strings.Contains(key, "MsgProc_") == false {
			continue
		}
		if _, ok := srcInfo[key]; !ok {
			fmt.Println("Delete=>fun name: ", key)
			DelFunc(desPath, int(val.Lbrace), int(val.Rbrace))
		}
	}

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MergeFun(dest, cont string, startPos int) {
	// 代码效率低, 代码待优化
	f, err := os.Open(dest)
	if err != nil {
		panic(err)
	}

	// 合并代码先
	var writeStr string
	if contents, err := ioutil.ReadAll(f); err == nil {

		preStr := []byte(contents)[:startPos]
		aftStr := []byte(contents)[startPos:]
		writeStr = fmt.Sprintf("%s%s%s", string(preStr), cont, string(aftStr))
	}
	//fmt.Println(writeStr)
	f.Close()

	file.WriteFile(dest, writeStr)
}

func DelFunc(path string, startPos, endPos int) {
	// 代码需优化
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var writeStr string
	if contents, err := ioutil.ReadAll(f); err == nil {
		preStr := []byte(contents)[:startPos]
		aftStr := []byte(contents)[endPos+1:]
		writeStr = fmt.Sprintf("%s%s", preStr, aftStr)
	}
	//fmt.Println(writeStr)
	f.Close()

	file.WriteFile(path, writeStr)
}
