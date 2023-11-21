package process

import (
	"context"
	"fmt"
	"sync"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

var _ Process = &Help{}

const (
	helpKey1 = "帮助"
	helpKey2 = "help"
)

type Help struct {
	once sync.Once

	keyMap map[string]Process
}

func (c *Help) Prefix() []string { return []string{helpKey1, helpKey2} }

func (c *Help) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	c.once.Do(func() { c.keyMap = PrefixMap(List...) })
	out := inMessage.CopyOut()

	msg := ReplaceAll(inMessage.Content, "", helpKey1, helpKey2, " ", "-")
	if len(msg) > 0 {
		p, err := CheckPrefix(msg, c.keyMap)
		if err != nil {
			return EventHandler(out.AddError(err))
		}

		return EventHandler(out.AddSprintf(p.Help()))
	}
	return EventHandler(out.AddSprintf(c.Help()))
}
func (c *Help) ShortHelp() string {
	return "获取帮助信息, 帮助-{指令}可获取详细帮助"
}

func (c *Help) Help() string {
	var out = ""
	for idx, p := range List {
		pfH := p.Prefix()[0]
		pfV := p.Prefix()[1:]
		if len(pfV) == 0 {
			out += fmt.Sprintf("%d. %s %s\n", idx+1, pfH, p.ShortHelp())
			continue
		}
		out += fmt.Sprintf("%d. %s%+v: %s\n", idx+1, pfH, pfV, p.ShortHelp())
	}
	return out
}
