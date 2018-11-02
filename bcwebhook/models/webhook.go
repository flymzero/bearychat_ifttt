package models

import (
	gocontext "context"
	"fmt"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/bearyinnovative/bearychat-go/openapi"
)

var BcClient *openapi.Client

//k名vid
var users map[string]string

func Webhook(ctx *context.Context) {
	if users == nil {
		users = map[string]string{}
	}
	getUserList()
	//
	p2p := createP2P("=bxcZ4")
	if p2p != nil {
		createMessage(p2p)
	}
	//
	ctx.Output.Body([]byte("0"))
}

// 发送一条消息到指定聊天会话。
func createMessage(p2p *openapi.P2P) {
	var opt = &openapi.MessageCreateOptions{
		VChannelID: *(p2p.VChannelID),
		Text:       "哈哈哈h",
	}
	message, _, err := BcClient.Message.Create(Ctx(), opt)
	if err != nil {
		fmt.Println(err.Error())

	}
	fmt.Println(message)
}

// 获取成员列表的id 和 名
func getUserList() {
	userList, _, err := BcClient.User.List(Ctx())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, value := range userList {
		if value.Name != nil && value.ID != nil {
			users[*(value.Name)] = *(value.ID)
		}
	}
	// fmt.Println(users)
}

// 创建一个 P2P 聊天
func createP2P(userId string) *openapi.P2P {
	p2p, _, err := BcClient.P2P.Create(Ctx(), &openapi.P2PCreateOptions{UserID: userId})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return p2p
}

func Ctx() gocontext.Context {
	ctx, cancel := gocontext.WithCancel(gocontext.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})
	return ctx
}
