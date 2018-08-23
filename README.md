### geth-dirver
    todo:该项目用于调用geth节点，提供api接口给相关业务.与geth节点进行交互。
        对项目进行解耦封包
        已实现创建离线地址
        已实现转账签名，支撑外部地址，内部地址，考虑了nonce问题，gas、gasLimit、gasPrice均可以设置，
        不设置，先有限请求链上接口获取建议值，失败后，会设置默认值
        已实现查询当前余额

        简易webserver接口服务



#### geth-testnet
    该目录为go-ethereum的运行testnet环境时的datadir保存节点目录
##### src/github.com/etherum/go-ethereum
    为了方便调试,为构建项目放在了github下
## 以下为eth节点相关启动和console交互命令

### product environment

    make geth
    运行正式环境，正式环境考虑安全问题适当开放权限
    src/go-ethereum/build/bin/geth  --datadir geth-1.8.10/datadir --rpc --rpcapi "db,eth,net,web3,miner,personal" --fast console

    可配合 attach 命令来进入与 geth 节点进行 js 交互的窗口。不关闭节点
    src/go-ethereum/build/bin/geth  attach rpc:src/geth-1.8.10/datadir/geth.ipc

    也可以通过console 命令可以直接启动 Geth 并进入交互窗口，但是退出时会关闭节点
    src/go-ethereum/build/bin/geth  --datadir src/geth-1.8.10/datadir --rpc --rpcapi "db,eth,net,web3,miner,personal" console

    personal ，admin 为高权限账户，正式环境考虑安全为题

### test environment

    当前时间为2018.6.8号。
    geth最新稳定版本为1.8.10.而我当前的是github上不稳定的包 1.8.11

    -dev 开启私有链
    src/github.com/ethereum/go-ethereum/build/bin/geth  --datadir geth-test/datadir --rpc --rpcapi "db,eth,net,web3,miner,personal"   --dev console


    eth.sendTransaction({from:"0xbaff87a555373dd0358035b77508c41eac84e8c8",to:"0x558FcdE4d3949880e0Ab240ba24cDd9f2c46aE1c",value:web3.toWei(50,"ether")})
    通过日志可以看到，在 dev 模式下，启动节点之后，
    会默认提供一个开发者账号：0x73d8e3e906f64103079cb9331a5274c288c633f5，

    --testnet 开启测试链，加入eth 测试节点
    src/github.com/ethereum/go-ethereum/build/bin/geth  --datadir geth-testnet/datadir --rpc --rpcapi "db,eth,net,web3,miner,personal"   --testnet

    通过geth.ipc进入JavaScript控制台
    src/github.com/ethereum/go-ethereum/build/bin/geth  attach rpc:geth-testnet/datadir/geth.ipc

    --rpcaddr "192.168.0.29" 设置暴露对外的ip节点

    miner.start(1)
    miner.stop()

    govendor add +e 增加所有外包依赖，手动复制C的libsecp256k1文件夹，增加了所有geth的依赖

    go调用包的熟悉，优先本项目->该路劲下的vendor->该上级目录的的包（gopath的外部包）



### javascript操作命令
    这个账号会作为当前的 coinbase 账号，在 keystore 目录下也有对应的加密私钥文件。

    eth.accounts
    eth.blockNumber
    eth.getBlock(11)
    eth.getTransaction(hash)
    eth.getTransactionReceipt(hash)
    personal.newAccount("zhangzhi") //0x8408cbf887d3de46a76c83db0538231e5aa4cdb4
    personal.newAccount("123456") //0x83a7fb46762881a4b73bf1f2da7eae8b7809b50f
    personal.newAccount()

    eth.getBalance("0x83a7fb46762881a4b73bf1f2da7eae8b7809b50f")
    admin = eth.accounts[0]
    eth.getBalance(eth.accounts[0])

    eth.sendTransaction({from:"0xbaff87a555373dd0358035b77508c41eac84e8c8",to:"0x5DE5c59330DAe02A48dE228AAdbfDAFCC0b10e1C",value:web3.toWei(50,"ether")})

    eth.sendTransaction({from:"0x8408cbf887d3de46a76c83db0538231e5aa4cdb4",to:"0x83a7fb46762881a4b73bf1f2da7eae8b7809b50f",value:web3.toWei(50,"ether")})
### 个人权限操作交易，输入密码形式，慎用
    personal.sendTransaction({from:"0x59949d4b02d02161b4e5df59943027c1ea2bfbb1",to:"0x24f89f5c62ea5edf4f39eff6e096dbae5540cc34",value:web3.toWei(104,"ether")},'123456')

    personal.unlockAccount("0x8408cbf887d3de46a76c83db0538231e5aa4cdb4")
    personal.lockAccount("0x8408cbf887d3de46a76c83db0538231e5aa4cdb4")

    0xc05ac42237742f5038e72fdbbd930815f6d79034ee862f0b3e711d1146cf7c36
## geth源码相关信息记录

正式地址 0xc62b0b4d09144659eb65a435bf82f458bf043ae3


