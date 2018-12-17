package controller

func MainEntry() {

	pwd := inputPwd(1)
	secretKey := DecodeSecretKey(pwd)
	//debug
	//ioutil.WriteFile("decode", secretKey, 0666)
	//解密后的秘钥存入内存，并且分段一一对应到地址
	secretKey = secretKey

}
