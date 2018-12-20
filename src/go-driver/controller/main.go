package controller

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"go-driver/common"
	"os"
	"strconv"
)

/*
	1.扫块，读取地址相关的交易
	2.主动发起转币申请，然后链上转币
	3. 1，2 支持合约
	4. todo 运行时修改运行时参数扩展，
*/
//解析后最终秘钥key
var orginKey []byte

func MainEntry() {
	common.IoBr()
	common.IoStartLog("启动钱包...")

	//写入runtimeEnv 用户后续控制台来操作运行时配置
	debug := viper.GetBool(`base.debug`)
	singleton := common.GetInstance()
	singleton.SetEnv(`debug`, debug)
	//解析地址存入单例
	setAddresses()
	//扫块

}

func setAddresses() {

	address, _ := common.IoReadFile(addressFile)
	addresses := bytes.Split(bytes.TrimSpace(address), []byte("\n"))
	addressLen := len(addresses)

	//启动输出密码解密秘钥
	//pwd := inputPwd(1)
	//secretKey := DecodeSecretKey(pwd)

	//debug
	//ioutil.WriteFile("orginKey_debug", secretKey, 0666)
	orginKey, _ := common.IoReadFile("orginKey_debug")

	orginKeys := bytes.Split(bytes.TrimSpace(orginKey), []byte("\n"))
	orginKeyLen := len(orginKeys)

	if addressLen != orginKeyLen {
		fmt.Println(`地址和秘钥不匹配`)
		common.IoStartLog(`地址和秘钥不匹配`)
		os.Exit(1)
	}
	//设置秘钥到运行时
	key := make(map[string][]byte)
	for k, v := range addresses {
		key[string(v)] = orginKeys[k]
	}
	//存入单例
	common.GetInstance().SetAddresses(addresses[0:systemNum], addresses, key)

	common.IoStartLog("成功解析地址数量" + strconv.Itoa(orginKeyLen))

}
