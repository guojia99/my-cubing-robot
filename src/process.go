package src

import (
	"fmt"
	"log"
	"strings"
)

// https://github.com/fuqiuluo/unidbg-fetch-qsign 签名服务器

func (c *Client) setMessage(msg Message) {
	if len(msg.Message) < 1 {
		return
	}
	if msg.Message[0] == '*' {
		c.ch <- msg
	}
}

var processFnsKey = []string{
	processProjectKey,
	processPlayerKey,
	processContestListKey,
	processContestKey,
	processSorKey,
	processHelpKey,
}

func (c *Client) initProcess() {
	c.processFns = map[string]func(msg Message) error{
		processProjectKey:     c.processProject,
		processPlayerKey:      c.processPlayer,
		processContestListKey: c.processContestList,
		processContestKey:     c.processContest,
		processSorKey:         c.processSor,
		processHelpKey:        c.processHelp,
	}
}

func (c *Client) messageProcess(msg Message) error {
	defer func() {
		if result := recover(); result != nil {
			log.Println(result)
		}
	}()
	for _, key := range processFnsKey {
		if strings.HasPrefix(msg.Message, key) {
			return c.processFns[key](msg)
		}
	}
	return fmt.Errorf("not suppr")
}
