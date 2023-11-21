package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	sorXKey = "sor-x"
)

type SorX struct {
}

func (s SorX) Prefix() []string {
	return []string{sorXKey}
}

func (s SorX) ShortHelp() string {
	return "成绩综合实力分数, SOR-X-{细项}, 默认WCA"
}

func (s SorX) Help() string {
	return ""
}

func (s SorX) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	return nil
}
