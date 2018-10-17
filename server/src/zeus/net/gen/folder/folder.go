package folder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Generate directory if it does not exist.
func GenDir(dir string) error {
	if ok := dirExists(dir); ok == false {
		fmt.Println(fmt.Sprintf("路径不存在....创建(%s)", dir))

		if err := os.Mkdir(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

// Determine if the path exists.
func dirExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err == nil {
		return true
	}

	if os.IsExist(err) {
		return true
	}

	return false
}

// GetCurImportDir 返回当前执行目录对应的import串.
// 如果非空，则以'/'结尾，用于拼接import路径，
// 如 GOPATH/src -> "", GOPATH/src/a/b -> "a/b/"
func GetCurImportDir() string {
	GOPATH := os.Getenv("GOPATH")
	dir, _ := filepath.Abs(fmt.Sprintf("%s/src", GOPATH))
	dir = strings.Replace(dir, "\\", "/", -1) //将\替换成/
	curPath := getCurrentDirectory()          // 不带结尾的 '/'
	relPath := strings.Replace(curPath, dir, "", -1)
	if relPath == "" {
		return ""
	}
	return relPath[1:] + "/" // 如 "a/b/"
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
