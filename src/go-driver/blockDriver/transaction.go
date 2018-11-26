package blockDriver

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"strconv"
)

//获取当前金额
func GetBalance(address string) (float64, error) {
	targetAddress := common.HexToAddress(address)
	ctx := context.Background()
	balance, err := Client.BalanceAt(ctx, targetAddress, nil)

	return FromWei(balance), err
}

//转化成ether单位，主要用于友好调试，具体计算不能以这个为准
func FromWei(wei *big.Int) float64 {
	//fmt.Println(wei)
	float, err := strconv.ParseFloat(wei.String(), 64)
	if err != nil {
		fmt.Println("转化成float64", err)
	}
	return float / params.Ether
}

//eth.sendTransaction({from:"0xbaff87a555373dd0358035b77508c41eac84e8c8",to:"0x558FcdE4d3949880e0Ab240ba24cDd9f2c46aE1c",value:web3.toWei(50,"ether")})
/*
	输入秘钥，离线自动签名转账
	from  		   转出地址
	fromPrivateKey 转出地址秘钥 16进制 无"0x"前缀
	to             转入地址
	amount         转账数量 单位wei
	gasLimit       最大消化gas数量 默认不带data转账21000个gas 0为自动获取评估gas数量
	gasPrice       每个gas对应的eth的价格，单位wei。 0默认1*10^11 = 0.0000001eth

*/
func SendRowTransaction(from string, fromPrivateKey string, to string, amount int64, gasLimit uint64, gasPrice int64) string {

	//发送地址
	fromAddress := common.HexToAddress(from)
	//接受地址
	toAddress := common.HexToAddress(to)
	ctx := context.Background()
	//适用范围在同一个地址，同一个节点内
	//防止覆盖，相同nonce会覆盖交易，相同nonce，且手续费用大于之前的可以覆盖，否则报错replacement transaction underpriced异常
	//pending值会出现在交易拥挤，手续费不高
	//queued 会出现在nonce值过高，只有值连续不断才能把nonce高值加入pending
	nonce, err := Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		fmt.Println("nonce", err)
	}
	//最多消化gas的数据，多余退回，默认一次转账不带data为 21000
	if gasLimit == 0 {
		call := ethereum.CallMsg{fromAddress, &toAddress, uint64(0), new(big.Int), big.NewInt(amount), nil}
		estimategas, err := Client.EstimateGas(ctx, call)
		if err != nil {
			fmt.Println("estimategas", err)
		}
		gasLimit = estimategas
	}
	//每个gas对应的eth的价格，单位wei 1*10^11 = 0.0000001eth 1e-7*21000
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = Client.SuggestGasPrice(ctx)
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

	err, rel := Client.SendTransaction(ctx, Transaction)
	if err != nil {
		fmt.Println("SendTransaction", err)
		return ""
	} else {

		fmt.Println("txid:", rel)
		return rel
	}

}

//todo : 交易成功getTransactionReceipt中的status 状态为0x1 0x0为失败

//todo ：获取交易记录和到账通知，考虑使用 eth_newPendingTransactionFilter 来获取所有到达的交易，
//todo ：交易过滤出内部地址，然后入库同时发起到账通知，同时考虑确认数，然后考虑因为确认数从大于2变到0，之前成功的充值需要回滚
//todo : eth_newBlockFilter 确认数 需要扫描块，并且parentHash能串起来，redis？eth_getFilterChanges 接收 eth_getFilterLogs

//todo 每次开启，新增一次filter，包含块，包含交易。每20秒刷新新块，刷新新块匹配链正常（不正常，回滚，通知回滚接口），更新最新块高度。
//todo 每个块来了以后，开启10个线程 处理多个交易，扫描符合的交易记录，

//做监听，每次监听都生成不同的filterId
//注意失效，只有创建后的交易才会记录，该id存在内存中，重启节点后，需要重启开启
func WatchNewBlock() {
	ctx := context.Background()
	err, filterId := Client.WatchNewBlockFilter(ctx)
	if err != nil {
		fmt.Println("watchNewBlock", err)
	}
	fmt.Println(filterId)
	//0x71014d3030a97f16e54e5dced87441a6 0 timeout
	//0x27c81ece4c628f6064d1b705880cf355 1 5.42 - 58 timeout 16

	//0xef7586e6f3db8dc6ec635d7f29e6e3f -- 5.44 -59 timeout 15
	//0xffdcd82615f0bfc817621ca140950132 -- 5.42
	//0x20490d36d2b8091047de433c622ac7bc  6.02-6.07
	//0xe2640ac978ce082a5bdb5316c3b00d11 6.13
}

func GetNewBlock() {

	ctx := context.Background()
	//err, rel := client.GetNewFilterChanges(ctx, "0xe2640ac978ce082a5bdb5316c3b00d11")
	err, rel := Client.GetNewFilterLog(ctx, "0xe2640ac978ce082a5bdb5316c3b00d11")
	if err != nil {
		fmt.Println("getNewBlock", err)
	}
	fmt.Println(rel)

}
