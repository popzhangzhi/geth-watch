package common

import (
	"os"
	"time"
)

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

/*
获取当前时间
返回格式Y-m-d H:i:s
*/

func GetDatetime() string {
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	return formatTimeStr

}

func GetUnix() int64 {
	return time.Now().Unix()
}
