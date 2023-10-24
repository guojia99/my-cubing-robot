package process

import (
	"fmt"
	"strings"

	coreModel "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const NotPlayerProject = "未参与"

func GetNotPlayerProject(db *gorm.DB, core coreModel.Core, inMessage string, qq string) (outMessage string) {
	if !strings.HasPrefix(inMessage, NotPlayerProject) {
		return ""
	}

	var playerUser model.PlayerUser
	if err := db.Where("qq = ?", qq).First(&playerUser).Error; err != nil {
		return fmt.Sprintf("`%s` 未登记，请联系浩浩登记", qq)
	}

	var player model.Player
	_ = db.Where("id = ?", playerUser.PlayerID).First(&player)

	var contest model.Contest
	if err := db.Where("is_end = ?", false).Where("name like ?", fmt.Sprintf("%%%s%%", "群赛")).First(&contest).Error; err != nil {
		return "没有开启最新的群赛，请联系浩浩开启"
	}

	playerContest, _ := core.GetScoreByPlayerContest(playerUser.PlayerID, contest.ID)
	var cache map[model.Project]struct{}
	for _, val := range playerContest {
		cache[val.Project] = struct{}{}
	}

	wca := strings.Contains(inMessage, "wca")
	xcube := strings.Contains(inMessage, "趣味")

	var out = ""
	count := 0
	for _, val := range model.AllProjectItem() {
		if count > 15 {
			out += "..."
			break
		}
		if !val.IsWca && wca {
			continue
		}
		if val.IsWca && xcube {
			continue
		}
		if _, ok := cache[val.Project]; !ok {
			out += fmt.Sprintf("%s(%s) ", val.Project.Cn(), val.Project)
			count += 1
		}
	}
	if out == "" {
		return fmt.Sprintf("%s 已参与所有项目", player.Name)
	}

	return fmt.Sprintf("%s 未参与项目如下： %s", player.Name, out)
}
