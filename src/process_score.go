package src

import (
	"fmt"
	"strings"

	"github.com/guojia99/my-cubing/src/core/model"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

const processProjectKey = "*项目"

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

func (c *Client) processProject(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processProjectKey, "")
	in = strings.ReplaceAll(in, " ", "")

	pj := getProjectString(in)
	if pj == "" {
		return SendMessage(msg.GroupId, "请输入正确的项目")
	}

	allBest, allAvg := c.core.GetAllPlayerBestScore()

	bests, ok := allBest[pj]
	if !ok || len(bests) == 0 {
		return SendMessage(msg.GroupId, "该项目无人参加, 欢迎积极参赛")
	}

	avgs, ok := allAvg[pj]
	out := fmt.Sprintf("%s Top10\n", pj.Cn())
	for idx, best := range bests {
		if idx >= 10 {
			break
		}

		if best.Project.RouteType() == model.RouteTypeRepeatedly {
			out += fmt.Sprintf("%d、%s\t %2.0f / %2.0f %s", idx+1, best.PlayerName, best.Result1, best.Result2, utils.TimeParser(best, true))
		} else {
			out += fmt.Sprintf("%d、%s %s", idx+1, best.PlayerName, utils.TimeParser(best, false))
			if ok && idx < len(avgs) {
				out += fmt.Sprintf(" || %s %s", utils.TimeParser(avgs[idx], true), avgs[idx].PlayerName)
			}
		}
		out += "\n"
	}

	out += fmt.Sprintf("详情请查看: http://mycube.club/statistics/best?tabs=best_all&cubes=%s", pj)

	return SendMessage(msg.GroupId, out)
}
