package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// 判断文件或文件夹否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// 创建文件夹（递归）
// @param perm 权限 0777（读4写2执行1）
func CreateDir(dirPath string, perms ...os.FileMode) error {
	if dirPath == "" {
		return errors.New("path can not be empty")
	}
	ok := IsExist(dirPath)
	if ok {
		return nil
	}
	var perm os.FileMode
	if len(perms) > 0 {
		perm = perms[0]
	} else {
		perm = 0777
	}
	err := os.MkdirAll(dirPath, perm)
	if err != nil {
		return err
	}
	return nil
}

//写入文本文件内容
// @param force 文件夹不存在时自动创建
func WriteFile(filePath string, body string, forces ...bool) error {
	if len(forces) > 0 && forces[0] {
		dir := strings.Replace(filePath, `\`, "/", -1)
		index := strings.LastIndex(dir, "/")
		dir = dir[:index]
		err := CreateDir(dir)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filePath, []byte(body), 0777)
}
