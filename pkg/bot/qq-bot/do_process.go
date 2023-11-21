package qq_bot

import (
	"context"
	"log"
	"os"
	"slices"

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

func (q *QQBotClient) sendMsgFn() process.SendEventHandler {
	return func(message *process.OutMessage) (err error) {
		if q.conf.Group {
			_, err = q.api.PostGroupMessage(
				context.TODO(), message.GroupID, &GroupMessageToCreate{
					Content: message.OutContent,
					MsgID:   message.MessageID,
				},
			)
			if err != nil {
				klog.Error(err)
			}

			if message.Image != "" {
				klog.Infof("image %s", message.Image)
				data, err2 := os.ReadFile(message.Image)
				if err2 != nil {
					return err2
				}
				_, err = q.api.PostGroupRichMediaMessage(
					context.TODO(), message.GroupID, &GroupRichMediaMessageToCreate{
						FileType:   1,
						Url:        "",
						SrvSendMsg: false,
						FileData:   data,
					},
				)
			}

		} else {
			_, err = q.api.PostMessage(
				context.TODO(), message.ChannelID, &MessageToCreate{
					Content: message.OutContent,
					MsgID:   message.MessageID,
					Image:   message.Image,
				},
			)
		}

		klog.Infof("send Msg `%s`", message.OutContent)
		return err
	}
}
