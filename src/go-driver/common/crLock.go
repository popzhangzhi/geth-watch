package common

import (
	"github.com/spf13/viper"
	"io/ioutil"
)

var crLockFile string

func RecordPid(data []byte) {

	ioLock(data)
}
func ReadPid() (string, error) {
	crLockFile = viper.GetString("path.crLockFile")
	strb, err := ioutil.ReadFile(crLockFile)
	return string(strb), err

}
func ClearPid() {
	ioLock([]byte("0"))
}
func ioLock(data []byte) {
	crLockFile = viper.GetString("path.crLockFile")
	ioutil.WriteFile(crLockFile, data, 0666)
}