## geth 目录结构
    account  该包实现了高层级的Ethereum账号管理
    account/abi 该包实现了Ethereum的ABI(应用程序二进制接口)
    account/abi/bind 该包生成Ethereum合约的Go绑定
    account/abi/bind/backends	        --
    account/keystore 实现了Secp256k1私钥的加密存储
    account/usbwallet	该包实现了支持USB硬件钱包

    btm	该包实现了二叉merkle树

    cmd/abigen	--
    cmd/bootnode	该节点为Ethereum发现协议运行一个引导节点
    cmd/ethkey
    cmd/evm	执行EVM代码片段
    cmd/faucet	faucet是以太faucet支持的轻量级客户
    cmd/geth	geth是Ethereum的官方客户端命令行
    cmd/p2psim	p2psim为客户端命令行模拟	HTTP API
    cmd/puppeth	puppeth是一个命令组装和维护私人网路
    cmd/rlpdump	rlpdump能更好的打印出RLP格式的数据
    cmd/swarm	bzzhash命令能够更好的计算出swarm哈希树
    cmd/utils	为Go-Ethereum命令提供说明
    cmd/wnode	--
    common	包含一些帮助函数
    common/bitutil	该包实现快速位操作
    common/compiler	包装了Solity编译器可执行文件
    common/fdllimit	--
    common/hexutil	以0x为前缀的十六进制编码
    common/math	数学工具
    common/number	--
    compression/rle	实现run-length encoding编码用于Ethereum数据
    comsensus	实现了不同以太共识引擎
    comsensus/clique	实现了权威共识引擎
    comsensus/ethash	发动机工作的共识ethash证明
    comsensus/misc	--
    console	--
    contracts/chequebook	'支票薄'以太智能合约
    contracts/chequebook/contract    --
    contracts/ens	--
    contracts/ens/contract --
    core 实现以太合约接口
    core/asm	汇编和反汇编接口
    core/bloombits	Bloom过滤批量数据
    core/state	封装在以太状态树之上的一种缓存结构
    core/types	以太合约支持的数据类型
    core/vm	以太虚拟机
    core/vm/runtime	一种用于执行EVM代码的基本执行模型
    crypto	--
    crypto/bn256	最优的ATE配对在256位Barreto-Naehrig曲线上
    crypto/bn256/cloudflare	在128位安全级别上的特殊双线性组
    crypto/bn256/google	在128位安全级别上的特殊双线性组
    crypto/ecies	--
    crypto/randentropy	--
    crypto/secp256k1 封装比特币secp256k1的C库
    crypto/sha3	Sha-3固定输出长度散列函数 and 由FIPS-202定义的抖动变量输出长度散列函数

    dashboard	        --

    eth 以太坊协议
    ethclient	以太坊RPC AIP客户端
    ethdb --
    eth/downloader	手动全链同步
    eth/fetcher	基于块通知的同步
    eth/filters	用于区块，交易和日志事件的过滤
    eth/gasprice	--
    eth/stats	网络统计报告服务
    eth/tracers	收集JavaScript交易追踪
    event	处理时时事件的费用
    event/filter 事件过滤

    internal/build	--
    internal/cmdtest	        --
    internal/debug	调试接口Go运行时调试功能
    internal/ethapi	常用的以太坊API函数
    internal/guide	小测试套件，以确保开发指南工作中的代码段正常运行
    internal/jsre	JavaScript执行环境
    internal/jsre/deps	控制台JavaScript依赖项Go嵌入
    internal/web3ext	geth确保web3.js延伸

    les	轻量级Ethereum子协议
    les/flowcontrol	客户端流程控制机制
    light	EtalumLight客户端实现按需检索能力的状态和链对象

    log	log输出日志
    log/term                                    --

    metrics	Coda Hale度量库的Go端口
    metrics/exp 表达式相关操作
    metrics/influxdb	        --
    metrics/librato	--
    miner	以太坊块创建和挖矿
    mobile	geth的移动端API

    node	设置多维接口节点

    p2p p2p网络协议
    p2p/discover	节点发现协议
    p2p/discv5	RLPx v5主题相关的协议
    p2p/enr 实现EIP-778中的以太坊节点记录
    p2p/nat	提供网络端口映射协议的权限
    p2p/netutil 网络包拓展
    p2p/protocols p2p拓展
    p2p/simulations	实现模拟p2p网络
    p2p/simulations/adapters	        --
    p2p/simulations/examples	        --
    p2p/testing	--
    params	--

    rlp RLP系列化格式
    rpc 通过网络或者I/O链接来访问接口

    swarm	--
    swarm/api --
    swarm/api/client --
    swarm/api/http HTML格式错误的处理
    swarm/fuse	--
    swarm/metrics --
    swarm/network --
    swarm/network/kademlia --
    swarm/services/swap	--
    swarm/services/swap/swap --
    swarm/storage	--
    swarm/testutil	--

    tests	以太坊JSON测试
    trie	Merkle Patricia树实现

    vendor/gopkg.in/check.v1 Go更深的测试
    vendor/gopkg.in/urfave/cli.v1 Go命令行应用的框架

    whisper/mailserver	--
    whisper/shhclient	--
    whisper/whisperv2	Whisper Poc-1实现
    whisper/whisperv5	Whisper协议(版本5)
    whisper/whisperv6	--


