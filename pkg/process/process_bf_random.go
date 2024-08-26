package process

import (
	"context"
	"sync"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

var _ Process = &BfRandom{}

const (
	bfRandomKey1 = "盲拧"
	bfRandomKey2 = "bf"
)

type BfRandom struct {
	once   sync.Once
	keyMap map[string]Process
}

func (c *BfRandom) IsGroup() bool {
	return false
}

func (c *BfRandom) CheckPrefix(in string) bool { return false }
func (c *BfRandom) Prefix() []string           { return []string{bfRandomKey1, bfRandomKey2} }

func (c *BfRandom) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	//out := inMessage.CopyOut()
	//msg := ReplaceAll(inMessage.Content, "", bfRandomKey1, bfRandomKey2)
	//return EventHandler(out)
	return nil
}

func (c *BfRandom) ShortHelp() string {
	return "获取帮助信息, 帮助-{指令}可获取详细帮助"
}

func (c *BfRandom) Help() string {
	return `

`
}
