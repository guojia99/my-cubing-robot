package process

import (
	"fmt"
	"strings"

	"github.com/guojia99/my-cubing/src/core"
	"github.com/guojia99/my-cubing/src/core/model"
	"gorm.io/gorm"
)

var sorKeys = []string{
	model.SorWCA,
	model.SorXCube,
	model.SorWCACubeLowLevel,
	model.SorWCACubeAllLevel,
	model.SorWCAAlien,
	model.SorWCA333,
	model.SorWCABf,
}

var SorCn = map[string]string{
	model.SorWCA:             "全项目",
	model.SorXCube:           "趣味",
	model.SorWCACubeLowLevel: "二至五",
	model.SorWCACubeAllLevel: "二至七",
	model.SorWCAAlien:        "异形",
	model.SorWCA333:          "全三阶",
	model.SorWCABf:           "盲拧",
}

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

const SorKey = "sor-"
const SorKey2 = "sor"

func Sor(db *gorm.DB, core core.Core, inMessage string) (outMessage string) {
	if !strings.Contains(inMessage, SorKey) && !strings.Contains(inMessage, SorKey2) {
		return ""
	}

	in := strings.ReplaceAll(inMessage, SorKey, "")
	in = strings.ReplaceAll(inMessage, SorKey2, "")
	in = strings.ReplaceAll(in, " ", "")

	var key = model.SorWCA
	if val, ok := sorKeyMap[in]; ok {
		key = val
	}

	allBest, allAvg := core.GetSorScore()

	b := allBest[key]
	a := allAvg[key]

	outMessage = fmt.Sprintf("--------- Sor %s ----------\n", key)
	outMessage += fmt.Sprintf("详情请查看 :http://mycube.club/statistics/sor?sor_tabs=%s\n", key)
	for idx, best := range b {
		if idx >= 10 {
			break
		}

		outMessage += fmt.Sprintf("%d、 %s %d", idx+1, best.Player.Name, best.SingleCount)
		outMessage += fmt.Sprintf(" || %d %s\n", a[idx].AvgCount, a[idx].Player.Name)
	}

	return outMessage
}
