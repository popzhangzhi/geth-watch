package controller

import "go-driver/blockDriver"
import (
	"fmt"
)

var (
	systemNum = 3
)

func init() {}

func Generate(userNum int) {

	for i := 0; i < systemNum+userNum; i++ {
		addr, priKey, err := blockDriver.DoCreate()
		if err != nil {
			fmt.Println("number "+(string(i)), err)
			break
		}
		fmt.Println(addr, priKey)
	}

}
