package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	. "go-driver/common"
	"go-driver/controller"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const commandName = "gethWatch"

//root命令声明
var rootCmd = &cobra.Command{
	Use: commandName,
	Run: func(cmd *cobra.Command, args []string) {
		//println("输入 " + commandName + " help 来查看更多命令")

	},
}

//start命令声明
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start " + commandName,
	Run: func(cmd *cobra.Command, args []string) {
		IoBr()
		IoStartLog("startCmd...")
		command := exec.Command(commandName, "main")
		error := command.Start()
		if error != nil {
			IoStartLogErr("main", fmt.Sprint(error))
		}
		IoStartLog(fmt.Sprintf(commandName+" start, [PID] %d running...", command.Process.Pid))

		RecordPid([]byte(fmt.Sprintf("%d", command.Process.Pid)))

		//os.Exit(0)

	},
}

//main 命令声明
var mainCmd = &cobra.Command{
	Use:    "main",
	Short:  "The " + commandName + " main action",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		controller.MainEntry()
	},
}

//stop命令声明
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop " + commandName,
	Run: func(cmd *cobra.Command, args []string) {

		strb, err := ReadPid()
		if err != nil {
			IoStartLogErr("readPid", fmt.Sprint(err))
		}
		if strb == "0" {
			println(commandName + " is not running")
			os.Exit(0)
		}
		command := exec.Command("kill", strb)
		err = command.Start()

		if err != nil {
			IoStartLogErr("stop "+commandName+" err:", fmt.Sprint(err))
		} else {
			//清空pid
			ClearPid()
			IoStartLog(commandName + " stop")
		}

		IoBr()
	},
}

//restart命令声明
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart " + commandName,
	Run: func(cmd *cobra.Command, args []string) {
		IoBr()
		IoStartLog("restartCmd...")
		stab, err := ReadPid()
		if err != nil {
			IoStartLogErr("readPid", fmt.Sprint(err))
		}
		if stab != "0" {
			command := exec.Command(commandName, "stop")
			error := command.Start()
			if error != nil {
				IoStartLogErr(commandName+" stop", fmt.Sprint(error))
			}
		}
		command := exec.Command(commandName, "start")
		error := command.Start()
		if error != nil {
			IoStartLogErr(commandName+" start", fmt.Sprint(error))
		}

	},
}

var GenerateAddrCmd = &cobra.Command{
	Use:   "generateAddr",
	Short: "generateAddr [num] 建议测试时生成1000",
	Run: func(cmd *cobra.Command, args []string) {

		IoBr()
		IoStartLog("generateAddr 开始...")
		num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("输入数量非整形", err)
			os.Exit(1)
		}
		controller.Generate(num)
		IoStartLog("generateAddr 生成地址完成")
		IoBr()
	},
}

//初始化函数
func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(mainCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(GenerateAddrCmd)

}

/**
读取配置文件
如果没有配置cfgFile文件路径，将会自动获取home目录路径
*/

var cfgFile string

func initConfig() {
	// 如果cfgFile不为空，则寻找特定的配置文件路径
	if cfgFile != "" {
		//设置指定的文件，全路径
		viper.SetConfigFile(cfgFile)
	} else {
		//多路径查找配置，查找当前目录
		//当前目录指的go install的目录。e.g本例是在get-watch/目录下
		viper.AddConfigPath("./")

		//查找名字为cr的文件名（不包含扩展名）
		viper.SetConfigName("cr")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("找不到配置文件cr.yaml", err)
		os.Exit(1)
	}

}

//主路口运行
func main() {
	// rootCmd.Execute()执行中才会触发执行initConfig（），意味着在下面那个if之前取不到配置文件里的信息
	if err := rootCmd.Execute(); err != nil {

		IoStartLogErr("rootCmd execute", fmt.Sprint(err))

	}

	//生成20w地址耗时3分钟，mac 4核
	//controller.Generate(3)
	//解密20w秘钥耗时1分钟，环境同上

	controller.MainEntry()

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
