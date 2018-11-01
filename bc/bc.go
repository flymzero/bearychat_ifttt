package bc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	bearychat "github.com/bearyinnovative/bearychat-go"
)

const bcToken = "f75c3e3a4cd04ce18cb8f14771eeefcb"
const bcUrl = "https://api.bearychat.com/v1"

func Run() {
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

func GetMessageInfo(vchannel_id, message_key string) (map[string]interface{}, error) {
	url := bcUrl + "/message.info?token=" + bcToken + "&vchannel_id=" + vchannel_id + "&message_key=" + message_key
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var referMessage map[string]interface{}
	err = json.Unmarshal(body, &referMessage)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if code, exist := referMessage["code"]; exist {
		if err, exist1 := referMessage["error"]; exist1 {
			return nil, errors.New(err.(string))
		} else {
			return nil, errors.New("code : " + strconv.Itoa(code.(int)))
		}
	}
	return referMessage, nil
}
