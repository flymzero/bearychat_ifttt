package bc

import (
	"fmt"
	"log"

	bearychat "github.com/bearyinnovative/bearychat-go"
)

const BcToken = "880b2f6f5ad949689aa5fcca5a874f53"
const BcUrl = "https://api.bearychat.com/v1"

func Run() {
	context, err := bearychat.NewRTMContext(BcToken)
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
