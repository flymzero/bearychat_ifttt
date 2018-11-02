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
	//
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

// func GetAttachmentTextOrFileRealUrl(orgMessage []byte) (v3 string) {
// 	v3 = ""
// 	// fmt.Println(string(orgMessage))
// 	var msg bearychat.UpdateAttachments
// 	err := json.Unmarshal(orgMessage, &msg)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return v3
// 	}
// 	fmt.Println(msg)
// 	if len(msg.Data.Attachments) == 0 {
// 		fmt.Println(errors.New("没有附件内容"))
// 		return v3
// 	}
// 	v3 = msg.Data.Attachments[0].Text
// 	if msg.Data.Attachments[0].File == nil {
// 		fmt.Println(errors.New("没有附件内容"))
// 		return
// 	}
// 	fileKey := msg.Data.Attachments[0].File.Key
// 	v3 = bcUrl + "/file.location?file_key=" + fileKey + "&token=" + bcToken
// 	return v3
// }

// func GetMessageInfo(vchannel_id, message_key string) (map[string]interface{}, error) {
// 	url := bcUrl + "/message.info?token=" + bcToken + "&vchannel_id=" + vchannel_id + "&message_key=" + message_key
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil, err
// 	}
// 	var referMessage map[string]interface{}
// 	err = json.Unmarshal(body, &referMessage)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil, err
// 	}
// 	if code, exist := referMessage["code"]; exist {
// 		if err, exist1 := referMessage["error"]; exist1 {
// 			return nil, errors.New(err.(string))
// 		} else {
// 			return nil, errors.New("code : " + strconv.Itoa(code.(int)))
// 		}
// 	}
// 	return referMessage, nil
// }
