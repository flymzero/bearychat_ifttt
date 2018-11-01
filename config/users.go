package config

type BCObject struct { //团队成员 key uid
	Uid      string                 `json:"uid"`
	Nickname string                 `json:"nickname"`
	IftttKey string                 `json:"iftttkey"`
	Others   map[string]OtherObject `json:"others"`
}

type OtherObject struct { //编外成员 key nickname
	Nickname string `json:"nickname"`
	IftttKey string `json:"iftttkey"`
}

var Users map[string]BCObject
