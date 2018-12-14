package controller

import (
	"bufio"
	"errors"
	"fmt"
	"go-driver/blockDriver"
	"go-driver/common"
	"go-driver/tools"
	"io/ioutil"
	"os"
)

var (
	systemNum     = 3
	addressFile   = "./address.txt"
	secretKeyFile = "./secretKey"
)

func init() {}

func Generate(userNum int) {

	//解密秘钥文件
	rsaOrginData, err := ioutil.ReadFile("test")

	if err != nil {
		fmt.Println(err)
	}
	rsaData, err := tools.RsaDecrypt(rsaOrginData[0:128])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(rsaData))

	os.Exit(1)

	os.Remove(addressFile)
	os.Remove(secretKeyFile)
	//检查文件是否存在，存在不生成
	check, err := checkBeforeenerate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !check {
		fmt.Println(addressFile + " " + secretKeyFile + " 文件存在，不生成秘钥对")
		os.Exit(1)
	}
	//开始生成地址，输入2次密码
	//pwd := inputPwd()
	//fmt.Println(pwd)
	var addresses []byte
	var secretKey []byte

	for i := 0; i < systemNum+userNum; i++ {
		addr, priKey, err := blockDriver.DoCreate()
		if err != nil {
			fmt.Println("number "+string(i), "创建地址错误", err)
			break
		}
		addresses = append(addresses, []byte(addr+"\n")...)
		secretKey = append(secretKey, []byte(priKey+"\n")...)

	}
	//保持成文件，并加密秘钥
	ioutil.WriteFile(addressFile, addresses, 0666)

	ioutil.WriteFile(secretKeyFile, secretKey, 0666)

	fmt.Println(len(secretKey))
	//分段加密 deter = total%17 +100 保证减少加密rsa次数

	status, err := tools.RsaEncryptBigData(secretKey, 100)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(status)

}

func checkBeforeenerate() (bool, error) {
	addrIsset, _ := common.PathExists(addressFile)
	secretkeyIsset, _ := common.PathExists(secretKeyFile)
	if !addrIsset && !secretkeyIsset {
		//两者路径不存在
		return true, nil

	} else if addrIsset && secretkeyIsset {
		//两者路径存在
		return false, nil
	} else {
		//其一路径不存在。返回警告
		err := errors.New(addressFile + " " + secretKeyFile + " 有其一存在，中断生成新的秘钥对")
		return false, err
	}

}

func inputPwd() string {
	fmt.Println("输入生成秘钥对的自定义密码(6位以上)")
	fmt.Print("密码:")
	//todo 密文输入密码，以下方法报错
	//bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(bytePassword))

	scanner := bufio.NewScanner(os.Stdin)
	var pwd string
	pwdTime := "first"

	for scanner.Scan() {
		input := scanner.Text()
		if len(input) < 6 {
			fmt.Println("密码小于6位，请重新输入")
			fmt.Print("密码:")
			continue
		}

		if pwdTime == "first" {
			pwd = input
			pwdTime = "second"
			fmt.Print("再次输入密码:")
		} else {
			if input != pwd {
				fmt.Println("两次密码输入不一致，请重新设置密码（请输入第一次密码）")
				fmt.Print("密码:")
				pwdTime = "first"
				continue
			}
			break
		}

	}

	return pwd
}
