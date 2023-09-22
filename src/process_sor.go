package src

import (
	"fmt"
	"strings"

	"github.com/guojia99/my-cubing/src/core/model"
)

const processSorKey = "*排位"

var sorKeyMap = map[string]model.SorStatisticsKey{
	model.SorWCA:             model.SorWCA,
	"全项目":                    model.SorWCA,
	model.SorXCube:           model.SorXCube,
	"趣味":                     model.SorXCube,
	model.SorWCACubeLowLevel: model.SorWCACubeLowLevel,
	"二至五":                    model.SorWCACubeLowLevel,
	model.SorWCACubeAllLevel: model.SorWCACubeAllLevel,
	"二至七":                    model.SorWCACubeAllLevel,
	model.SorWCAAlien:        model.SorWCAAlien,
	"异形":                     model.SorWCAAlien,
	model.SorWCA333:          model.SorWCA333,
	"全三阶":                    model.SorWCA333,
	model.SorWCABf:           model.SorWCABf,
	"盲拧":                     model.SorWCABf,
}

func (c *Client) processSor(msg Message) error {
	in := strings.ReplaceAll(msg.Message, processSorKey, "")
	in = strings.ReplaceAll(in, " ", "")

	var key = model.SorWCA
	if val, ok := sorKeyMap[in]; ok {
		key = val
	}

	allBest, allAvg := c.core.GetSorScore()

	b := allBest[key]
	a := allAvg[key]

	out := fmt.Sprintf("--------------- %s ---------------\n", key)
	out += fmt.Sprintf("详情请查看 :http://mycube.club/statistics/sor?sor_tabs=%s", key)
	for idx, best := range b {
		if idx >= 10 {
			break
		}

		out += fmt.Sprintf("%d、 %s %d", idx+1, best.Player.Name, best.SingleCount)
		out += fmt.Sprintf(" || %d %s\n", a[idx].AvgCount, a[idx].Player.Name)
	}
	return SendMessage(msg.GroupId, out)
}
