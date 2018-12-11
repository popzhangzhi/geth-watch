package common

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
	"time"
)

/*
	记录正确日志，追加
	起始日志在cr.yaml配置
*/
func IoStartLog(data string, arg ...string) []byte {
	turnToByte := []byte(data)
	context, err := ioStart(turnToByte, arg...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return context
}

/*
	记录错误日志，追加
*/
func IoStartLogErr(errPosition string, data string, arg ...string) {

	turnToByte := []byte(data)
	msg := append(append(append([]byte("ERR:"), errPosition...), "|"...), turnToByte...)
	_, err := ioStart(msg, arg...)
	if err != nil {
		fmt.Println(err)
		// todo 写入日志出错，如何提示
	}

	//fmt.Println(string(msg))
	//os.Exit(1)

}

func ioStart(data []byte, arg ...string) ([]byte, error) {
	startLog := viper.GetString("path.startLog")
	n, err := IoLog(startLog, data, arg...)
	return n, err
}

/*
 写入指定文件，指定内容，追加
*/
func IoLog(filename string, data []byte, arg ...string) ([]byte, error) {
	later := append(data, "\n"...)
	context := append([]byte(time.Now().Format("2006-01-02 03:04:05 PM")+"| "), later...)

	for k, v := range arg {
		if k == 0 && v == "true" {
			return context, nil
		}

	}

	fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	n, err := fl.Write(context)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return nil, err

}
