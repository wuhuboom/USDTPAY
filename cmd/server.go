/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package cmd

import (
	"example.com/m/common"
	"example.com/m/dao/mysql"
	"example.com/m/dao/redis"
	"example.com/m/logger"
	"example.com/m/process"
	"example.com/m/router"
	"example.com/m/tools"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zh-five/xdaemon"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	port     string
	daemon   bool
	rootPath string
)
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "启动服务器",
	Example: "server",
	Run:     run,
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&rootPath, "rootPath", "r", "", "程序根目录")
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", viper.GetString("project.listeningPort"), "监听端口号")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "是否为守护进程模式")
}

func run(cmd *cobra.Command, args []string) {
	//初始化目录
	initDir()
	//初始化守护进程
	initDaemon()
	if noExist, _ := tools.IsFileNotExist(common.LogDirPath); noExist {
		if err := os.MkdirAll(common.LogDirPath, 0777); err != nil {
			log.Println(err.Error())
		}
	}
	isMainUploadExist, _ := tools.IsFileExist(common.UploadDirPath)
	if !isMainUploadExist {
		os.Mkdir(common.UploadDirPath, os.ModePerm)
	}

	//初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("日志初始化失败", err)
		return
	}
	defer zap.L().Sync() //缓存日志追加到日志文件中

	//链接数据库
	if err := mysql.Init(); err != nil {
		fmt.Println("mysql 链接失败,", err)
		return
	}
	defer mysql.Close()
	//redis 初始化
	//4.初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Println("redis文件初始化失败：", err)
		return
	}
	defer redis.Rdb.Close()

	//注册进程
	go process.CratedPoolAddress(mysql.DB)
	go process.CheckLastGetMoneyTime(mysql.DB)

	router.Setup()
}

// 初始化目录
func initDir() {
	if rootPath == "" {
		rootPath = tools.GetRootPath()
	}
	log.Println("程序运行路径:" + rootPath)
	common.RootPath = rootPath
	common.LogDirPath = rootPath + "/logs/"
	common.ConfigDirPath = rootPath + "/config/"
	common.StaticDirPath = rootPath + "/static/"
	common.UploadDirPath = rootPath + "/static/upload/"
}

// 初始化守护进程
func initDaemon() {
	//启动前 先杀死之前的 程序
	pid, err := ioutil.ReadFile(viper.GetString("project.portFile"))
	if err == nil {
		pidSlice := strings.Split(string(pid), ",")
		var command *exec.Cmd
		for _, pid := range pidSlice {
			if runtime.GOOS == "windows" {
				command = exec.Command("taskkill.exe", "/f", "/pid", pid)
			} else {
				fmt.Println("成功结束进程:", pid)
				command = exec.Command("kill", pid)
			}
			command.Start()
		}
	}
	if daemon == true {
		d := xdaemon.NewDaemon(common.LogDirPath + common.LogFileName)
		d.MaxError = 10
		d.Run()
	}
	//记录pid
	ioutil.WriteFile(common.RootPath+"/"+viper.GetString("project.portFile"), []byte(fmt.Sprintf("%d,%d", os.Getppid(), os.Getpid())), 0666)

}
