package bc

import (
	"bearychat_ifttt/config"
	"fmt"
	"regexp"
	"strings"

	bearychat "github.com/bearyinnovative/bearychat-go"
)

type MessageType string

const (
	NullType MessageType = ""
	HelpType MessageType = "-h"
	ListType MessageType = "-ls" //列出所有自己的成员
	SetType  MessageType = "-s"  //设置对象信息
	DelType  MessageType = "-d"  //删除对象信息
	DoType   MessageType = "$"   //触发操作
)

const helpText = `[bearychat_ifttt](https://github.com/flymzero/bearychat_ifttt) 使用帮助

**更新**

> 去除email数据绑定，一个对象只绑定**名称**和**key**，名称唯一

**对象数据绑定操作**

把自己或者别人（统一称为对象）先绑定对应的ifttt的key执行操作
一个对象包含**名称**和**key**两个字段 （名称唯一）

> **-h**    获取使用帮助
> **-ls**   列出所有自己及自己设置对象的信息(名称&key)
> **-s [-m]  n:名称  k:ifttt上的key **   设置自己或者对象的信息(n必填唯一，-m则表示设置自己信息)
- 例1 : -s -m n:我 k:xxxxxx  设置自己的名称称和ifttt的key
- 例2 : -s n:老婆 k:oooooo 设置一个对象：老婆及触发的key
> **-d 名称**  删除你对象中对应名称的这个对象
- 例 : -d 老婆 把老婆这个对象删了

**触发操作 ifttt webhooks 触发词 + 任意服务**
> **$触发词  [n:名称]  [v1:xx]  [v2:xx]  [v3:xx] ** 
对这个名称对象(n不填就是自己，v1,v2,v3都是可选的)，进行触发操作,并传输可选的3个参数，~~当存在引用附件时以附件的url优先作为v3的值~~,贝洽的文件在未登录的情况下无法访问链接,所以无法转存!!!
需在"对象"手机ifttt上创建对应的Applet

**基本上ifttt上能创造的东西(通知,发邮件,文件转存,发微博....),你都可以让这个机器人代劳, 话说ifttt不开放共享applet也真是坑**
更多使用方法见 https://github.com/flymzero/bearychat_ifttt
`

const noUserInfo = "没有关于您的用户信息请先设置 命令：**-s**"
const cmdError = "❎ 命令不匹配"
const keyNeedP2p = "设置命名中包含Ifttt的key，涉及加密信息，请私聊我添加"
const NeedSetMeFirst = "❎ 设置命名必须先设置自己的信息才能添加其他对象 命令：**-s -m**"
const SetNickNameError = "❎ 设置命名中**n:**没有数据"
const SetNickNameOtherAgin = "❎ 你的对象中有相同名称存在"
const SetNickNameSelfAgin = "❎ 和你的自己的名称一样"
const DelNotExistError = "❎ 你对象中不存在这个名称的对象(不能删除自己)"
const DoSelfKeyError = "❎ 请先设置自己的Ifttt的Key"
const DoOtherKeyError = "❎ 请先设置该对象的Ifttt的Key"
const DoOtherNotExist = "❎ 你对象中不存在这个名称的对象"

// 查看是否需要处理
func NeedManage(message bearychat.RTMMessage, uid string) bool {
	if message.Type() == bearychat.RTMMessageTypeP2PMessage {
		return true
	}
	if mentioned, _ := message.ParseMentionUID(uid); mentioned { //讨论组并且@机器人或者私聊
		return true
	}
	return false
}

// 消息去除@机器人字段
func getText(message bearychat.RTMMessage, uid string) string {
	text := message.Text()
	if message.Type() == bearychat.RTMMessageTypeP2PMessage {
		var mentionUserRegex = regexp.MustCompile("@<=(=[A-Za-z0-9]+)=> ")
		locs := mentionUserRegex.FindAllStringSubmatchIndex(text, -1)
		if len(locs) != 0 {
			for _, loc := range locs {
				// "@<==1=> xxx" -> [0 8 3 5]
				// [3:5] "=1" [8:] "xxx"
				if text[loc[2]:loc[3]] == uid {
					text = text[loc[1]:]
					break
				}
			}
		}
	} else if mentioned, content := message.ParseMentionUID(uid); mentioned { //讨论组并且@机器人或者私聊
		text = content
	}
	return text

}

// 获取 对话机器人消息text
func GetRealText(message bearychat.RTMMessage, uid string) string {
	text := getText(message, uid)
	realTextArray := strings.Fields(text)
	if len(realTextArray) == 0 {
		return ""
	} else {
		return realTextArray[0]
	}
}

// 获取  对话机器人 请求类型
func GetMessageType(realText string) MessageType {
	if strings.HasPrefix(realText, string(DoType)) {
		return DoType
	}
	//
	switch MessageType(realText) {
	case NullType:
		return HelpType

	case HelpType:
		return HelpType

	case ListType:
		return ListType

	case SetType:
		return SetType

	case DelType:
		return DelType

	default:
		return HelpType
	}
}

