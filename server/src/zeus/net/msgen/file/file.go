package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

// 根据URL获取文件名（带扩展名）
func GetFileNameByPath(path string) string {
	fileName := filepath.Base(path)
	return fileName
}

// 根据URL获取文件名(不带扩展名)
func GetFileMiniNameByPath(path string) string {
	baseName := GetFileNameByPath(path)

	srvName := strings.Split(baseName, ".")
	if len(srvName) < 2 {
		return ""
	}
	return srvName[0]
}

// 生成一个文件(如果文件存在覆盖)
func WriteFile(path, content string) {
	obj := newFileInfo(path)
	obj.AddContent(content)
	obj.Close()
}

// 拷贝文件
func CopyFile(desPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(desPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, ok := io.Copy(dst, src); ok != nil {
		return fmt.Errorf("Err:CopyFile io.copy fail", ok)
	}
	return nil
}
