package bc

import (
	"context"
	"fmt"
	"log"
	"time"

	bearychat "github.com/bearyinnovative/bearychat-go"
	"github.com/bearyinnovative/bearychat-go/openapi"
)

const bcToken = "f75c3e3a4cd04ce18cb8f14771eeefcb"

var BcClient *openapi.Client

func Run() {

	// http客户端
	BcClient = openapi.NewClient(bcToken)
	// Rtm
	context, err := bearychat.NewRTMContext(bcToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	err, messageC, errC := context.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("rtm 已启动")

	for {
		select {
		case err := <-errC:
			log.Printf("rtm loop error: %+v", err)
			if err := context.Loop.Stop(); err != nil {
				log.Fatal(err)
			}
			return
		case message := <-messageC:
			// 过滤不是消息的数据
			if !message.IsChatMessage() {
				continue
			}
			// 正在输入
			if message.Type() == bearychat.RTMMessageTypeP2PTyping {
				continue
			}
			// 过滤自己机器人的消息
			if message.IsFromUID(context.UID()) {
				continue
			}

			//
			if NeedManage(message, context.UID()) {
				//
				realText := GetRealText(message, context.UID())
				messageType := GetMessageType(realText)
				//
				switch messageType {
				case HelpType:
					fmt.Println("need help")
					SendHelpMessage(message, context)
				case ListType:
					fmt.Println("need show all list")
					SendLsMessage(message, context)
				case SetType:
					fmt.Println("need set info")
					SendSMessage(message, context)
				case DelType:
					fmt.Println("del user")
					SendDMessage(message, context)
				case DoType:
					fmt.Println("do something")
					SendDoMessage(message, context)
				}

			}
		}
	}
}

func Ctx() context.Context {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(10*time.Second, func() {
		cancel()
	})
	return ctx
}

func NewRefer(m, oldM bearychat.RTMMessage) bearychat.RTMMessage {
	m["refer_key"] = oldM["key"]
	return m
}

func GetRealFileUrl(message bearychat.RTMMessage) (string, error) {
	v3 := ""
	referKey := message["refer_key"].(string)
	messageInfoOptions := &openapi.MessageInfoOptions{
		VChannelID: message["vchannel_id"].(string),
		Key:        openapi.MessageKey(referKey),
	}
	referMessage, _, err := BcClient.Message.Info(Ctx(), messageInfoOptions)
	if err != nil {
		fmt.Println(err.Error())
		return v3, err
	}
	//
	if referMessage.Text != nil {
		v3 = *(referMessage.Text)
	}
	if referMessage.File != nil {
		if referMessage.File.Key != nil {
			fileKey := string(*(referMessage.File.Key))
			fileUrl := BcClient.BaseURL.String() + "file.location?file_key=" + fileKey + "&token=" + bcToken
			return fileUrl, nil
		}
	}
	return v3, nil
}
