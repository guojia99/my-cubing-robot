package process

import (
	"context"
	"fmt"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const (
	notPlayKey1 = "未参与"
	notPlayKey2 = "not-play"
)

var _ Process = &NotPlay{}

type NotPlay struct {
}

func (c *NotPlay) CheckPrefix(in string) bool {
	return false
}

func (c *NotPlay) Prefix() []string { return []string{notPlayKey1, notPlayKey2} }

func (c *NotPlay) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	// get detail
	var playerUser model.PlayerUser
	if err := db.Where("qq = ?", inMessage.UserID).Or("qq_bot_uni_id = ?", inMessage.UserID).First(&playerUser).Error; err != nil {
		return EventHandler(out.AddSprintf("`%s` 未登记，请联系浩浩登记", inMessage.UserID))
	}
	var player model.Player
	_ = db.Where("id = ?", playerUser.PlayerID).First(&player)

	var contest model.Contest
	if err := db.Where("is_end = ?", false).Where("name like ?", fmt.Sprintf("%%%s%%", "群赛")).First(&contest).Error; err != nil {
		return EventHandler(out.AddSprintf("没有开启最新的群赛，请联系浩浩开启"))
	}
	contestDetail, err := core.GetContest(contest.ID)
	if err != nil {
		return EventHandler(out.AddSprintf("没有开启最新的群赛，请联系浩浩开启"))
	}

	playerContest, _ := core.GetScoreByPlayerContest(playerUser.PlayerID, contest.ID)
	var cache = make(map[model.Project]struct{})
	for _, val := range playerContest {
		cache[val.Project] = struct{}{}
	}

	// class
	var class = model.ProjectClassWCA
	for _, val := range projectClass {
		if strings.Contains(inMessage.Content, string(val)) {
			class = string(val)
		}
	}
	out.AddSprintf("未参与%s项目\n", class)

	for _, val := range contestDetail.Rounds {
		ok := false
		for _, cs := range val.Project.Class() {
			if cs == class {
				ok = true
			}
		}
		if !ok {
			continue
		}
		if _, ok = cache[val.Project]; !ok {
			out.AddSprintf("%s(%s) ", val.Project.Cn(), val.Project)
		}
	}

	return EventHandler(out)
}
func (c *NotPlay) ShortHelp() string {
	return "(需注册) 可获取选手未参加的项目列表, 未参与-{细项}可获取某个详细的未参与列表,默认WCA."
}
func (c *NotPlay) Help() string {
	return fmt.Sprintf(
		`未参与
0. 未参与是指本期比赛未参与的项目列表, 需要群主录入账号到系统中
1. 未参与 将使用最新一场未结束的群赛作为未参与的群赛内容.
2. 未参与-{分类} 可指定某些类别的项目
3. 分类: %+v
	
`, projectClass,
	)
}
