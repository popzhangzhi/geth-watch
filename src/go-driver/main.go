package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

func main() {
	// Sprintf 返回值，并不输出
	//bigint := 16
	//a :=fmt.Sprintf("%#x", bigint)
	//fmt.Println
	//a:="1G9jn5cBMWD6y8CgRtAFFNzu9HmkMywyx4"
	//b:=strings.Count(a,"")-1
	//fmt.Println(b)

	//连接rpc客服端
	client, err := rpc.Dial("http://localhost:8545")

	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}
	//查询是、当前节点拥有的账号
	var account []string
	err = client.Call(&account, "eth_accounts")

	if err != nil {
		fmt.Println("client.Call err", err)
		return
	}
	//查询addr对应的余额
	var balance hexutil.Big
	for _, str := range account {
		err = client.Call(&balance, "eth_getBalance", str, "latest")
		//fmt.Println(balance)
		fmt.Printf("%v:%d\n", str, (*big.Int)(&balance))
	}

	//fmt.Printf("account[0]: %s\nbalance[0]: %s\n", account[0], result)
	//fmt.Printf("accounts: %s\n", account[0])
}
