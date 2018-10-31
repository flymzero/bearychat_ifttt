package main

import (
	"bearychat_ifttt/bc"
	"bearychat_ifttt/config"
	"fmt"
	"os"
)

func main() {
	// config.Users = map[string]config.BCObject{"key": config.BCObject{Nickname: "gg"}}
	// config.WriteUsers("./config/users.json", config.Users)

	// 获取用户储存数据
	if err := config.ReadUsers("./config/users.json", &config.Users); err != nil {
		os.Exit(1)
	}
	fmt.Println("已读取用户配置")
	// 启动rtm
	bc.Run()

}
