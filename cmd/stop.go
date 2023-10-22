package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

// 关闭程序
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止客服http服务",
	Run: func(cmd *cobra.Command, args []string) {
		pids, err := ioutil.ReadFile(viper.GetString("project.portFile"))
		if err != nil {
			return
		}
		pidSlice := strings.Split(string(pids), ",")
		var command *exec.Cmd
		for _, pid := range pidSlice {
			if runtime.GOOS == "windows" {
				command = exec.Command("taskkill.exe", "/f", "/pid", pid)
			} else {
				fmt.Printf("关闭pid %s", pid)
				command = exec.Command("kill", pid)
			}
			command.Start()
		}
	},
}
