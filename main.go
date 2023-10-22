package main

import (
	"example.com/m/cmd"
	"example.com/m/setting"
	"fmt"
)

func main() {
	//配置文件启动
	err := setting.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	//启动程序
	cmd.Execute()
}
