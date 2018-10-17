package file

import (
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetFileNameByPath(path string) string {
	fileName := filepath.Base(path)
	//fmt.Println(fileName)
	return fileName
}

func WriteFile(path, content string) {
	obj := newFileInfo(path)
	obj.AddContent(content)
	obj.Close()
}
