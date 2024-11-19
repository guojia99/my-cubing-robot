package process

import (
	"context"
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"sync"
)

const (
	randomPoKey1 = "选择"
	randomPoKey2 = "选择f"
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

func (c *RandomPo) Prefix() []string { return []string{randomPoKey1, randomPoKey2} }

func (c *RandomPo) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	fString := strings.Contains(inMessage.Content, randomPoKey2)
	msg := ReplaceAll(inMessage.Content, "", randomPoKey2, randomPoKey1)

	getList := func(sl []string) []string {
		var newList []string
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
		return newList
	}

	if !fString {
		sl := strings.Split(msg, " ")
		if len(sl) == 0 {
			return EventHandler(out.AddSprintf("空空如也"))
		}
		return EventHandler(out.AddSprintf(getList(sl)[0]))
	}

	fRe := regexp.MustCompile(`\{([^}]+)\}`)
	for {
		match := fRe.FindStringSubmatch(msg)
		if match == nil {
			break
		}
		options := strings.Split(match[1], " ")
		msg = strings.Replace(msg, match[0], getList(options)[0], 1)
	}
	return EventHandler(out.AddSprintf(msg))
}
