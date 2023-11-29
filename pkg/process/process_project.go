package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	projectKey  = "项目"
	projectKey2 = "project"
)

var _ Process = &Project{}

type Project struct {
}

func (c *Project) CheckPrefix(in string) bool {
	return false
}

func (c *Project) Prefix() []string { return []string{projectKey, projectKey2} }

func (c *Project) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {

	return nil
}
func (c *Project) ShortHelp() string {
	return "项目-{项目}可获取该项目排名, 项目-{项目}-{数字}可指定排名长度, 最大30."
}
func (c *Project) Help() string { return "" }
