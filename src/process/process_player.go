package process

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

var (
	playerCache []string
	lastUpdate  time.Time
)

func getPlayerList(db *gorm.DB) []string {
	if time.Now().Sub(lastUpdate) < time.Minute*5 {
		return playerCache
	}
	var p []string
	// 使用 Pluck 提取 name 字段到切片中
	db.Model(&model.Player{}).Pluck("name", &p)
	playerCache = p
	lastUpdate = time.Now()
	return p
}

const PlayerKey = "选手"

func Player(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string) {
	if !strings.HasPrefix(inMessage, PlayerKey) {
		return ""
	}

	inMessage = strings.ReplaceAll(inMessage, PlayerKey, "")
	if len(inMessage) > 20 || len(inMessage) == 0 {
		return
	}

	// 查数据
	var player model.Player

	id, err := strconv.Atoi(inMessage)
	if err == nil {
		db.Where("id = ?", id).First(&player)
	}

	if player.ID == 0 {
		for _, name := range getPlayerList(db) {
			if name == inMessage {
				db.Where("name = ?", inMessage).First(&player)
				break
			}
		}
	}

	if player.ID == 0 {
		var players []model.Player
		db.Where("name like ?", fmt.Sprintf("%%%s%%", inMessage)).Find(&players)
		if len(players) >= 2 {
			outMessage += "选择指定的选手进行查询\n"
			for _, val := range players {
				outMessage += fmt.Sprintf("%d、%s\n", val.ID, val.Name)
			}
			return outMessage
		}
		if len(players) == 0 {
			outMessage = "查询不到数据\n"
			return outMessage
		}
		player = players[0]
	}

	// 查成绩
	bestAll, avgAll := core.GetPlayerBestScore(player.ID)

	outMessage = player.Name + "\n"
	outMessage += "----------- 个人主页 -----------\n"
	outMessage += fmt.Sprintf("http://www.mycube.club/player?id=%d\n", player.ID)
	// 成绩渲染函数
	score := func(pj model.Project) string {
		best, ok := bestAll[pj]
		avg, ok2 := avgAll[pj]
		if !ok && !ok2 {
			return ""
		}

		out := ""
		if pj.RouteType() == model.RouteTypeRepeatedly {
			out += fmt.Sprintf("%s %s / %s %s",
				utils.TB(pj.Cn()+":", 5),
				fmt.Sprintf("%.0f", best.Score.Result1),
				fmt.Sprintf("%.0f", best.Score.Result2),
				utils.TB(utils.TimeParser(best.Score, false), 4),
			)
		} else {
			out += fmt.Sprintf("%s %s",
				utils.TB(pj.Cn()+":", 5),
				utils.TB(utils.TimeParser(best.Score, false), 4),
			)
			if ok2 {
				out += fmt.Sprintf(" | %s", utils.TB(utils.TimeParser(avg.Score, true), 4))
			}
		}
		return out + "\n"
	}

	// 渲染成绩
	wcaOut := ""
	for _, pj := range wcaRoute {
		wcaOut += score(pj)
	}
	if len(wcaOut) > 0 {
		wcaOut = "----------- 官方项目 -----------\n" + wcaOut
	}

	xcubeOut := ""
	for _, pj := range xCubeRoute {
		xcubeOut += score(pj)
	}
	if len(xcubeOut) > 0 {
		xcubeOut = "----------- 趣味项目 -----------\n" + xcubeOut
	}
	outMessage += wcaOut + xcubeOut

	// 渲染排位
	signalSor, avgSor := core.GetPlayerSor(player.ID)
	sorOut := "----------- 排位分数 -----------\n"
	sorOut += "项目    排名 单次  || 平均 排名\n"
	for _, val := range sorKeys {
		sorOut += fmt.Sprintf("%s\t %s %s || %s %s\n",
			utils.TB(SorCn[val], 3),
			utils.TB(signalSor[val].SingleRank, 2),
			utils.TB(signalSor[val].SingleCount, 2),
			utils.TB(avgSor[val].AvgCount, 2),
			utils.TB(avgSor[val].AvgRank, 2),
		)
	}
	outMessage += sorOut

	return outMessage
}

var wcaRoute = func() []model.Project {
	out := model.WCAProjectRoute()

	n := len(out)
	for i := 0; i < n-1; i++ {
		sw := false
		for j := 0; j < n-i-1; j++ {
			if utils.CountChineseCharacters(out[j].Cn()) > utils.CountChineseCharacters(out[j+1].Cn()) {
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
			if utils.CountChineseCharacters(out[j].Cn()) > utils.CountChineseCharacters(out[j+1].Cn()) {
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
