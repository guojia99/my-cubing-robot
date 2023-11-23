package qq_bot

import (
	"context"
	"log"
	"slices"
	"time"

	"k8s.io/klog"

	"github.com/guojia99/my_cubing_robot/pkg/process"
)

func (q *QQBotClient) doProcessLoop() {
	mp := process.PrefixMap(q.process...)

	for {
		select {
		case <-q.ctx.Done():
			return
		case msg := <-q.inputCh:
			klog.Infof("input msg with `%s`, by `%s` send `%s`", msg.GroupID, msg.UserID, msg.Content)
			if len(q.conf.GroupList) != 0 && !slices.Contains(q.conf.GroupList, msg.GroupID) {
				klog.Warning("continue")
				continue
			}
			func() {
				ctx, cancel := context.WithCancel(q.ctx)
				defer cancel()

				prs, err := process.CheckPrefix(msg.Content, mp)
				if err != nil {
					klog.Warning(msg.Content, err.Error())
					return
				}
				if err = prs.Do(ctx, q.db, q.core, msg, q.sendMsgFn()); err != nil {
					log.Printf("[debug] do process error %s", err)
				}
			}()

		}
	}
}

func (q *QQBotClient) getImageInfo(group, image string) (string, error) {

	value, ok := q.imageCache.Get(group + image)
	if ok {
		return value.(string), nil
	}

	out, err := q.api.PostGroupRichMediaMessage(
		q.ctx, group, &GroupRichMediaMessageToCreate{
			FileType:   1,
			Url:        image,
			SrvSendMsg: false,
		},
	)
	if err != nil {
		return "", err
	}

	q.imageCache.Set(out.FileInfo, group+image, time.Minute*2)

	klog.Infof("with file info %s, uuid %s", out.FileInfo, out.FileUuid)
	return out.FileInfo, nil
}

func (q *QQBotClient) groupMsg(message *process.OutMessage) (err error) {
	// 发送富媒体
	if message.Image != "" {
		var fileInfo string
		fileInfo, err = q.getImageInfo(message.GroupID, message.Image)
		if err != nil {
			return err
		}
		_, err = q.api.PostGroupMessage(
			q.ctx, message.GroupID, &GroupMessageToCreate{
				MsgType: 7,
				Content: message.OutContent,
				MsgID:   message.MessageID,
				Media: Media{
					FileInfo: fileInfo,
				},
			},
		)
		return
	}

	_, err = q.api.PostGroupMessage(
		q.ctx, message.GroupID, &GroupMessageToCreate{
			MsgType: 1,
			Content: message.OutContent,
			MsgID:   message.MessageID,
		},
	)
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
				klog.Error(result)
			}
			if err != nil {
				klog.Error(err)
			}
		}()
		klog.Infof("send Msg `%s` | `%s`", message.OutContent, message.Image)
		if q.conf.Group {
			return q.groupMsg(message)
		}
		return q.sendMsg(message)
	}
}
