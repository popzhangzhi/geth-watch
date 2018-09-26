package main

import (
	"fmt"

	"go-driver/controller"
	"time"
)

var (
	coinbase = "0xbaff87a555373dd0358035b77508c41eac84e8c8"
	//发送地址
	from           = "0x558FcdE4d3949880e0Ab240ba24cDd9f2c46aE1c"
	fromPrivateKey = "8e2cdff2c37ae8aad4c0ff102a84f8f0e0a23549a83cc01598d8089ad82e1a15"

	//接收地址
	address    = "0x032bbB648C56daE9370cA4F97D7D9f6019C84B9c"
	privateKey = "51486722177311552563720459288918193559318459571153646758180554017044071229487"

	address2    = "0xD5806F13709D6B6520f5E66a6969e833A0d98C72"
	privateKey2 = "36d6a41017e583ea93be3771d6084a4b96d1eb19d9a347633c6a154d655c7fcf"
)

func main() {
	controller.Generate()
	//blockDriver.DoCreate()

	//blockDriver.SendRowTransaction(from, fromPrivateKey, address2, int64(params.Ether), 0, 0)

	//account := [...]string{coinbase, from, address, address2}

	//for _, str := range account {
	//	a, _ := blockDriver.GetBalance(str)
	//	fmt.Println(str, a)
	//}

	//blockDriver.WatchNewBlock()

	//24到29
	//blockDriver.GetNewBlock()
	//daemon()

}

/**
守护进程
*/

func daemon() {
	timer1 := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer1.C:
			fmt.Println("123213")
		}
	}
}
