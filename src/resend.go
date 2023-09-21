package src

import (
	"fmt"
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

func (c *Client) messageProcess(msg Message) error {
	switch {
	case strings.HasPrefix(msg.Message, "*name"):

	case strings.HasPrefix(msg.Message, processContestListKey):
		return c.processContestList(msg)
	case strings.HasPrefix(msg.Message, processContestKey):
		return c.processContest(msg)
	}
	return fmt.Errorf("not suppr")
}
