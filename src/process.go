package src

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/guojia99/my-cubing/src/core/model"
)

const processContestListKey = "*比赛列表"

func (c *Client) processContestList(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processContestListKey, "")
	in = strings.ReplaceAll(in, " ", "")

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

	return SendMessage(SendMessageModel{GroupId: msg.GroupId, Message: out})
}

const processContestKey = "*比赛"

func (c *Client) processContest(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processContestKey, "")
	in = strings.ReplaceAll(in, " ", "")

	fmt.Println(in)
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
		return SendMessage(SendMessageModel{GroupId: msg.GroupId, Message: "未查询到数据"})
	}

	return SendMessage(
		SendMessageModel{
			GroupId: msg.GroupId,
			Message: fmt.Sprintf(
				`比赛详情
---------------------
比赛: %s
简介: %s
时间: (%s - %s)
状态: %v
网址: http://mycube.club/contest?id=%d`,
				contest.Name,
				contest.Description,
				contest.StartTime.Format("20060102"), contest.EndTime.Format("20060102"),
				contest.IsEnd,
				contest.ID,
			),
		},
	)
}
