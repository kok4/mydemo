package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type _FileInfo struct {
	bufio.Writer
	*os.File
}

func newFileInfo(path string) *_FileInfo {
	fi := &_FileInfo{}

	var err error
	fi.File, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("open file failed.", err.Error())
		return nil
	}
	return fi
}

func (fi *_FileInfo) AddContent(s string) {
	s += "\n"
	nr := strings.NewReader(s)
	br := bufio.NewReader(nr)
	b := bytes.NewBuffer(make([]byte, 0))
	br.WriteTo(b)
	fi.File.Write(b.Bytes())
}

func (fi *_FileInfo) Close() {
	fi.File.Close()
}