// 加工 回复信息
func addCommonMessageInfo(message bearychat.RTMMessage, sendMessage *bearychat.RTMMessage) {
	(*sendMessage)["vchannel_id"] = message["vchannel_id"]
	if message.IsP2P() {
		(*sendMessage)["type"] = bearychat.RTMMessageTypeP2PMessage
		(*sendMessage)["to_uid"] = message["uid"]
	} else {
		(*sendMessage)["type"] = bearychat.RTMMessageTypeChannelMessage
		(*sendMessage)["channel_id"] = message["channel_id"]
	}
}

// -h 发送
func SendHelpMessage(message bearychat.RTMMessage, context *bearychat.RTMContext) {
	m := bearychat.RTMMessage{}
	addCommonMessageInfo(message, &m)
	m["text"] = helpText
	if err := context.Loop.Send(m); err != nil {
		fmt.Println(err.Error())
	}
}

func setLsStr(key, value string) string {
	if key == "iftttkey" && len(value) > 0 {
		return value //"㊙️"  参赛期间可以看到key
	} else if len(value) > 0 {
		return value
	}
	return "空"
}

// -ls 发送
func SendLsMessage(message bearychat.RTMMessage, context *bearychat.RTMContext) {
	m := bearychat.RTMMessage{}
	addCommonMessageInfo(message, &m)
	//寻找用户信息
	userInfo, exist := config.Users[message["uid"].(string)]
	if !exist {
		m["text"] = noUserInfo
	} else {
		infoText := "[自己]   名称 : " + setLsStr("", userInfo.Nickname) + ",   key : " + setLsStr("iftttkey", userInfo.IftttKey) + "\n"
		if userInfo.Others != nil {
			for _, value := range userInfo.Others {
				tempText := "[对象]   名称 : " + setLsStr("", value.Nickname) + ",   key : " + setLsStr("iftttkey", value.IftttKey) + "\n"
				infoText += tempText
			}
		}
		m["text"] = infoText + "\n参赛期间可以看到key,方便测试"
	}
	if err := context.Loop.Send(m); err != nil {
		fmt.Println(err.Error())
	}
}

// -d 删除
func SendDMessage(message bearychat.RTMMessage, context *bearychat.RTMContext) {
	m := bearychat.RTMMessage{}
	addCommonMessageInfo(message, &m)
	text := getText(message, context.UID())
	realTextArray := strings.Fields(text)
	if len(realTextArray) != 2 {
		m["text"] = cmdError
	}
	nickname := realTextArray[1]
	//寻找用户信息
	userInfo, exist := config.Users[message["uid"].(string)]
	if !exist {
		m["text"] = DelNotExistError
	} else {
		if userInfo.Others == nil {
			m["text"] = DelNotExistError
		} else {
			isExist := false
			for k, _ := range userInfo.Others {
				if k == nickname {
					isExist = true
					break
				}
			}
			if !isExist {
				m["text"] = DelNotExistError
			} else {
				others := userInfo.Others
				delete(others, nickname)
				userInfo.Others = others
				config.Users[message["uid"].(string)] = userInfo
				config.WriteUsers("./config/users.json", config.Users)
				m["text"] = "✅   删除对象成功"
			}
		}
	}
	//
	if err := context.Loop.Send(m); err != nil {
		fmt.Println(err.Error())
	}
}

// -s 设置
func SendSMessage(message bearychat.RTMMessage, context *bearychat.RTMContext) {
	m := bearychat.RTMMessage{}
	addCommonMessageInfo(message, &m)
	text := getText(message, context.UID())
	realTextArray := strings.Fields(text)
	fmt.Println(realTextArray)
	if len(realTextArray) <= 1 || realTextArray[0] != string(SetType) {
		m["text"] = cmdError
	} else {
		isSelf := false
		nickname := ""
		key := ""
		for _, value := range realTextArray[1:] {
			if value == "-m" {
				isSelf = true
			} else {
				if strings.HasPrefix(value, "n:") && len(value) > 2 {
					nickname = value[2:]
				} else if strings.HasPrefix(value, "k:") && len(value) > 2 {
					//参赛期间可以看到key,方便测试
					// if !message.IsP2P() {
					// 	m["text"] = keyNeedP2p
					// 	if err := context.Loop.Send(m); err != nil {
					// 		fmt.Println(err.Error())
					// 	}
					// 	return
					// }
					key = value[2:]
				}
			}
		}
		//
		//寻找用户信息
		userInfo, exist := config.Users[message["uid"].(string)]
		if !isSelf && !exist {
			m["text"] = NeedSetMeFirst
		} else if isSelf {
			if len(nickname) > 0 {
				userInfo.Nickname = nickname
				//看名称是否重名
				if userInfo.Others != nil {
					for k, _ := range userInfo.Others {
						if k == nickname {
							m["text"] = SetNickNameOtherAgin
							if err := context.Loop.Send(m); err != nil {
								fmt.Println(err.Error())
							}
							return
						}
					}
				}
			}
			if len(key) > 0 {
				userInfo.IftttKey = key
			}
			config.Users[message["uid"].(string)] = userInfo
			config.WriteUsers("./config/users.json", config.Users)
			infoText := "[自己]   名称 : " + setLsStr("", userInfo.Nickname) + ",   key : " + setLsStr("iftttkey", userInfo.IftttKey) + "\n"
			m["text"] = "✅   " + infoText + "\n参赛期间可以看到key,方便测试"
		} else {
			if len(nickname) == 0 {
				m["text"] = SetNickNameError
			} else {
				//看名称是否重名
				if userInfo.Nickname == nickname {
					m["text"] = SetNickNameSelfAgin
					if err := context.Loop.Send(m); err != nil {
						fmt.Println(err.Error())
					}
					return
				}
				otherInfo, otherExist := userInfo.Others[nickname]
				if !otherExist {
					otherInfo = config.OtherObject{}
				}
				if len(nickname) > 0 {
					otherInfo.Nickname = nickname
				}
				if len(key) > 0 {
					otherInfo.IftttKey = key
				}
				if userInfo.Others == nil {
					userInfo.Others = map[string]config.OtherObject{}
				}
				userInfo.Others[nickname] = otherInfo
				config.Users[message["uid"].(string)] = userInfo
				config.WriteUsers("./config/users.json", config.Users)
				infoText := "[对象]   名称 : " + setLsStr("", otherInfo.Nickname) + ",   key : " + setLsStr("iftttkey", otherInfo.IftttKey) + "\n"
				m["text"] = "✅   " + infoText + "\n参赛期间可以看到key,方便测试"
			}
		}
	}
	//
	if err := context.Loop.Send(m); err != nil {
		fmt.Println(err.Error())
	}
}

