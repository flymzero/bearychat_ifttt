package main

import (
	"bearychat_ifttt/bcwebhook/models"
	_ "bearychat_ifttt/bcwebhook/routers"

	"github.com/astaxie/beego"
	"github.com/bearyinnovative/bearychat-go/openapi"
)

const bcToken = "f75c3e3a4cd04ce18cb8f14771eeefcb"

func main() {
	// http客户端
	models.BcClient = openapi.NewClient(bcToken)
	//
	beego.Run()
}
