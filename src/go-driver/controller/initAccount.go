package controller

import "go-driver/blockDriver"
import (
	"fmt"
)

var (
	systemNum = 3
	userNum   = 10
)

func init() {}

func Generate() {

	for i := 0; i < systemNum+userNum; i++ {
		addr, priKey, err := blockDriver.DoCreate()
		if err != nil {
			fmt.Println(2, "number "+(string(i)), err)
			break
		}
		fmt.Println(1, addr, priKey)
	}

}
