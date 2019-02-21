package controller

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"go-driver/blockDriver"
	. "go-driver/common"
	"os"
	"strconv"
)

/*
	1.扫块，读取地址相关的交易
	2.主动发起转币申请，然后链上转币
	3. 1，2 支持合约
	4. 运行时修改运行时参数扩展
*/
//解析后最终秘钥key
var orginKey []byte

func MainEntry() {
	IoBr()
	IoStartLog("启动钱包...")

	//写入runtimeEnv 用户后续控制台来操作运行时配置
	debug := viper.GetBool(`base.debug`)
	GetInstance().SetDebug(debug)
	//解析地址存入单例
	setAddresses()
	//扫块
	searchBlock()
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
func searchBlock() {

	//连接node
	blockDriver.DoEthclientDial()

	allAddresses := GetInstance().GetAllAddresses()

	baseAddress := [][]byte{[]byte(coinbase)}
	account := append(allAddresses, baseAddress...)

	for _, str := range account {
		amount, _ := blockDriver.GetBalance(string(str))
		fmt.Println(string(str), blockDriver.FromWei(amount), "ETH")
	}

	//监听事件，todo 如果filterId过期自动新生成一个filerId
	filterId := blockDriver.WatchNewBlock()
	dealBlockProcess(filterId)

	//24到29

	//daemon()
	//test_sendTransaction()
}

/*
	todo 接受到块后，根据块数量来开启线程。最小10个，最大20个

*/
func dealBlockProcess(filterId string) {

	//开10个线程来接受块
	for i := 0; i < 10; i++ {
		go func() {

		}()
	}
	blocks := blockDriver.GetNewBlock(filterId)
	fmt.Println(blocks)
}

//  eth.sendTransaction({from:"0x0b90ba04fc3520666297a1da31b1f5ff313a475b",to:"0x28172D45396753e4226D1F020849D97eEDB9bcEc",value:web3.toWei(50000,"ether")})
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
