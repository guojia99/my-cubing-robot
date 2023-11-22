package process

import (
	"context"
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

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	contestKey1 = "比赛"
	contestKey2 = "contest"
)

var _ Process = &Contest{}

type Contest struct {
}

func (c *Contest) Prefix() []string { return []string{contestKey1, contestKey2} }
func (c *Contest) ShortHelp() string {
	return "获取比赛信息, 比赛-{赛事ID/名称} 可获取某场比赛详细信息, 比赛列表可获取近期比赛场次"
}
func (c *Contest) Help() string {
	return `- 比赛使用方法:
1. 比赛: 比赛即可获取本期比赛成绩
2. 比赛-{num}：可以指定ID选择某场比赛的成绩
3. 比赛-{名称}: 指定某个比赛.
4. 比赛列表: 获取最近10期比赛的比赛列表
5. 比赛列表-{num}: 翻页, 代表第n页`
}

func (c *Contest) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	if strings.Contains(inMessage.Content, "列表") {
		return c.sendList(ctx, db, core, inMessage, EventHandler)
	}
	return c.sendContest(ctx, db, core, inMessage, EventHandler)
}

func (c *Contest) getContest(ctx context.Context, db *gorm.DB, core core.Core, msg string) (model.Contest, error) {
	var contestId int
	numbers := utils.GetNumbers(msg)
	if len(numbers) > 0 {
		contestId = int(numbers[0])
	}

	var contest model.Contest
	if contestId != 0 {
		return contest, db.Where("id = ?", contestId).First(&contest).Error
	}
	msg = ReplaceAll(msg, "", c.Prefix()...)
	if len(msg) > 0 {
		err := db.Model(&model.Contest{}).Where("name like ?", fmt.Sprintf("%%%s%%", msg)).First(&contest).Error
		if err == nil {
			return contest, nil
		}
	}

	err := db.Order("created_at DESC").First(&contest).Error
	return contest, err
}

func (c *Contest) sendContest(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	contest, err := c.getContest(ctx, db, core, inMessage.Content)
	if err != nil {
		return EventHandler(out.AddError("找不到比赛"))
	}

	// todo 上传到cos
	var contestFile = path.Join("/tmp", fmt.Sprintf("contest_%d_tab_nav_all_score_table.png", contest.ID))
	var url = fmt.Sprintf("https://mycube.club/contest?id=%d&contest_tab=tab_nav_all_score_table", contest.ID)
	//var url = ""
	out.AddSprintf("比赛: %s\n详情请查看 %s\n", contest.Name, url)

	status, err := os.Stat(contestFile)
	if err == nil && time.Since(status.ModTime()) < time.Minute*30 {
		return EventHandler(out.AddSprintf("【缓存图片】\n").AddImage(contestFile))
	}

	if err = exec.Command("python3", "/usr/local/bin/load_mycube_image.py", "--image", contestFile, "--url", url).Run(); err != nil {
		log.Println("error", err)
		out.AddSprintf("【图像生成失败】")
		return EventHandler(out)
	}
	return EventHandler(out.AddImage(contestFile))
}

func (c *Contest) sendList(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	var page = 0
	numbers := utils.GetNumbers(inMessage.Content)
	if len(numbers) > 0 {
		page = int(numbers[0])
	}

	count, contests, err := core.GetContests(page, 10, "")
	if err != nil {
		return err
	}
	out := inMessage.CopyOut()
	for _, contest := range contests {
		out.AddSprintf(
			"%d. %s (%+v)\n", contest.ID, contest.Name, func() string {
				if contest.IsEnd {
					return "结束"
				}
				return "进行中"
			}(),
		)
	}
	out.AddSprintf("本赛事共收录%d场比赛", count)
	return EventHandler(out)
}