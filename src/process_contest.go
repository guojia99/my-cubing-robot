package src

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/guojia99/my-cubing/src/core/model"
)

const processContestListKey = "*比赛列表"

func (c *Client) processContestList(msg Message) error {
	//in := strings.ReplaceAll(msg.Message, processContestListKey, "")
	//in = strings.ReplaceAll(in, " ", "")

	// todo 换页选择
	var contests []model.Contest
	if err := c.db.Limit(10).Order("created_at DESC").Find(&contests).Error; err != nil {
		return err
	}
	out := "比赛列表[序号|名称]\n"
	for _, val := range contests {
		out += fmt.Sprintf("%d.%s\n", val.ID, val.Name)
	}
	out += "----------------- \n"
	out += "回复 `*比赛 +{编号}` 即可获取比赛详情"

	return SendMessage(msg.GroupId, out)
}

const processContestKey = "*比赛"

func (c *Client) processContest(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processContestKey, "")
	in = strings.ReplaceAll(in, " ", "")

	var contest model.Contest
	if len(in) == 0 {
		c.db.Order("created_at DESC").First(&contest)
	} else {
		id, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		c.db.Where("id = ?", id).First(&contest)
	}

	if contest.ID == 0 {
		return SendMessage(msg.GroupId, "未查询到数据")
	}

	// 基础信息
	out := fmt.Sprintf(
		"比赛详情\n---------- 基础信息 -----------\n比赛: %s\n简介: %s\n时间: (%s - %s)\n状态: %v\n网址: http://mycube.club/contest?id=%d\n",
		contest.Name,
		contest.Description,
		contest.StartTime.Format("20060102"), contest.EndTime.Format("20060102"),
		contest.IsEnd,
		contest.ID,
	)

	// sor
	best, avg := c.core.GetSorScoreByContest(contest.ID)

	wcaSorOut := ""
	wcaBest, ok := best[model.SorWCA]
	if ok {
		wcaAvg, wcaAvgOk := avg[model.SorWCA]
		for idx, val := range wcaBest {
			if idx >= 10 {
				break
			}
			wcaSorOut += fmt.Sprintf("%d、%s %d", idx+1, val.Player.Name, val.SingleCount)
			if wcaAvgOk && len(wcaAvg) > idx {
				wcaSorOut += fmt.Sprintf("  || %d %s", wcaAvg[idx].AvgCount, wcaAvg[idx].Player.Name)
			}
			wcaSorOut += "\n"
		}
	}
	if len(wcaSorOut) >= 0 {
		wcaSorOut = "---------- 标准排位 -----------\n" + wcaSorOut
	}

	xCubeSorOut := ""
	xCubeBest, ok := best[model.SorXCube]
	if ok {
		xCubeAvg, xCubeAvgOk := avg[model.SorXCube]
		for idx, val := range xCubeBest {
			if idx >= 10 {
				break
			}
			xCubeSorOut += fmt.Sprintf("%d、%s %d", idx+1, val.Player.Name, val.SingleCount)
			if xCubeAvgOk && len(xCubeAvg) > idx {
				xCubeSorOut += fmt.Sprintf("  || %d %s", xCubeAvg[idx].AvgCount, xCubeAvg[idx].Player.Name)
			}
			xCubeSorOut += "\n"
		}
	}
	if len(xCubeSorOut) > 0 {
		xCubeSorOut = "---------- 趣味排位 -----------\n" + xCubeSorOut
	}

	return SendMessage(msg.GroupId, out+wcaSorOut+xCubeSorOut)
}
