package blockDriver

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"math/big"
)

var Client *ethclient.Client

//ethclient方式连接
func DoEthclientDial() {
	url := viper.GetString(`node.rpcUrl`)
	var err error
	Client, err = ethclient.Dial(url)
	if err != nil {
		fmt.Println("ethclient.Dial err", err)
	}

}

//原生rpc连接
//查询当前节点账号和余额
func DoRpcRowDial() {
	url := viper.GetString(`node.rpcUrl`)
	//连接rpc客服端
	client, err := rpc.Dial(url)

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
