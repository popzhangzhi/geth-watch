package blockDriver

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

var Client *ethclient.Client

func init() {
	doEthclientDial()
}

//ethclient方式连接
func doEthclientDial() {
	var err error
	Client, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("ethclient.Dial err", err)
	}

}

//原生rpc连接
//查询当前节点账号和余额
func doRpcRowDial() {

	//连接rpc客服端
	client, err := rpc.Dial("http://localhost:8545")

	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}
	//查询是、当前节点拥有的账号
	var account []string
	err = client.Call(&account, "eth_accounts")
	//account = append(account, from)
	//account = append(account, address)
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
	fmt.Println("------------------------")

}
