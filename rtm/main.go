package main

import (
	"bearychat_ifttt/rtm/bc"
	"bearychat_ifttt/rtm/config"
	"fmt"
	"os"
)

func main() {
	// cYm2vsJ3YxJOJdOX7V0kQG
	// carsh with https://github.com/maxcnunes/gaper
	// go get -u github.com/maxcnunes/gaper/cmd/gaper
	// 获取用户储存数据
	if err := config.ReadUsers("./config/users.json", &config.Users); err != nil {
		os.Exit(1)
	}
	fmt.Println("已读取用户配置")
	// 启动rtm
	bc.Run()

}
