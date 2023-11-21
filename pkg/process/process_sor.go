package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	sorKey = "sor"
)

var _ Process = &Sor{}

type Sor struct {
}

func (c *Sor) Prefix() []string { return []string{sorKey} }

func (c *Sor) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	return nil
}
func (c *Sor) ShortHelp() string { return "成绩排名总和分数汇总, SOR-{细项}, 默认WCA" }
func (c *Sor) Help() string      { return "" }
