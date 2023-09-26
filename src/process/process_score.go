package process

import (
	"fmt"
	"strings"

	"github.com/guojia99/my-cubing/src/core"
	"github.com/guojia99/my-cubing/src/core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

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

type ProjectDetail struct {
	Idx     int    `json:"idx" table:"序号"`
	Player1 string `json:"player1" table:"选手"`
	Best    string `json:"score1" table:"单次"`
	Avg     string `json:"score2" table:"平均"`
	Player2 string `json:"player2" table:"选手"`
}

const ScoreKey = "rank-"

func Rank(db *gorm.DB, core core.Core, inMessage string) (outMessage string) {
	inMessage = strings.ReplaceAll(inMessage, " ", "")
	if !strings.HasPrefix(inMessage, ScoreKey) {
		return ""
	}

	pj := getProjectString(strings.ReplaceAll(inMessage, ScoreKey, ""))
	if pj == "" {
		return ""
	}

	allBest, allAvg := core.GetAllPlayerBestScore()
	bests, ok := allBest[pj]
	if !ok || len(bests) == 0 {
		return "该项目无人参加, 欢迎积极参赛"
	}

	avgs, ok := allAvg[pj]
	outMessage = fmt.Sprintf("%s Top10\n", pj.Cn())
	for idx, best := range bests {
		if idx >= 10 {
			break
		}

		outMessage += fmt.Sprintf("%d、%s %s", idx+1, best.PlayerName, utils.TimeParser(best, false))
		if ok && idx < len(avgs) {
			outMessage += fmt.Sprintf(" || %s %s", utils.TimeParser(avgs[idx], true), avgs[idx].PlayerName)
		}
		
		outMessage += "\n"
	}

	outMessage += fmt.Sprintf("详情请查看: http://mycube.club/statistics/best?tabs=best_all&cubes=%s", pj)

	return outMessage
}
