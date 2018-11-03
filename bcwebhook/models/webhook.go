package models

import (
	gocontext "context"
	"errors"
	"fmt"
	"strings"
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
	//
	requestText := string(ctx.Input.RequestBody)
	if len(requestText) == 0 {
		return
	}
	name, text, err := getRequestInfo(requestText)
	if err != nil {
		return
	}
	// 更新用户列表
	getUserList()
	//
	id, exist := users[name]
	if !exist {
		fmt.Println("找不到对应用户", name)
		return
	}
	// 创建聊天会话
	p2p := createP2P(id)
	if p2p != nil {
		// 发送消息
		createMessage(p2p, text)
	}
}

func getRequestInfo(requestText string) (name, text string, err error) {
	array1 := strings.Fields(requestText)
	if len(array1) == 0 {
		err = errors.New("request text is err")
		fmt.Println(err.Error())
		return "", "", err
	}
	name = array1[0]
	text = ""
	array2 := strings.SplitN(requestText, " ", 2)
	if len(array2) >= 2 {
		text = array2[1]
	}
	return name, text, nil
}

// 发送一条消息到指定聊天会话。
func createMessage(p2p *openapi.P2P, text string) {
	var opt = &openapi.MessageCreateOptions{
		VChannelID: *(p2p.VChannelID),
		Text:       "🤖 消息 : " + text,
	}
	_, _, err := BcClient.Message.Create(Ctx(), opt)
	if err != nil {
		fmt.Println(err.Error())
	}
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
