package tools

import (
	"crypto/md5"
	"fmt"
)

/*
	对秘钥进行分片xor混淆 16位byte
*/
func OxrSecrectKey(data []byte, pwd []byte) []byte {

	return EncodeXOR(data, md5.Sum(pwd))

}

/*
	分片，进行异或
*/
func EncodeXOR(data []byte, xorBit [16]byte) []byte {

	total := len(data)
	time := CountDividCeil(total, 16)

	fmt.Println(`xor分片`, total, 16, time)

	var xorData []byte
	for i := 0; i < int(time); i++ {
		var tempData []byte
		if i == int(time)-1 {
			tempData = data[i*16:]
		} else {
			tempData = data[i*16 : (i+1)*16]
		}
		tempRsa := byteXor(tempData, xorBit)

		xorData = append(xorData, tempRsa...)

	}

	return xorData

}

/*
	异或算法，不足16位的，只按对应位进行异或
*/
func byteXor(data []byte, xorBit [16]byte) []byte {

	length := len(data)
	rel := make([]byte, length)

	for k, v := range data {
		rel[k] = v ^ xorBit[k]
	}

	return rel

}
