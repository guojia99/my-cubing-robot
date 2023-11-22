package process

import (
	"context"
	"fmt"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

var List = []Process{
	&Help{},         // 帮助
	&Contest{},      // 比赛信息
	&Project{},      // 项目列表
	&ProjectClass{}, // 项目分类
	&Player{},       // 玩家
	&Rank{},         // 排名
	&MRank{},        // 月排名
	//&NotPlay{}, // 未参与
	&PreEnter{}, // 预录入
	&PK{},       // 成绩对比
	//&Sor{},          // 排名分
	//&SorX{},         // 排位分
	//&Record{},       // 记录
	//&Export{},       // 导出
}

const MaxKeyLength = 8

type (
	InMessage struct {
		// qq bot
		MessageID string `json:"messageID"`
		ChannelID string `json:"channelID"`

		// base message
		GroupID          string `json:"groupID"` // 群ID， 频道ID
		UserID           string `json:"userID"`  // 发起的用户ID
		UserName         string `json:"userName"`
		Content          string `json:"content"`
		NotPrefixContent string `json:"notPrefixContent"`
	}

	OutMessage struct {
		InMessage

		Image      string   `json:"image"`
		Files      []string `json:"files"`
		OutContent string   `json:"content"`
	}

	SendEventHandler func(*OutMessage) error

	Process interface {
		Prefix() []string
		ShortHelp() string
		Help() string
		Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error
	}
)

func (i InMessage) CopyOut() *OutMessage {
	return &OutMessage{
		InMessage:  i,
		Image:      "",
		Files:      make([]string, 0),
		OutContent: "\n",
	}
}

// AddSprintf formats according to a format specifier and returns the resulting string.
func (o *OutMessage) AddSprintf(format string, a ...any) *OutMessage {
	o.OutContent += fmt.Sprintf(format, a...)
	return o
}

func (o *OutMessage) AddError(err any) *OutMessage {
	o.AddSprintf("错误: %+v", err)
	return o
}

func (o *OutMessage) AddImage(file string) *OutMessage {
	o.Image = file
	return o
}