// do
func SendDoMessage(message bearychat.RTMMessage, context *bearychat.RTMContext) {
	m := bearychat.RTMMessage{}
	addCommonMessageInfo(message, &m)
	text := getText(message, context.UID())
	realTextArray := strings.Fields(text)
	//寻找用户信息
	userInfo, exist := config.Users[message["uid"].(string)]
	if !exist {
		m["text"] = NeedSetMeFirst
		if err := context.Loop.Send(m); err != nil {
			fmt.Println(err.Error())
		}
		return
	} else if len(realTextArray) == 0 || !strings.HasPrefix(realTextArray[0], string(DoType)) || realTextArray[0] == string(DoType) {
		m["text"] = cmdError
		if err := context.Loop.Send(m); err != nil {
			fmt.Println(err.Error())
		}
		return
	} else {
		trigger := realTextArray[0][1:]
		nickname := ""
		key := ""
		v1 := ""
		v2 := ""
		v3 := ""
		for _, value := range realTextArray[1:] {
			if strings.HasPrefix(value, "n:") && len(value) > 2 {
				nickname = value[2:]
			} else if strings.HasPrefix(value, "v1:") && len(value) > 3 {
				v1 = value[3:]
			} else if strings.HasPrefix(value, "v2:") && len(value) > 3 {
				v2 = value[3:]
			} else if strings.HasPrefix(value, "v3:") && len(value) > 3 {
				v3 = value[3:]
			}
		}
		// nickname可能是对象看看对象存在挖，存在取key
		if len(nickname) == 0 { //自己
			if len(userInfo.IftttKey) == 0 {
				m["text"] = DoSelfKeyError
				if err := context.Loop.Send(m); err != nil {
					fmt.Println(err.Error())
				}
				return
			} else {
				key = userInfo.IftttKey
			}
		} else {
			if nickname == userInfo.Nickname {
				if len(userInfo.IftttKey) == 0 {
					m["text"] = DoSelfKeyError
					if err := context.Loop.Send(m); err != nil {
						fmt.Println(err.Error())
					}
					return
				} else {
					key = userInfo.IftttKey
				}
			} else {
				if userInfo.Others == nil {
					m["text"] = DoOtherNotExist
					if err := context.Loop.Send(m); err != nil {
						fmt.Println(err.Error())
					}
					return
				} else {
					for k, v := range userInfo.Others {
						if k == nickname {
							if len(v.IftttKey) == 0 {
								m["text"] = DoOtherKeyError
								if err := context.Loop.Send(m); err != nil {
									fmt.Println(err.Error())
								}
								return
							} else {
								key = v.IftttKey
								break
							}
						}
					}
				}
			}
		}
		//判断是否有引用
		if refer_key, exist := message["refer_key"]; exist {
			if refer_key != nil {
				// fmt.Println(message)
				referKey := refer_key.(string)
				//获取 引用的具体内容
				referMessage, err := GetMessageInfo(message["vchannel_id"].(string), referKey)
				if err != nil {
					m["text"] = "❎ 获取引用消息错误 : " + err.Error()
					if err := context.Loop.Send(m); err != nil {
						fmt.Println(err.Error())
					}
					return
				}
				fmt.Println(referMessage)
			}

		}
		//ifttt 请求
		if err := config.IftttPost(trigger, key, v1, v2, v3); err != nil {
			m["text"] = "❎ Ifttt 请求错误 : " + err.Error()
			if err := context.Loop.Send(m); err != nil {
				fmt.Println(err.Error())
			}
			return
		}
	}

	//
	m["text"] = "✅  触发操作成功"
	if err := context.Loop.Send(m); err != nil {
		fmt.Println(err.Error())
	}
}
