package cq_http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my_cubing_robot/pkg/bot/process_tools"
	"github.com/guojia99/my_cubing_robot/pkg/process"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type CQHttpClient struct {
	process_tools.ProcessClient

	conf Config

	api *gin.Engine
}

func NewCQHttpClient(conf Config, db *gorm.DB) *CQHttpClient {
	return &CQHttpClient{
		conf: conf,
		ProcessClient: process_tools.ProcessClient{
			Ctx:      nil,
			Db:       db,
			Core:     core.NewCore(db, false, time.Second),
			InputCh:  make(chan process.InMessage, 512),
			OutputCh: make(chan interface{}, 128),
		},
	}
}

func (c *CQHttpClient) Run(ctx context.Context) error {
	c.Ctx = ctx
	c.api = gin.Default()
	c.SendMsgFn = c.sendMsgFn()

	c.api.NoRoute(c.route)

	go c.DoProcessLoop()

	return c.api.Run(fmt.Sprintf("127.0.0.1:%d", c.conf.Port))
}

func (c *CQHttpClient) route(ctx *gin.Context) {
	var r CqInMessage
	_ = ctx.Bind(&r)
	ctx.JSON(http.StatusOK, gin.H{})

	if r.MetaEventType == "heartbeat" {
		return
	}
	if r.MessageType != MessageTypeGroup {
		return
	}
	// 判断是不是艾特自己
	if r.SelfID == r.UserId {
		return
	}

	msg := process.InMessage{
		MessageID:        fmt.Sprintf("%d", r.MessageId),
		ChannelID:        "",
		GroupID:          fmt.Sprintf("%d", r.GroupID),
		UserID:           fmt.Sprintf("%d", r.UserId),
		UserName:         r.Sender.NickName,
		Content:          "",
		NotPrefixContent: "",
	}

	for _, v := range r.Message {
		switch v.Type {
		case "text":
			msg.Content += v.Data.Text
		case "at":
			// 解析raw message 是否有艾特自己
			// [CQ:at,qq=2854216320]
			atMe := fmt.Sprintf("[CQ:at,qq=%d]", r.SelfID)
			if v.Data.QQ == r.SelfID ||
				(strings.Contains(r.RawMessage, "[CQ:at,qq=") && !strings.Contains(r.RawMessage, atMe)) {
				return
			}
		}
	}

	msg.Content = strings.TrimPrefix(msg.Content, " ")

	select {
	case <-ctx.Done():
		return
	case c.InputCh <- msg:
		return
	}
}
