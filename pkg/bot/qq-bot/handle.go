package qq_bot

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	jsoniter "github.com/json-iterator/go"

	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/event"
	"github.com/guojia99/my_cubing_robot/pkg/process"
	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

func (q *QQBotClient) checkAnyMessage(input string, at bool) (result string, err error) {
	// 移除前面所有空格
	result = strings.TrimLeftFunc(input, unicode.IsSpace)

	// 频道at的情况
	pattern := `<@!\d+> `
	result = regexp.MustCompile(pattern).ReplaceAllString(result, "")

	err = errors.New("not prefix")
	if utils.HasPrefixList(result, "/", "*", ".") {
		result = result[1:]
		err = nil
	}

	// 移除多余空格
	result = strings.Join(strings.Fields(result), " ")

	// at 强制通过
	if at {
		err = nil
	}
	return result, err
}

func (q *QQBotClient) _atGroupMessageEventHandler() event.GroupAtMessageEventHandler {
	return func(event *WSPayload, data *WSGroupATMessageData) (err error) {
		if data.Content, err = q.checkAnyMessage(data.Content, true); err == nil {
			q.inputCh <- anyMessageToMessage(data)
		}
		return
	}
}

func (q *QQBotClient) _groupMessageEventHandler() event.GroupMessageEventHandler {
	return func(event *WSPayload, data *WSGroupMessageData) (err error) {
		if data.Content, err = q.checkAnyMessage(data.Content, false); err == nil {
			q.inputCh <- anyMessageToMessage(data)
		}
		return
	}
}

func (q *QQBotClient) _atMessageEventHandler() event.ATMessageEventHandler {
	return func(event *WSPayload, data *WSATMessageData) (err error) {
		if data.Content, err = q.checkAnyMessage(data.Content, true); err == nil {
			q.inputCh <- anyMessageToMessage(data)
		}
		return
	}
}

func (q *QQBotClient) _messageEventHandler() event.MessageEventHandler {
	return func(event *WSPayload, data *WSMessageData) (err error) {
		if data.Content, err = q.checkAnyMessage(data.Content, false); err == nil {
			q.inputCh <- anyMessageToMessage(data)
		}
		return
	}
}

func anyMessageToMessage(in interface{}) (out process.InMessage) {
	switch in.(type) {
	case *WSMessageData, *WSATMessageData:
		var msg Message
		data, _ := jsoniter.Marshal(in)
		_ = jsoniter.Unmarshal(data, &msg)

		out = process.InMessage{
			MessageID: msg.ID,
			ChannelID: msg.ChannelID,
			GroupID:   msg.GuildID,
			UserID:    msg.Author.ID,
			Content:   msg.Content,
		}
	case *WSGroupMessageData, *WSGroupATMessageData:
		var msg GroupMessage
		data, _ := jsoniter.Marshal(in)
		_ = jsoniter.Unmarshal(data, &msg)

		out = process.InMessage{
			MessageID: msg.MsgId,
			ChannelID: "",
			GroupID:   msg.GroupId,
			UserID:    msg.Author.UserId,
			UserName:  msg.Author.UserOpenId,
			Content:   msg.Content,
		}
	}
	return
}
