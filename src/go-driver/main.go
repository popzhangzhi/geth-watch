package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strconv"
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
var client *ethclient.Client

func main() {
	//DoCreate()
	//doRpcRowDial()

	doEthclientDial()
	sendRowTransaction(from, fromPrivateKey, address2, int64(params.Ether), 0, 0)

	account := [...]string{coinbase, from, address, address2}

	for _, str := range account {
		a, _ := getBalance(str)
		fmt.Println(str, a)
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
	account = append(account, from)
	account = append(account, address)
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

}

//ethclient方式连接
func doEthclientDial() {
	var err error
	client, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("ethclient.Dial err", err)
	}

}
func getBalance(address string) (float64, error) {
	targetAddress := common.HexToAddress(address)
	ctx := context.Background()
	balance, err := client.BalanceAt(ctx, targetAddress, nil)

	return FromWei(balance), err
}

func FromWei(wei *big.Int) float64 {
	fmt.Println(wei)
	float, err := strconv.ParseFloat(wei.String(), 64)
	if err != nil {
		fmt.Println("转化成float64", err)
	}
	return float / params.Ether
}

//eth.sendTransaction({from:"0xbaff87a555373dd0358035b77508c41eac84e8c8",to:"0x558FcdE4d3949880e0Ab240ba24cDd9f2c46aE1c",value:web3.toWei(50,"ether")})
/*
	输入秘钥，自动签名转账
	from  		   转出地址
	fromPrivateKey 转出地址秘钥 16进制 无"0x"前缀
	to             转入地址
	amount         转账数量 单位wei
	gasLimit       最大消化gas数量 默认不带data转账21000个gas 0为自动获取评估gas数量
	gasPrice       每个gas对应的eth的价格，单位wei。 0默认1*10^11 = 0.0000001eth

*/
func sendRowTransaction(from string, fromPrivateKey string, to string, amount int64, gasLimit uint64, gasPrice int64) {

	//发送地址
	fromAddress := common.HexToAddress(from)
	//接受地址
	toAddress := common.HexToAddress(to)
	ctx := context.Background()
	//适用范围在同一个地址，同一个节点内
	//防止覆盖，相同nonce会覆盖交易，相同nonce，且手续费用大于之前的可以覆盖，否则报错replacement transaction underpriced异常
	//pending值会出现在交易拥挤，手续费不高
	//queued 会出现在nonce值过高，只有值连续不断才能把nonce高值加入pending
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		fmt.Println("nonce", err)
	}
	//最多消化gas的数据，多余退回，默认一次转账不带data为 21000
	if gasLimit == 0 {
		call := ethereum.CallMsg{fromAddress, &toAddress, uint64(0), new(big.Int), big.NewInt(amount), nil}
		estimategas, err := client.EstimateGas(ctx, call)
		if err != nil {
			fmt.Println("estimategas", err)
		}
		gasLimit = estimategas
	}
	//每个gas对应的eth的价格，单位wei 1*10^11 = 0.0000001eth 1e-7*21000
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = client.SuggestGasPrice(ctx)
		if err != nil {
			fmt.Println("suggestGasPrice", err)
			gasPrice = 100000000000
			gasPriceBigInt = big.NewInt(gasPrice)
		}
	} else {
		gasPriceBigInt = big.NewInt(gasPrice)
	}

	Transaction := types.NewTransaction(nonce, toAddress, big.NewInt(amount), gasLimit, gasPriceBigInt, nil)
	//验证秘钥是否合理
	key, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		utils.Fatalf("Failed to load the private key: %v", err)
	}
	Transaction, _ = types.SignTx(Transaction, types.HomesteadSigner{}, key)

	fmt.Println("transaction.data", Transaction.Data())
	fmt.Println("transaction.gas", Transaction.Gas())
	fmt.Println("transaction.gasprice", Transaction.GasPrice())
	fmt.Println("transaction.amount", Transaction.Value())
	fmt.Println("transaction.nonce", Transaction.Nonce())

	err, rel := client.SendTransaction(ctx, Transaction)
	if err != nil {
		fmt.Println("SendTransaction", err)
	}
	fmt.Println("txid:", rel)

}
