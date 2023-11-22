package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	recordKey  = "记录"
	recordKey2 = "record"
)

type Record struct {
}

func (r Record) Prefix() []string { return []string{recordKey, recordKey2} }

func (r Record) ShortHelp() string {
	return "获取赛季记录"
}

func (r Record) Help() string {
	return ""
}

func (r Record) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	return nil
}
