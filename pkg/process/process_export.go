package process

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	core "github.com/guojia99/my-cubing-core"
	exports "github.com/guojia99/my-cubing-core/export"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

func init() {
	_ = os.MkdirAll("/data/x-file/robot_image", 0755)
}

type Export struct {
	sync.Mutex
}

func (e *Export) CheckPrefix(in string) bool {
	return false
}

const (
	exportKey  = "导出"
	exportKey2 = "export"
)

func (e *Export) Prefix() []string { return []string{exportKey, exportKey2} }

func (e *Export) ShortHelp() string {
	return " (管理员) 导出成绩表格或者其他数据, 导出-{比赛ID}, 该功能群主无法使用"
}

func (e *Export) Help() string {
	return `导出功能:
1. 导出-{比赛ID}: 导出该场比赛的表格数据`
	//2. 导出选手-{选手ID}: 导出选手成绩, 没有ID代表所有选手的数据.
	//3. 导出记录： 导出记录历史和当前的记录`
}

func (e *Export) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	if !e.Mutex.TryLock() {
		return EventHandler(out.AddSprintf("导出中，请勿重复发送指令"))
	}
	defer e.Unlock()

	msg := ReplaceAll(inMessage.Content, "", exportKey, exportKey2, "-")
	if len(msg) == 0 {
		return nil
	}

	numbers := utils.GetNumbers(msg)
	if len(numbers) == 0 {
		return EventHandler(out.AddSprintf("无法导出 `%s`", msg))
	}
	id := int(numbers[0])

	file := path.Join("/data/x-file/robot_image", fmt.Sprintf("contest_%d_value.xlsx", id))
	imageUrl := fmt.Sprintf("https://mycube.club/x-file/robot_image/contest_%d_value.xlsx", id)
	if err := exports.ExportContestScoreXlsx(core, uint(id), file); err != nil {
		return EventHandler(out.AddError(err))
	}

	out.AddSprintf("导出成功: %s", imageUrl)
	return EventHandler(out)
}
