package process

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	coreUtils "github.com/guojia99/my-cubing-core/utils"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	preEnterKey  = "录入"
	preEnterKey2 = "enter"
)

var _ Process = &PreEnter{}

type PreEnter struct {
	once sync.Once

	mp map[string]Process
}

func (c *PreEnter) CheckPrefix(in string) bool {
	c.once.Do(
		func() {
			c.mp = make(map[string]Process)
			for _, pj := range model.AllProjectRoute() {
				c.mp[pj.Cn()] = c
				c.mp[string(pj)] = c
			}
		},
	)

	in = strings.TrimSpace(in)
	key, _, err := CheckPrefix(in, c.mp)
	if err != nil {
		return false
	}
	in = strings.ReplaceAll(in, key, "")

	if len(utils.GetNumbers(in)) > 0 {
		return true
	}

	return false
}

func (c *PreEnter) Prefix() []string { return []string{preEnterKey, preEnterKey2} }
func (c *PreEnter) ShortHelp() string {
	return "(需注册) 可录入某场比赛项目成绩, 录入 {项目1} {成绩列表1} / {项目2} {成绩列表2} ..."
}
func (c *PreEnter) Help() string {
	return `录入
* 你可以使用 [登记] 功能进行登记你的qq帐号

0. 如存在注册信息异常无法录入的情况下, 请使用第四个方式进行录入.
1. 快速录入:  录入 {项目} {成绩...}
2. 多项目:  录入 {项目} {成绩...} / {项目2} {成绩2}
3. 指定比赛ID: 录入- {比赛ID} {项目} {成绩...}
4. 为指定选手录入成绩: 录入[ID/名称]-{比赛ID} {项目} {成绩...} / {项目2} {成绩2}`
}

func (c *PreEnter) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	// 过滤无效数据
	msg := ReplaceAll(inMessage.Content, "", append(c.Prefix(), "(", ")", "\n")...)
	msg = strings.ReplaceAll(msg, "：", ":")
	msg = strings.ReplaceAll(msg, "，", ",")
	msg = strings.ReplaceAll(msg, "\\", "/")
	msg = strings.ReplaceAll(msg, "。", ".")
	msg = ReplaceAll(msg, "[", "【", "〔", "〈", "［")
	msg = ReplaceAll(msg, "]", "】", "〕", "〉", "］")

	// 获取是否存在比赛
	var (
		contest    model.Contest
		playerUser model.PlayerUser
		player     model.Player
		err        error
	)

	// 解析用户
	if len(msg) > 1 && msg[0] == '[' {
		split := strings.Split(msg[1:], "]")
		if len(split) != 2 {
			return EventHandler(out.AddSprintf("格式错误"))
		}
		player, err = getPlayer(db, split[0])
		if err != nil {
			return EventHandler(out.AddSprintf("找不到对应玩家"))
		}
		msg = split[1]
	} else {
		if err = db.Where("qq = ?", inMessage.UserID).Or("qq_bot_uni_id = ?", inMessage.UserID).First(&playerUser).Error; err != nil {
			return EventHandler(out.AddSprintf("`%s` 未登记，请联系浩浩登记", inMessage.UserID))
		}
		var _ = db.Where("id = ?", playerUser.PlayerID).First(&player)
	}

	// 解析比赛
	if len(msg) > 1 && msg[0] == '-' {
		if num := utils.GetNumbers(msg[1:]); len(num) != 0 {
			id := int(num[0])
			err = db.Where("id = ?", id).Where("is_end = ?", false).First(&contest).Error
			msg = strings.Replace(msg[1:], fmt.Sprintf("%d", id), "", 1)
		}
	} else {
		err = db.Where("is_end = ?", false).Where("name like ?", fmt.Sprintf("%%%s%%", "群赛")).First(&contest).Error
	}

	if err != nil {
		return EventHandler(out.AddError(errors.New("输入无效的比赛或找不到该比赛")))
	}

	return EventHandler(out.AddSprintf(_simpleAddPreScore(db, core, player, contest, msg, inMessage.UserID)))
}

