package controller

import (
	"fmt"
	"go-driver/tools"
	"io/ioutil"
)

/*
	读取加密文件解密成原始秘钥
*/
func DecodeSecretKey(pwd []byte) []byte {

	//解密秘钥文件
	rsaOrginData, err := ioutil.ReadFile(secretKeyFile)

	if err != nil {
		fmt.Println(err)
	}
	rsaData, err := tools.RsaDecryptBigData(rsaOrginData)
	if err != nil {
		fmt.Println(err)
	}
	//
	xorSecretKey := tools.OxrSecrectKey(rsaData, pwd)

	return xorSecretKey

}
