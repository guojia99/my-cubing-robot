package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const ContestKey = "contest"
const ContestKey2 = "contest-"

const ContestSubKeySor = "-sor"
const ContestSubKeyRank = "-rank"
const ContestSubKeyRecord = "-record"

func Contest(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string, outImage string) {

	if !strings.HasPrefix(inMessage, ContestKey) {
		return
	}

	var contestId int
	if strings.Contains(inMessage, ContestKey2) {
		inMessage = strings.ReplaceAll(inMessage, ContestKey2, "")
		numbers := _getNumbers(inMessage)
		if len(numbers) > 0 {
			contestId = int(numbers[0])
		}
	}

	var contest model.Contest
	if contestId != 0 {
		if err := db.Where("id = ?", contestId).First(&contest).Error; err != nil {
			return "无法找到该比赛", ""
		}
	} else {
		if err := db.Order("created_at DESC").Order("id DESC").First(&contest).Error; err != nil {
			return "获取比赛错误", ""
		}
	}

	var contestFile = path.Join("/tmp", fmt.Sprintf("contest_%d_tab_nav_all_score_table.png", contest.ID))
	var url = fmt.Sprintf("https://mycube.club/contest?id=%d&contest_tab=tab_nav_all_score_table", contest.ID)

	outMessage = fmt.Sprintf("比赛: %s\n详情请查看 %s\n", contest.Name, url)
	status, err := os.Stat(contestFile)
	if err == nil && time.Since(status.ModTime()) < time.Minute*30 {
		return outMessage, contestFile
	}

	if err = exec.Command("python3", "/usr/local/bin/load_mycube_image.py", "--image", contestFile, "--url", url).Run(); err != nil {
		log.Println("error", err)
		return outMessage, ""
	}

	return outMessage, contestFile
}
