package process

import (
	"context"
	"strconv"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	rankKey  = "排名"
	rankKey2 = "rank"
)

type Rank struct {
}

func (r Rank) Prefix() []string {
	return []string{rankKey, rankKey2}
}

func (r Rank) ShortHelp() string {
	return "可获取排名信息, 排名-{项目}"
}

func (r Rank) Help() string {
	return `排名
1. rank-{项目} : 获取项目的排名
2. rank-{项目} {数字}: 获取项目的排名长度, 最多30`
}

func (r Rank) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", "-", rankKey, rankKey2)

	split := strings.Split(msg, " ")
	var number = "10"
	if len(split) >= 2 {
		number = split[1]
	}

	var limit = 10
	if i, err := strconv.ParseInt(number, 10, 64); err == nil {
		limit = int(i)
	}
	if limit > 30 {
		limit = 30
	}

	project := getProjectString(split[0])
	if len(project) == 0 {
		return nil
	}

	bests, avgs := core.GetBestScoreByProject(project)
	if len(bests) == 0 {
		return EventHandler(out.AddSprintf("该项目无人参加, 欢迎积极参赛"))
	}

	out.AddSprintf("%s Top%d\n", project.Cn(), limit)
	for idx, best := range bests {
		if idx >= limit {
			break
		}
		out.AddSprintf("%d、%s %s", idx+1, best.PlayerName, utils.TimeParser(best, false))
		if idx < len(avgs) {
			out.AddSprintf(" || %s %s", utils.TimeParser(avgs[idx], true), avgs[idx].PlayerName)
		}
		out.AddSprintf("\n")
	}
	out.AddSprintf("详情请查看: http://www.mycube.club/statistics/best?tabs=best_all&cubes=%s", project)
	return EventHandler(out)
}

func getProjectString(in string) model.Project {
	for _, val := range model.AllProjectRoute() {
		if string(val) == in {
			return val
		}
		if val.Cn() == in {
			return val
		}
	}
	return ""
}
