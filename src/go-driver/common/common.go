package common

import "os"

/*
判断文件或者目录存在不存在
*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//存在路径
	if err == nil {
		return true, nil
	}
	//不存在路径
	if os.IsNotExist(err) {
		return false, nil
	}
	//未知错误
	return false, err
}
