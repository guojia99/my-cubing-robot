package process

import (
	"context"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

type TestFile struct {
}

func (t TestFile) CheckPrefix(in string) bool {
	return false
}

func (t TestFile) Prefix() []string { return []string{"测试"} }

func (t TestFile) ShortHelp() string { return "测试用的" }

func (t TestFile) Help() string { return "测试用的" }

func (t TestFile) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	if strings.Contains(inMessage.Content, "image") {
		out.AddImage("https://mycube.club/x-file/robot_image/contest_37_tab_nav_all_score_table.png")
		out.AddSprintf("测试图文")
		return EventHandler(out)
	}
	return nil

}
