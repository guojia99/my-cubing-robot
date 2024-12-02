package qq_bot

import (
	"context"
	"log"
	"time"

	"github.com/guojia99/my_cubing_robot/pkg/process"
)

func (q *QQBotClient) getImageInfo(group, image string) (string, error) {
	value, ok := q.imageCache.Get(group + image)
	if ok {
		return value.(string), nil
	}

	out, err := q.api.PostGroupRichMediaMessage(
		q.Ctx, group, &GroupRichMediaMessageToCreate{
			FileType:   1,
			Url:        image,
			SrvSendMsg: false,
		},
	)
	if err != nil {
		return "", err
	}

	q.imageCache.Set(out.FileInfo, group+image, time.Minute*2)

	log.Printf("with file info %s, uuid %s\n", out.FileInfo, out.FileUuid)
	return out.FileInfo, nil
}

func (q *QQBotClient) groupMsg(message *process.OutMessage) (err error) {
	msg := &GroupMessageToCreate{
		Content: message.OutContent,
		MsgID:   message.MessageID,
	}

	// 发送富媒体
	if message.Image != "" {
		msg.MsgType = 7
		msg.Media.FileInfo, err = q.getImageInfo(message.GroupID, message.Image)
		if err != nil {
			return err
		}
	}
	_, err = q.api.PostGroupMessage(q.Ctx, message.GroupID, msg)
	return err
}

func (q *QQBotClient) sendMsg(message *process.OutMessage) error {
	_, err := q.api.PostMessage(
		context.TODO(), message.ChannelID, &MessageToCreate{
			Content: message.OutContent,
			MsgID:   message.MessageID,
			Image:   message.Image,
		},
	)
	return err
}

func (q *QQBotClient) sendMsgFn() process.SendEventHandler {
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
		if q.conf.Group {
			return q.groupMsg(message)
		}
		return q.sendMsg(message)
	}
}
