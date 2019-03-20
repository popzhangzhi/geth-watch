package controller

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/viper"
	"go-driver/blockDriver"
	. "go-driver/common"
	"math/big"
	"os"
	"strconv"
)

/*
	1.扫块，读取地址相关的交易
	2.主动发起转币申请，然后链上转币
	3. 1，2 支持合约
	4. 运行时修改运行时参数扩展
*/
//filterid
var filterId = ""

func MainEntry() {
	IoBr()
	IoStartLog("启动钱包...")

	//写入runtimeEnv 用户后续控制台来操作运行时配置,注意yaml文件字符串可以省略双引号，但是key与value间的空格不能少，少了
	//无法解析成功该文件
	debug := viper.GetBool(`base.debug`)
	GetInstance().SetDebug(debug)
	//解析地址存入单例
	setAddresses()
	//链接rpc-node
	blockDriver.DoEthclientDial()

	getDebugInfo()
	//扫块

	initWatchBlock()
}

func getDebugInfo() {
	//连接node
	blockDriver.DoEthclientDial()

	allAddresses := GetInstance().GetAllAddresses()
	//加入测试节点主地址
	baseAddress := [][]byte{[]byte(coinbase)}

	account := append(allAddresses, baseAddress...)

	for k, str := range account {
		amount, _ := blockDriver.GetBalance(string(str))
		if string(str) == coinbase {
			fmt.Println("coinbase")
		}
		fmt.Println(k, string(str), blockDriver.FromWei(amount), "ETH")

	}
}

/*
	解析地址和秘钥，存入运行时。
*/
func setAddresses() {

	address, _ := IoReadFile(addressFile)
	addresses := bytes.Split(bytes.TrimSpace(address), []byte("\n"))
	addressLen := len(addresses)
	debug := GetInstance().Debug

	var secretKey []byte

	if debug {
		//debug
		//ioutil.WriteFile("orginKey_debug", secretKey, 0666)
		secretKey, _ = IoReadFile("orginKey_debug")
	} else {
		//启动输出密码解密秘钥
		pwd := inputPwd(1)
		secretKey = DecodeSecretKey(pwd)
	}

	orginKeys := bytes.Split(bytes.TrimSpace(secretKey), []byte("\n"))
	orginKeyLen := len(orginKeys)

	if addressLen != orginKeyLen {
		fmt.Println(`地址和秘钥不匹配`, addressLen, orginKeyLen)
		IoStartLog(`地址和秘钥不匹配`)
		os.Exit(1)
	}
	//设置秘钥到运行时
	key := make(map[string][]byte)
	for k, v := range addresses {
		key[string(v)] = orginKeys[k]
	}
	//存入单例
	GetInstance().SetAddresses(addresses[0:systemNum], addresses, key)

	IoStartLog("成功解析地址数量" + strconv.Itoa(orginKeyLen))

}

/*
 扫块
*/
func initWatchBlock() {

	//获取起始块
	blockStart := viper.GetInt(`wallet.height`)
	IoStartLog("扫块起始块高" + string(blockStart))
	blockStartNumber := big.NewInt(int64(blockStart))
	//获取终止块
	blockEndNumber := blockDriver.GetCurrentBlockNumber()

	watchBlock(blockStartNumber, blockEndNumber)

	//test_sendTransaction()

}

/**
特定项查找链上数据
*/
func getFilterData() {
	//监听事件，当前只扫相关的块,filterid存在，开启getfilterchange线程，否则再次生词filterid在开启getfilterchange线程
checkFilterId:
	blockHeight := viper.GetInt(`wallet.height`)
	blockHeightOX := hexutil.EncodeBig(big.NewInt(int64(blockHeight)))
	//GetInstance().GetAllAddressesToString()
	object := blockDriver.FileterObject{
		blockHeightOX, "latest", coinbase, nil}
	filterId = blockDriver.WatchFilterBlock(object)
	if filterId == "" {
		goto checkFilterId
	} else {
		fmt.Println("filterId", filterId)
		blockDriver.GetFilterData(filterId)
	}

}

/*
	todo 接受到块后，根据块数量来开启线程。最小10个，最大20个

*/
func watchBlock(blockHeight *big.Int, blockEndNumber *big.Int) {

	if blockHeight.Cmp(blockEndNumber) >= 0 {
		//起始块大于终止块，return
		IoStartLogErr("watchBlock", "起始块大于终止块退出扫块")
		return
	}
	//声明协程池
	p := NewRouinePoor(10)
	start := blockHeight.Int64()
	end := blockEndNumber.Int64()

	//赋值块高和协程允许的闭包函数到channel
	go func() {

		for i := start; i <= end; i++ {

			arg := make(map[string]string)
			arg["blockNumber"] = strconv.FormatInt(i, 10)
			arg["taskId"] = arg["blockNumber"]

			t := NewTask(func(params map[string]string) {
				fmt.Println("workId:" + params["workId"] + " 当前块高:" + params["blockNumber"])
			}, arg)

			p.ReveiceChannel <- t

		}
		//赋值成功后关闭channel
		p.Close()

	}()

	p.Run()

}

//eth.sendTransaction({from:"0x0b90ba04fc3520666297a1da31b1f5ff313a475b",to:"0x28172D45396753e4226D1F020849D97eEDB9bcEc",value:web3.toWei(50000,"ether")})
func test_sendTransaction() {
	allAddress := GetInstance().GetAllAddresses()
	key := GetInstance().GetKey(string(allAddress[0]))
	txid := blockDriver.SendRowTransaction(string(allAddress[0]), string(key), string(allAddress[1]), `1`, 0, 0)

	fmt.Println(txid)
}

var (
	coinbase = "0x0b90ba04fc3520666297a1da31b1f5ff313a475b"
	//发送地址
	from = "0x28172D45396753e4226D1F020849D97eEDB9bcEc"

	//接收地址
	address = "0xd5404e27a125434b6390B67D40B6697C62A3D131"

	address2 = "0xfBeC43A114d8412C1E5Fc122d85164b81bbcfCF2"
)
