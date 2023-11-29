package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	mRankKey  = "月排名"
	mRankKey2 = "m_rank"
)

type MRank struct {
}

func (M MRank) CheckPrefix(in string) bool {
	return false
}

func (M MRank) Prefix() []string { return []string{mRankKey, mRankKey2} }

func (M MRank) ShortHelp() string {
	return "可获取某月排名信息,默认今年, 月排名[月份]-{项目} {数量}, 月份格式：1~12或202301~202312"
}

func (M MRank) Help() string {
	return `月排名
1. mrank-{项目}: 获取本月的月排行
2. mrank[月份]-{项目}: 获取某月的月排行, 格式1~12或202301~202312
3. mrank-{项目} {数量}, 获取项目的排名长度, 最多30`
}

func (M MRank) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", mRankKey, mRankKey2, "-")
	msg = strings.ReplaceAll(msg, "【", "[")
	msg = strings.ReplaceAll(msg, "】", "]")

	// 时间
	var first, last = getFirstDateAndLastDate(time.Now())
	if strings.Contains(msg, "[") {
		sp := strings.Split(msg[1:], "]")
		if len(sp) != 2 {
			return EventHandler(out.AddSprintf("输入格式错误"))
		}
		fmt.Println(sp[0])
		first, last = getFirstDateAndLastDate(getStringTime(sp[0]))
		msg = sp[1]
	}

	// 数量
	msg = strings.TrimLeftFunc(msg, unicode.IsSpace)
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

	// 项目
	project := getProjectString(split[0])
	if len(project) == 0 {
		return nil
	}

	// 获取
	bestsM, avgsM := core.GetBestScoreByTimes(first, last)
	bests, avgs := bestsM[project], avgsM[project]
	if len(bests) == 0 {
		return EventHandler(out.AddSprintf("该月份此项目无人参加, 欢迎积极参赛"))
	}

	// 渲染
	out.AddSprintf("%s %d-%d月度排行Top%d\n", project.Cn(), first.Year(), first.Month(), limit)
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
	//out.AddSprintf("详情请查看: http://www.mycube.club/statistics/best?tabs=best_all&cubes=%s", project)
	return EventHandler(out)
}

func getStringTime(input string) time.Time {
	year := time.Now().Year()
	month := time.Now().Month()
	day := 1

	// 情况1: 只有年
	if len(input) <= 2 {
		monthInt, _ := strconv.Atoi(input)
		month = time.Month(monthInt)
	} else if len(input) == 4 {
		year, _ = strconv.Atoi(input)
		month = 1
	} else if len(input) == 6 {
		year, _ = strconv.Atoi(input[:4])
		monthInt, _ := strconv.Atoi(input[4:6])
		month = time.Month(monthInt)
	} else {
		month = time.Now().Month()
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func getFirstDateAndLastDate(baseTime time.Time) (time.Time, time.Time) {
	firstDateTime := baseTime.AddDate(0, 0, -baseTime.Day()+1)
	firstDateZeroTime := time.Date(firstDateTime.Year(), firstDateTime.Month(), firstDateTime.Day(), 0, 0, 0, 0, firstDateTime.Location())

	lastDateTime := baseTime.AddDate(0, 1, -baseTime.Day())
	lastDateZeroTime := time.Date(lastDateTime.Year(), lastDateTime.Month(), lastDateTime.Day(), 0, 0, 0, 0, firstDateTime.Location())

	return firstDateZeroTime, lastDateZeroTime
}
