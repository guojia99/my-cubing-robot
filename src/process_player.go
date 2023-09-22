package src

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guojia99/my-cubing/src/core/model"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

const processPlayerKey = "*选手"

func countChineseCharacters(text string) int {
	count := 0
	for _, char := range text {
		if utf8.RuneLen(char) >= 1 { // 判断字符是否占用多个字节（多个码点）
			count++
		}
	}
	return count
}

var wcaRoute = func() []model.Project {
	out := model.WCAProjectRoute()

	n := len(out)
	for i := 0; i < n-1; i++ {
		sw := false
		for j := 0; j < n-i-1; j++ {
			if countChineseCharacters(out[j].Cn()) > countChineseCharacters(out[j+1].Cn()) {
				out[j], out[j+1] = out[j+1], out[j]
				sw = true
			}
		}
		if !sw {
			break
		}
	}
	return out
}()

var xCubeRoute = func() []model.Project {
	out := model.XCubeProjectRoute()

	n := len(out)
	for i := 0; i < n-1; i++ {
		sw := false
		for j := 0; j < n-i-1; j++ {
			if countChineseCharacters(out[j].Cn()) > countChineseCharacters(out[j+1].Cn()) {
				out[j], out[j+1] = out[j+1], out[j]
				sw = true
			}
		}
		if !sw {
			break
		}
	}
	return out
}()

func (c *Client) processPlayer(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processPlayerKey, "")
	in = strings.ReplaceAll(in, " ", "")

	if len(in) == 0 {
		return SendMessage(msg.GroupId, fmt.Sprintf("请输入: %s {选手ID或名称}", processPlayerKey))
	}

	id, err := strconv.Atoi(in)
	var player model.Player
	if err == nil {
		err = c.db.Where("id = ?", id).First(&player).Error
	} else {
		err = c.db.Where("name = ?", in).First(&player).Error
	}

	if player.ID == 0 || err != nil {
		var players []model.Player
		if err = c.db.Where("name like ?", fmt.Sprintf("%%%s%%", in)).Find(&players).Error; err != nil || len(players) == 0 {
			return SendMessage(msg.GroupId, "查询不到选手")
		}

		out := "请选择指定的选手进行查询\n"
		for _, val := range players {
			out += fmt.Sprintf("%d、%s\n", val.ID, val.Name)
		}

		if len(players) == 1 {
			player = players[0]
		} else {
			return SendMessage(msg.GroupId, out)
		}
	}

	bestAll, avgAll := c.core.GetPlayerBestScore(player.ID)
	out := player.Name + "\n"
	out += "----------- 个人主页 -----------\n"
	out += fmt.Sprintf("http://mycube.club/player?id=%d\n\n", player.ID)

	// 成绩渲染函数
	score := func(pj model.Project) string {
		best, ok := bestAll[pj]
		avg, ok2 := avgAll[pj]
		if !ok && !ok2 {
			return ""
		}

		out := ""
		if pj.RouteType() == model.RouteTypeRepeatedly {
			out += fmt.Sprintf("%s: %.0f / %.0f %s", pj.Cn(), best.Score.Result1, best.Score.Result2, utils.TimeParser(best.Score, false))
		} else {
			out += fmt.Sprintf("%s: %s", pj.Cn(), utils.TimeParser(best.Score, false))
			if ok2 {
				out += fmt.Sprintf(" | %s", utils.TimeParser(avg.Score, true))
			}
		}
		return out + "\n"
	}

	// 渲染成绩
	wcaOut := ""
	for _, pj := range wcaRoute {
		if val := score(pj); out != "" {
			wcaOut += val
		}
	}
	if len(wcaOut) > 0 {
		wcaOut = "----------- 官方项目 -----------\n" + wcaOut
	}

	xcubeOut := ""
	for _, pj := range xCubeRoute {
		if val := score(pj); out != "" {
			xcubeOut += val
		}
	}
	if len(xcubeOut) > 0 {
		xcubeOut = "----------- 趣味项目 -----------\n" + xcubeOut
	}

	return SendMessage(msg.GroupId, out+wcaOut+xcubeOut)
}
