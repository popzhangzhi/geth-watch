package main

import (
	. "core/common"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-driver/controller"
	"os"
	"os/exec"
	"time"
)

//root命令声明
var rootCmd = &cobra.Command{
	Use: "cr",
	Run: func(cmd *cobra.Command, args []string) {
		println("cr help to see more command")

	},
}

//start命令声明
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cr",
	Run: func(cmd *cobra.Command, args []string) {
		IoStartLog("-------------------------------------------------")
		IoStartLog("startCmd...")
		command := exec.Command("cr_go", "main")
		error := command.Start()
		if error != nil {
			IoStartLogErr("cr_go main", fmt.Sprint(error))
		}
		IoStartLog(fmt.Sprintf("cr start, [PID] %d running...", command.Process.Pid))

		RecordPid([]byte(fmt.Sprintf("%d", command.Process.Pid)))

		//os.Exit(0)

	},
}

//main 命令声明
var mainCmd = &cobra.Command{
	Use:    "main",
	Short:  "The cr main action",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		controller.Generate()
	},
}

//stop命令声明
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop cr",
	Run: func(cmd *cobra.Command, args []string) {

		strb, err := ReadPid()
		if err != nil {

			IoStartLogErr("readPid", fmt.Sprint(err))
		}
		if strb == "0" {
			println("cr is not running")
			os.Exit(0)
		}
		command := exec.Command("kill", strb)
		err = command.Start()

		if err != nil {
			IoStartLogErr("stop cr err:", fmt.Sprint(err))
		} else {
			//清空pid
			ClearPid()
			IoStartLog("cr stop")
		}

		IoStartLog("-------------------------------------------------")
	},
}

//restart命令声明
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart cr",
	Run: func(cmd *cobra.Command, args []string) {
		IoStartLog("-------------------------------------------------")
		IoStartLog("restartCmd...")
		stab, err := ReadPid()
		if err != nil {
			IoStartLogErr("readPid", fmt.Sprint(err))
		}
		if stab != "0" {
			command := exec.Command("cr_go", "stop")
			error := command.Start()
			if error != nil {
				IoStartLogErr("cr_go stop", fmt.Sprint(error))
			}
		}
		command := exec.Command("cr_go", "start")
		error := command.Start()
		if error != nil {
			IoStartLogErr("cr_go start", fmt.Sprint(error))
		}

	},
}

var (
	coinbase = "0x0b90ba04fc3520666297a1da31b1f5ff313a475b"
	//发送地址
	from           = "0x558FcdE4d3949880e0Ab240ba24cDd9f2c46aE1c"
	fromPrivateKey = "8e2cdff2c37ae8aad4c0ff102a84f8f0e0a23549a83cc01598d8089ad82e1a15"

	//接收地址
	address    = "0x032bbB648C56daE9370cA4F97D7D9f6019C84B9c"
	privateKey = "51486722177311552563720459288918193559318459571153646758180554017044071229487"

	address2    = "0xD5806F13709D6B6520f5E66a6969e833A0d98C72"
	privateKey2 = "36d6a41017e583ea93be3771d6084a4b96d1eb19d9a347633c6a154d655c7fcf"
)

//初始化函数
func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(mainCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)

}

/**
读取配置文件
如果没有配置cfgFile文件路径，将会自动获取home目录路径
*/
var RuntimeViper *viper.Viper

//var cfgFile string = "/Users/zhangzhi/go/codeReview/src/core/cr.yaml"
var cfgFile string

func initConfig() {
	// 如果cfgFile不为空，则寻找特定的配置文件路径
	if cfgFile != "" {
		//设置指定的文件，全路径
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.

		//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		//fmt.Println(dir)
		//home, err := homedir.Dir()
		//fmt.Println(home)
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
		//多路径查找配置，查找当前目录
		viper.AddConfigPath(".")
		//查找名字为cr的文件名（不包含扩展名）
		viper.SetConfigName("cr")
	}

	if err := viper.ReadInConfig(); err != nil {
		println("找不到配置文件cr", err)
		os.Exit(1)
	}
	RuntimeViper = viper.GetViper()

}

//主路口运行
func main() {
	// rootCmd.Execute()执行中才会触发执行initConfig（），意味着在下面那个if之前取不到配置文件里的信息
	if err := rootCmd.Execute(); err != nil {

		IoStartLogErr("rootCmd execute", fmt.Sprint(err))

	}

	//controller.Generate()

	//account := [...]string{coinbase, from, address, address2}
	//
	//for _, str := range account {
	//	a, _ := blockDriver.GetBalance(str)
	//	fmt.Println(str, a, "ETH")
	//}

	//blockDriver.WatchNewBlock()

	//24到29
	//blockDriver.GetNewBlock()
	//daemon()

}

/*
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
