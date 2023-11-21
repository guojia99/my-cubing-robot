package process

import (
	"context"
	"sync"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

type Export struct {
	sync.Mutex
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
1. 导出-{比赛ID}: 导出该场比赛的表格数据, 没有ID代表所有.
2. 导出选手-{选手ID}: 导出选手成绩, 没有ID代表所有选手的数据.
3. 导出记录： 导出记录历史和当前的记录
`
}

func (e *Export) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	if !e.Mutex.TryLock() {
		return EventHandler(out.AddSprintf("导出中，请勿重复发送指令"))
	}
	defer e.Unlock()

	return nil
}
