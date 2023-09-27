package process

import (
	"fmt"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/src/utils"
)

const star = "★ "
const PkKey = "PK"

func PK(db *gorm.DB, core core.Core, inMessage string) (outMessage string) {
	if !strings.HasPrefix(strings.ToUpper(inMessage), PkKey) {
		return ""
	}

	inMessage = strings.ReplaceAll(inMessage, PkKey, "")
	inMessage = strings.ReplaceAll(inMessage, strings.ToLower(PkKey), "")

	data := strings.Split(inMessage, "vs")
	if len(data) < 2 {
		return "请输入两位选手 *PK {选手1} vs {选手2}"
	}

	name1, name2 := strings.ReplaceAll(data[0], " ", ""), strings.ReplaceAll(data[1], " ", "")

	var player1, player2 model.Player
	if db.Where("name = ?", name1).First(&player1); player1.ID == 0 {
		return fmt.Sprintf("找不到选手`%s`", name1)
	}
	if db.Where("name = ?", name2).First(&player2); player2.ID == 0 {
		return fmt.Sprintf("找不到选手`%s`", name2)
	}

	outMessage = fmt.Sprintf("%s VS %s\n", player1.Name, player2.Name)

	p1Best, p1Avg := core.GetPlayerBestScore(player1.ID)
	p2Best, p2Avg := core.GetPlayerBestScore(player2.ID)
	p1Count, p2Count := 0, 0

	for _, pj := range model.AllProjectRoute() {
		p1B, p1Bok1 := p1Best[pj]
		p2B, p2Bok1 := p2Best[pj]
		p1A, p1Aok2 := p1Avg[pj]
		p2A, p2Aok2 := p2Avg[pj]
		if !p1Bok1 && !p2Bok1 {
			continue
		}

		outMessage += fmt.Sprintf("%s", utils.TB(pj.Cn(), 5))
		if p1Bok1 && !p2Bok1 {
			outMessage += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB("-", 5))
			p1Count += 1
		} else if !p1Bok1 && p2Bok1 {
			outMessage += fmt.Sprintf("%s   || %s", utils.TB("-", 5), star+utils.TB(utils.TimeParser(p2B.Score, false), 5))
			p2Count += 1
		} else if p1B.Score.Best == p2B.Score.Best {
			outMessage += fmt.Sprintf("%s || %s", utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB(utils.TimeParser(p2B.Score, false), 5))
			p1Count += 1
			p2Count += 1
		} else if p1B.Score.IsBestScore(p2B.Score) {
			outMessage += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1B.Score, false), 5), utils.TB(utils.TimeParser(p2B.Score, false), 5))
			p1Count += 1
		} else {
			outMessage += fmt.Sprintf("%s   || %s", utils.TB(utils.TimeParser(p1B.Score, false), 5), star+utils.TB(utils.TimeParser(p2B.Score, false), 5))
			p2Count += 1
		}
		outMessage += "\n"

		if !p1Aok2 && !p2Aok2 {
			continue
		}

		outMessage += utils.TB("", 5)
		if p1Aok2 && !p2Aok2 {
			outMessage += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB("-", 5))
			p1Count += 1
		} else if !p1Aok2 && p2Aok2 {
			outMessage += fmt.Sprintf("%s   || %s", utils.TB("-", 5), star+utils.TB(utils.TimeParser(p2A.Score, true), 5))
			p2Count += 1
		} else if p1A.Score.Avg == p2A.Score.Avg {
			outMessage += fmt.Sprintf("%s || %s", utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB(utils.TimeParser(p2A.Score, true), 5))
			p1Count += 1
			p2Count += 1
		} else if p1A.Score.IsBestAvgScore(p2A.Score) {
			outMessage += fmt.Sprintf("%s || %s", star+utils.TB(utils.TimeParser(p1A.Score, true), 5), utils.TB(utils.TimeParser(p2A.Score, true), 5))
			p1Count += 1
		} else {
			outMessage += fmt.Sprintf("%s   || %s", utils.TB(utils.TimeParser(p1A.Score, true), 5), star+utils.TB(utils.TimeParser(p2A.Score, true), 5))
			p2Count += 1
		}
		outMessage += "\n"
	}

	if p1Count == p2Count {
		outMessage += fmt.Sprintf("%d 平局 %d\n", p1Count, p2Count)
	} else if p1Count > p2Count {
		outMessage += fmt.Sprintf("胜利 %s %d vs %d\n", star, p1Count, p2Count)
	} else {
		outMessage += fmt.Sprintf("%d vs %d 胜利 %s\n", p1Count, p2Count, star)
	}

	return outMessage
}
