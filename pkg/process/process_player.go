package process

import (
	"context"
	"errors"
	"fmt"
	"slices"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

var _ Process = &Player{}

const (
	playerKey  = "选手"
	playerKey2 = "玩家"
	playerKey3 = "player"
)

type Player struct {
}

func (c *Player) CheckPrefix(in string) bool {
	return false
}

func (c *Player) Prefix() []string { return []string{playerKey, playerKey2, playerKey3} }

func (c *Player) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", playerKey, playerKey2, playerKey3, " ", "-")

	if len(msg) > 1 && msg[0] == '-' {
		msg = msg[1:]
	}

	// id查询
	var player model.Player
	number := utils.GetNumbers(msg)
	if len(number) > 0 && number[0] > 0 {
		id := int(number[0])
		if err := db.Where("id = ?", id).First(&player).Error; err == nil {
			return EventHandler(out.AddSprintf("查询不到玩家"))
		}
	}

	// 模糊查询
	if player.ID == 0 {
		var players []model.Player
		db.Where("name like ?", fmt.Sprintf("%%%s%%", msg)).Find(&players)

		if len(players) >= 2 {
			out.AddSprintf("选择指定的选手进行查询\n")
			for _, val := range players {
				out.AddSprintf("%d、%s\n", val.ID, val.Name)
			}
			return EventHandler(out)
		}
		if len(players) == 0 {
			return EventHandler(out.AddSprintf("查询不到玩家"))
		}
		player = players[0]
	}

	// 渲染
	bestAll, avgAll := core.GetPlayerBestScore(player.ID)

	score := func(pj model.Project) string {
		best, ok := bestAll[pj]
		avg, ok2 := avgAll[pj]
		if !ok && !ok2 {
			return ""
		}

		output := ""
		output += fmt.Sprintf("%s (%d) %s", pj.Cn()+":", best.Rank, utils.TimeParser(best.Score, false))
		if ok2 {
			output += fmt.Sprintf(" | %s (%d)", utils.TimeParser(avg.Score, true), avg.Score.Rank)
		}
		return output + "\n"
	}

	out.AddSprintf("%s\n", player.Name)
	//out.AddSprintf("----------- 个人主页 -----------\n")
	//out.AddSprintf("http://www.mycube.club/player?id=%d\n", player.ID)

	for _, class := range projectClass {
		cur := ""
		for _, pj := range model.AllProjectRoute() {
			if slices.Contains(pj.Class(), string(class)) {
				cur += score(pj)
			}
		}

		if len(cur) > 0 {
			cur = fmt.Sprintf("----------- %s -----------\n", class) + cur
		}
		out.AddSprintf(cur)
	}
	return EventHandler(out)
}

func (c *Player) ShortHelp() string {
	return "获取选手信息, 选手-{选手ID/名称}可获取选手的信息列表."
}
func (c *Player) Help() string { return c.ShortHelp() }

func getPlayer(db *gorm.DB, in string) (model.Player, error) {
	if len(in) == 0 {
		return model.Player{}, errors.New("无法查询空的选手")
	}

	var player model.Player

	number := utils.GetNumbers(in)
	if len(number) > 0 && number[0] > 0 {
		id := int(number[0])
		if err := db.Where("id = ?", id).First(&player).Error; err == nil {
			return player, nil
		}
	}

	if err := db.Where("name like ?", fmt.Sprintf("%%%s%%", in)).First(&player).Error; err == nil {
		return player, nil
	}

	return model.Player{}, fmt.Errorf("查找不到`%s`", in)
}
