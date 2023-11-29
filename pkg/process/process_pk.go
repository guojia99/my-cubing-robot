package process

import (
	"context"
	"fmt"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const star = "★ "
const (
	pkKey  = "PK"
	pkKey2 = "pk"
)

type PK struct {
}

func (P PK) CheckPrefix(in string) bool {
	return false
}

func (P PK) Prefix() []string { return []string{pkKey, pkKey2} }

func (P PK) ShortHelp() string {
	return "对某两个选手的成绩进行PK, PK-{细项} {选手1:ID/名称}vs{选手2:ID/名称}"
}

func (P PK) Help() string {
	return `PK
1. PK: PK {选手1: ID/名称} vs {选手2：ID/名称}
2. PK单类型PK: PK[细项] {选手1: ID/名称} vs {选手2：ID/名称}
3. 指定项目PK: 指定项目PK， PK[项目] {选手1: ID/名称} vs {选手2：ID/名称}`
}

func (P PK) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	msg := ReplaceAll(inMessage.Content, "", "-")
	msg = ReplaceAll(msg, ",", "，", ".")
	msg = ReplaceAll(msg, "vs", "VS", "Vs", "vS")
	msg = ReplaceAll(msg, "", " ")

	_, cl, players := CutMsgWithFields(msg, "vs")
	if len(players) != 2 {
		return EventHandler(out.AddSprintf("格式错误"))
	}

	var classValueOrPj []string
	for _, val := range projectClass {
		classValueOrPj = append(classValueOrPj, string(val))
	}
	if len(cl) != 0 {
		classValueOrPj = []string{}
		for _, val := range strings.Split(cl, ",") {
			key := strings.TrimSpace(val)
			if pj, ok := projectMap[key]; ok {
				classValueOrPj = append(classValueOrPj, string(pj))
				continue
			}
			if pjs := getClassProjects(model.ProjectClass(key)); len(pjs) > 0 {
				classValueOrPj = append(classValueOrPj, key)
			}
		}
	}

	if len(classValueOrPj) == 0 {
		return EventHandler(out.AddSprintf("项目细项不能为空"))
	}

	player1, err1 := getPlayer(db, players[0])
	if err1 != nil {
		return EventHandler(out.AddError(err1))
	}
	player2, err2 := getPlayer(db, players[1])
	if err1 != nil || err2 != nil {
		return EventHandler(out.AddError(err2))
	}

	out.AddSprintf("%s VS %s\n", player1.Name, player2.Name)

	p1Best, p1Avg := core.GetPlayerBestScore(player1.ID)
	p2Best, p2Avg := core.GetPlayerBestScore(player2.ID)
	p1Count, p2Count := 0, 0

	for _, class := range classValueOrPj {
		var pjs = getClassProjects(model.ProjectClass(class))
		if len(pjs) == 0 {
			pjs = append(pjs, model.Project(class))
		}

		setMsg := ""
		p1C, p2C := 0, 0
		for _, pj := range pjs {
			p1B, p1Bok1 := p1Best[pj]
			p2B, p2Bok1 := p2Best[pj]
			p1A, p1Aok2 := p1Avg[pj]
			p2A, p2Aok2 := p2Avg[pj]
			if !p1Bok1 && !p2Bok1 {
				continue
			}

			setMsg += fmt.Sprintf("%s", utils.TB(pj.Cn(), 5))

			if p1Bok1 && !p2Bok1 {
				setMsg += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB("-", 5))
				p1C += 1
			} else if !p1Bok1 && p2Bok1 {
				setMsg += fmt.Sprintf("%s   || %s", utils.TB("-", 5), star+utils.TB(utils.TimeParser(p2B.Score, false), 5))
				p2C += 1
			} else if p1B.Score.Best == p2B.Score.Best {
				setMsg += fmt.Sprintf("%s || %s", utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB(utils.TimeParser(p2B.Score, false), 5))
				p1C += 1
				p2C += 1
			} else if p1B.Score.IsBestScore(p2B.Score) {
				setMsg += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB(utils.TimeParser(p2B.Score, false), 5))
				p1C += 1
			} else {
				setMsg += fmt.Sprintf("%s   || %s", utils.TB(utils.TimeParser(p1B.Score, false), 5), star+utils.TB(utils.TimeParser(p2B.Score, false), 5))
				p2C += 1
			}
			setMsg += "\n"
			if !p1Aok2 && !p2Aok2 {
				continue
			}
			setMsg += utils.TB("", 8)
			if p1Aok2 && !p2Aok2 {
				setMsg += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB("-", 5))
				p1C += 1
			} else if !p1Aok2 && p2Aok2 {
				setMsg += fmt.Sprintf("%s   || %s", utils.TB("-", 5), star+utils.TB(utils.TimeParser(p2A.Score, true), 5))
				p2C += 1
			} else if p1A.Score.Avg == p2A.Score.Avg {
				setMsg += fmt.Sprintf("%s || %s", utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB(utils.TimeParser(p2A.Score, true), 5))
				p1C += 1
				p2C += 1
			} else if p1A.Score.IsBestAvgScore(p2A.Score) {
				setMsg += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB(utils.TimeParser(p2A.Score, true), 5))
				p1C += 1
			} else {
				setMsg += fmt.Sprintf("%s   || %s", utils.TB(utils.TimeParser(p1A.Score, true), 5), star+utils.TB(utils.TimeParser(p2A.Score, true), 5))
				p2C += 1
			}
			setMsg += "\n"
		}

		p1Count += p1C
		p2Count += p2C

		if len(setMsg) > 0 {
			out.AddSprintf("-------------- %s -------------\n", class)
			out.AddSprintf(setMsg)
			if p1C == p2C {
				out.AddSprintf("%d 平局 %d\n", p1C, p2C)
			} else if p1C > p2C {
				out.AddSprintf("胜利 %s %d vs %d\n", star, p1C, p2C)
			} else {
				out.AddSprintf("%d vs %d 胜利 %s\n", p1C, p2C, star)
			}
		}
	}
	out.AddSprintf("-------------- %s -------------\n", "总分")

	if p1Count == p2Count {
		out.AddSprintf("%d 平局 %d\n", p1Count, p2Count)
	} else if p1Count > p2Count {
		out.AddSprintf("胜利 %s %d vs %d\n", star, p1Count, p2Count)
	} else {
		out.AddSprintf("%d vs %d 胜利 %s\n", p1Count, p2Count, star)
	}

	return EventHandler(out)
}
