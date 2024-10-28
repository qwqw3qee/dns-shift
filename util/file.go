package util

import (
	"os"
	"path/filepath"
)

func GetCurDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	return dir
}
