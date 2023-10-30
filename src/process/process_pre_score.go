package process

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	coreModel "github.com/guojia99/my-cubing-core"

	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

const AddPreScoreKey = "录入"

/*
1, 快速录入（多个录入用  / 分开）：
*录入  333 1.1,1.2,1:03.10,DNF,DNS
*录入   333 1.1,1.2,1:03.10,DNF,DNS  / 444 1.2,1.3,1.2,1.2,1.43 / ....

2, 详细录入:
*录入-31  333(2）1.2(1,2,3),1.1(2,4),1.2,1:03.10,DNF,DNS
*录入-{比赛ID}  {项目名}({轮次})  {成绩1}(判罚列表...), {成绩2}(判罚列表...), {成绩3}(判罚列表...), ...
*/

func AddPreScore(db *gorm.DB, core coreModel.Core, inMessage string, qq string) (outMessage string) {
	if !strings.HasPrefix(inMessage, AddPreScoreKey) {
		return ""
	}
	inMessage = strings.ReplaceAll(inMessage, AddPreScoreKey, "")
	inMessage = strings.ReplaceAll(inMessage, "\n", "")
	inMessage = strings.ReplaceAll(inMessage, "：", ":")
	inMessage = strings.ReplaceAll(inMessage, "，", ",")
	inMessage = strings.ReplaceAll(inMessage, "\\", "/")
	inMessage = strings.ReplaceAll(inMessage, "。", ".")

	if len(inMessage) <= 1 {
		return "无效的输入"
	}
	var contest model.Contest
	if inMessage[0] == '-' {
		if num := _getNumbers(inMessage[1:]); len(num) != 0 {
			id := int(num[0])
			if err := db.Where("id = ?", id).Where("is_end = ?", false).First(&contest).Error; err == nil {
				inMessage = strings.ReplaceAll(inMessage[1:], fmt.Sprintf("%d", id), "")
				return _simpleAddPreScore(db, core, contest, inMessage, qq)
			}
		}
		return "输入无效的比赛或找不到该比赛"
	}

	if err := db.Where("is_end = ?", false).Where("name like ?", fmt.Sprintf("%%%s%%", "群赛")).First(&contest).Error; err == nil {
		return _simpleAddPreScore(db, core, contest, inMessage, qq)
	}

	return "没有开启比赛或找不到该比赛"
}

func _simpleAddPreScore(db *gorm.DB, core coreModel.Core, contest model.Contest, inMessage string, qq string) (outMessage string) {
	// 1, 快速录入（多个录入用  / 分开）：
	//*录入  333 1.1,1.2,1:03.10,DNF,DNS
	//*录入   333 1.1,1.2,1:03.10,DNF,DNS  / 444 1.2,1.3,1.2,1.2,1.43 / ....
	var playerUser model.PlayerUser
	if err := db.Where("qq = ?", qq).First(&playerUser).Error; err != nil {
		return fmt.Sprintf("`%s` 未登记，请联系浩浩登记", qq)
	}

	var player model.Player
	_ = db.Where("id = ?", playerUser.PlayerID).First(&player)

	preScores, err := _preScoresParser(db, contest, inMessage)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	var out = fmt.Sprintf("比赛：%s\n", contest.Name)
	for _, val := range preScores {
		val.PlayerID = playerUser.PlayerID
		val.ContestID = contest.ID
		val.Source = "QQ-robot"
		val.Recorder = qq

		if err = core.AddPreScore(val); err != nil {
			out += fmt.Sprintf("%s %s 录入失败： %s\n", player.Name, val.Project.Cn(), err)
		} else {
			out += fmt.Sprintf("%s %s 录入成功！\n", player.Name, val.Project.Cn())
		}
	}
	return out
}

func _preScoresParser(db *gorm.DB, contest model.Contest, inMessage string) ([]coreModel.AddPreScoreRequest, error) {
	scores := strings.Split(inMessage, "/")
	fmt.Println(scores)
	if len(scores) == 0 {
		return nil, errors.New("请输入正确的录入:\n 如：*录入 333 1.1,1.2,1:03.10,DNF,DNS")
	}

	var preScores []coreModel.AddPreScoreRequest
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
		var preScore = coreModel.AddPreScoreRequest{
			AddScoreRequest: coreModel.AddScoreRequest{
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

func _getNumbers(in string) []float64 {
	re := regexp.MustCompile("(-?\\d+)(\\.\\d+)?")
	numbers := re.FindAllString(in, -1)

	var out []float64
	for _, num := range numbers {
		f, err := strconv.ParseFloat(num, 64)
		if err == nil {
			out = append(out, f)
		}
	}
	return out
}
