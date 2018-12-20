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

/*
	生成地址路口
*/
func Generate(userNum int) {
	//debug
	//os.Remove(addressFile)
	//os.Remove(secretKeyFile)
	//os.Remove("./orginKey")

	//检查文件是否存在，存在不生成
	check, err := checkBeforeGenerate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !check {
		fmt.Println(addressFile + " " + secretKeyFile + " 文件存在，不生成秘钥对")
		os.Exit(1)
	}
	//开始生成地址，输入2次密码
	pwd := inputPwd(0)
	//生成地址
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
	//保存未加密秘钥文件，正式环境应该注释
	//ioutil.WriteFile(`orginKey`, secretKey, 0666)
	// 进行异或混淆。对一值，异或2次同一值，得到最初值。
	xorSecretKey := tools.OxrSecrectKey(secretKey, pwd)
	//对文件进行rsa加密
	rsaData, err := tools.RsaEncryptBigData(xorSecretKey, 117)
	//保存混淆后的文件
	ioutil.WriteFile(secretKeyFile, rsaData, 0666)
	if err != nil {
		fmt.Println(err)
	}

}

/*
	生成地址对之前的检查
*/
func checkBeforeGenerate() (bool, error) {
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

/*
	0 新建密码，二次确认
	1 输入密码
*/
func inputPwd(flag int) []byte {
	if flag == 0 {
		fmt.Println("输入生成秘钥对的自定义密码(6位以上)")
	}
	fmt.Print("密码:")
	//todo 密文输入密码，以下方法报错 需要考虑win linux mac
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
			if flag == 1 {
				break
			}
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

	return []byte(pwd)
}
