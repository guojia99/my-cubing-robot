package src

import (
	"github.com/guojia99/my_cubing_robot/src/utils"
)

type Message struct {
	Anonymous   interface{} `json:"anonymous"`
	Font        int         `json:"font"`
	GroupId     int         `json:"group_id"`
	Message     string      `json:"message"`
	MessageId   int         `json:"message_id"`
	MessageSeq  int         `json:"message_seq"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	SelfId      int64       `json:"self_id"`
	Sender      Sender      `json:"sender"`
	SubType     string      `json:"sub_type"`
	Time        int         `json:"time"`
	UserId      int         `json:"user_id"`
}

type Sender struct {
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	UserId   int    `json:"user_id"`
}

type Message2 struct {
	PostType      string    `json:"post_type"`
	MetaEventType string    `json:"meta_event_type"`
	Time          int       `json:"time"`
	SelfId        int64     `json:"self_id"`
	St            MessageST `json:"st"`
	Interval      int       `json:"interval"`
}

type MessageST struct {
	AppEnabled     bool        `json:"app_enabled"`
	AppGood        bool        `json:"app_good"`
	AppInitialized bool        `json:"app_initialized"`
	Good           bool        `json:"good"`
	Online         bool        `json:"online"`
	PluginsGood    interface{} `json:"plugins_good"`
	Stat           struct {
		PacketReceived  int `json:"packet_received"`
		PacketSent      int `json:"packet_sent"`
		PacketLost      int `json:"packet_lost"`
		MessageReceived int `json:"message_received"`
		MessageSent     int `json:"message_sent"`
		DisconnectTimes int `json:"disconnect_times"`
		LostTimes       int `json:"lost_times"`
		LastMessageTime int `json:"last_message_time"`
	} `json:"stat"`
}

type SendMessageModel struct {
	GroupId int    `json:"group_id"`
	Message string `json:"message"`
}

func SendMessage(groupId int, message string) error {
	if message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}
	_, err := utils.HTTPRequest("POST", "http://127.0.0.1:5700/send_group_msg", nil, nil, SendMessageModel{
		GroupId: groupId,
		Message: message,
	})
	return err
}
