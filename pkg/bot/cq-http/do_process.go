package cq_http

import (
	"fmt"
	"github.com/guojia99/my_cubing_robot/pkg/process"
	"github.com/guojia99/my_cubing_robot/pkg/utils"
	"log"
	"strconv"
)

func (c *CQHttpClient) sendMsg(out *process.OutMessage) error {
	gid, _ := strconv.Atoi(out.GroupID)

	msg := CQSendMessage{
		GroupId:    int64(gid),
		Message:    []Message{},
		AutoEscape: false,
	}
	msg.Message = append(msg.Message, Message{
		Data: MessageData{
			Text: out.OutContent,
		},
		Type: "text",
	})

	msg.Message = append(msg.Message, Message{
		Data: MessageData{
			File:    "file://" + out.Image,
			Type:    "show",
			SubType: 0,
		},
		Type: "image",
	})

	// todo 使用CQ码

	_, err := utils.HTTPRequest(
		"POST", fmt.Sprintf("%s:%d/send_group_msg", c.conf.Address, c.conf.SendPort), nil, nil, msg,
	)

	//data, _ := json.Marshal(msg)
	//logger.Infof("[Robot][CQ] 发送消息 %s", string(data))

	return err
}

func (c *CQHttpClient) sendMsgFn() process.SendEventHandler {
	return func(message *process.OutMessage) (err error) {
		defer func() {
			if result := recover(); result != nil {
				log.Print(result)
			}
			if err != nil {
				log.Print(err)
			}
		}()
		log.Printf("send Msg `%s` | `%s`\n", message.OutContent, message.Image)
		//if c.conf.Group {
		//	return q.groupMsg(message)
		//}
		return c.sendMsg(message)
	}
}
