package process

import (
	"context"
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
	"strings"
	"sync"
)

const (
	randomPoKey1 = "选择"
)

type RandomPo struct {
	once sync.Once
}

func (c *RandomPo) ShortHelp() string {
	return ""
}

func (c *RandomPo) Help() string {
	return ""
}

func (c *RandomPo) IsGroup() bool {
	return false
}

func (c *RandomPo) CheckPrefix(in string) bool {
	return false
}

func (c *RandomPo) Prefix() []string { return []string{randomPoKey1} }

func (c *RandomPo) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", randomPoKey1)

	sl := strings.Split(msg, " ")
	if len(sl) == 0 {
		return EventHandler(out.AddSprintf("空空如也"))
	}

	var newList = []string{}
	for _, v := range sl {
		v = strings.TrimLeft(v, " ")
		if len(v) == 0 {
			continue
		}
		if v == " " {
			continue
		}
		newList = append(newList, v)
	}

	newList = shuffledCopy(newList, false)
	return EventHandler(out.AddSprintf(newList[0]))
}
