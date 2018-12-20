package common

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"time"
)

//默认写入日志的时候带上时间前缀
var prefixTime = true

//Iolog 以字符串的形式返回而非写入文件，用于在协程中执行完协程后在写入。默认不开启
var ioLogReturnStr = false

/*
	记录正确日志，追加
	起始日志在cr.yaml配置
*/
func IoStartLog(data string) []byte {
	turnToByte := []byte(data)
	context, err := ioStart(turnToByte)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return context
}

/*
  返回初始化日志格式的内容
*/
func IoLogFmtStr(data string) []byte {
	ioLogReturnStr = true
	defer closeIoLogReturnStr()
	context := IoStartLog(data)
	return context

}

/*
  关闭返回日志
*/
func closeIoLogReturnStr() {
	ioLogReturnStr = false
}

/*
	记录错误日志，追加
*/
func IoStartLogErr(errPosition string, data string) {

	turnToByte := []byte(data)
	msg := append(append(append([]byte("ERR:"), errPosition...), "|"...), turnToByte...)
	_, err := ioStart(msg)
	if err != nil {
		fmt.Println(err)
		// todo 写入日志出错，如何提示
	}

	//fmt.Println(string(msg))
	//os.Exit(1)

}

/**
写入分割线
*/
func IoBr() {
	prefixTime = false
	defer openPrefixTime()
	ioStart([]byte("-------------------------------------------------"))

}

/*
打开日志格式，时间前缀
*/
func openPrefixTime() {
	prefixTime = true
}

/**
  读取配置中的startLog并写入
*/
func ioStart(data []byte) ([]byte, error) {
	startLog := viper.GetString("path.startLog")
	n, err := IoLog(startLog, data)
	return n, err
}

/*
 写入指定文件，指定内容，追加
*/
func IoLog(filename string, data []byte) ([]byte, error) {
	context := append(data, "\n"...)

	if prefixTime {
		context = append([]byte(time.Now().Format("2006-01-02 03:04:05 PM")+"| "), context...)
	}
	if ioLogReturnStr {
		return context, nil
	}

	return IoFile(filename, context)
}

/*
	写入文件，无格式
*/
func IoFile(filename string, context []byte) ([]byte, error) {

	fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	n, err := fl.Write(context)
	if err == nil && n < len(context) {
		err = io.ErrShortWrite
	}
	return nil, err
}

/*
	读取文件
*/
func IoReadFile(filename string) ([]byte, error) {
	fl, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	return ioutil.ReadAll(fl)
}