func _simpleAddPreScore(db *gorm.DB, core core.Core, player model.Player, contest model.Contest, inMessage string, qq string) (outMessage string) {
	// 快速录入（多个录入用  / 分开）：
	//*录入  333 1.1,1.2,1:03.10,DNF,DNS
	//*录入   333 1.1,1.2,1:03.10,DNF,DNS  / 444 1.2,1.3,1.2,1.2,1.43 / ....
	//

	preScores, err := _preScoresParser(db, contest, inMessage)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	var out = fmt.Sprintf("比赛：%s\n", contest.Name)

	bests, avgs := core.GetAllProjectBestScores()
	sBests, sAvgs := core.GetPlayerBestScore(player.ID)

	for _, val := range preScores {
		val.PlayerID = player.ID
		val.ContestID = contest.ID
		val.Source = "QQ-robot"
		val.Recorder = qq

		if err = core.AddPreScore(val); err != nil {
			out += fmt.Sprintf("%s %s 录入失败： %s\n", player.Name, val.Project.Cn(), err)
			continue
		}

		score := model.Score{Project: val.Project}
		score.SetResult(val.Result, model.ScorePenalty{})
		out += fmt.Sprintf("%s %s (%s / %s) 录入成功！\n", player.Name, val.Project.Cn(), coreUtils.BestOrAvgParser(score, false), coreUtils.BestOrAvgParser(score, true))

		// 刷新记录成绩提示
		best, bestOk := bests[val.Project]
		avg, avgOk := avgs[val.Project]
		if !bestOk && !avgOk && !score.DBest() && !score.DAvg() {
			out += fmt.Sprintf("(该成绩是第一个成功有单次和平均有效成绩)\n")
			continue
		} else if !bestOk && !score.DBest() {
			out += fmt.Sprintf("(该成绩是第一个成功有单次成绩)\n")
		} else if !avgOk && !score.DAvg() {
			out += fmt.Sprintf("(该成绩是第一个成功有平均成绩)\n")
		}
		if bestOk && avgOk && score.IsBestScore(best) && score.IsBestAvgScore(avg) { // 双刷提示
			out += fmt.Sprintf("(该成绩双刷了`%s | %s` 的历史最佳成绩 (%s / %s))\n", best.PlayerName, avg.PlayerName, coreUtils.BestOrAvgParser(best, false), coreUtils.BestOrAvgParser(avg, true))
			continue
		} else if bestOk && score.IsBestScore(best) { // 刷了单次
			out += fmt.Sprintf("(该成绩打破了`%s` 的历史单次最佳成绩 %s)\n", best.PlayerName, coreUtils.BestOrAvgParser(best, false))
			continue
		} else if avgOk && score.IsBestAvgScore(avg) { // 刷了平均
			if val.Project.RouteType() != model.RouteType1rounds && val.Project.RouteType() != model.RouteTypeRepeatedly {
				out += fmt.Sprintf("(该成绩打破了`%s` 的历史平均最佳成绩 %s)\n", best.PlayerName, coreUtils.BestOrAvgParser(avg, true))
			}
			continue
		}

		// 刷新pb
		sBest, sBestOk := sBests[val.Project]
		sAvg, sAvgOk := sAvgs[val.Project]

		if sBestOk && sAvgOk && score.IsBestScore(sBest.Score) && score.IsBestAvgScore(sAvg.Score) { // 双刷提示
			out += fmt.Sprintf("(该成绩双刷了自己的历史最佳成绩 (%s / %s))\n", coreUtils.BestOrAvgParser(sBest.Score, false), coreUtils.BestOrAvgParser(sAvg.Score, true))
		} else if sBestOk && score.IsBestScore(sBest.Score) { // 刷了单次
			out += fmt.Sprintf("(该成绩打破了自己的历史单次最佳成绩 %s)\n", coreUtils.BestOrAvgParser(sBest.Score, false))
		} else if sAvgOk && score.IsBestAvgScore(sAvg.Score) { // 刷了平均
			if val.Project.RouteType() != model.RouteType1rounds && val.Project.RouteType() != model.RouteTypeRepeatedly {
				out += fmt.Sprintf("(该成绩打破了自己的历史平均最佳成绩 %s)\n", coreUtils.BestOrAvgParser(sAvg.Score, true))
			}
		}
	}
	return out
}

func _preScoresParser(db *gorm.DB, contest model.Contest, inMessage string) ([]core.AddPreScoreRequest, error) {
	scores := strings.Split(inMessage, "/")
	fmt.Println(scores)
	if len(scores) == 0 {
		return nil, errors.New("请输入正确的录入:\n 如：*录入 333 1.1,1.2,1:03.10,DNF,DNS")
	}

	var preScores []core.AddPreScoreRequest
	for _, score := range scores {
		pj := _getProject(score)
		if pj == "" {
			return nil, fmt.Errorf("`%s`有不存在的项目", score)
		}

		// todo 解析 轮次
		var roundNumber = 1
		var round model.Round
		if err := db.Where("contest_id = ?", contest.ID).Where("project = ?", pj).Where("number = ?", roundNumber).Where("is_start = ?", true).First(&round).Error; err != nil {
			return nil, fmt.Errorf("`%s` 项目 轮次`%d` 不存在或者未开启该项目", pj, roundNumber)
		}

		// 移除所有成绩无关内容
		cache := strings.ReplaceAll(score, string(pj), "")
		cache = strings.ReplaceAll(cache, pj.Cn(), "")

		// 解析成绩分隔断
		var ss []string
		if strings.Contains(cache, ",") {
			cache = strings.ReplaceAll(cache, " ", "")
			ss = strings.Split(cache, ",")
		} else {
			ss = strings.Split(cache, " ")
		}
		var newSs []string
		for _, val := range ss {
			if len(val) > 0 {
				newSs = append(newSs, val)
			}
		}
		ss = newSs
		if len(ss) == 0 {
			return nil, fmt.Errorf("`%s` 无法执行无成绩的内容", score)
		}

		// 数据处理
		var preScore = core.AddPreScoreRequest{
			AddScoreRequest: core.AddScoreRequest{
				Project: pj,
				RoundId: round.ID,
				Result:  []float64{},
				Penalty: model.ScorePenalty{},
			},
		}
		// 提取成绩：
		// 1:03.10, DNF, DNS
		// 1:03.10(1,2), DNF, DNS
		for _, s := range ss {
			// todo 成绩解析 带penalty的
			preScore.Result = append(preScore.Result, utils.ParserTimeToSeconds(s))
		}

		preScores = append(preScores, preScore)
	}

	return preScores, nil
}

var pjMap = func() map[string]model.Project {
	var out = make(map[string]model.Project)

	for _, pj := range model.AllProjectRoute() {
		out[pj.Cn()] = pj
		out[string(pj)] = pj
	}
	return out
}()

func _getProject(in string) model.Project {
	// 333 1.1,1.2,1:03.10,DNF,DNS
	for {
		if len(in) == 0 {
			return ""
		}
		if in[0] == ' ' {
			in = in[1:]
			continue
		}
		break
	}

	split := strings.Split(in, " ")
	if len(split) == 0 {
		return ""
	}

	key := split[0]
	val, _ := pjMap[key]
	return val
}
