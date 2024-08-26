package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	groupIDKey  = "群ID"
	groupIDKey2 = "group_id"
)

type GetGroupID struct {
}

func (g GetGroupID) IsGroup() bool {
	return true
}

func (g GetGroupID) CheckPrefix(in string) bool { return false }

func (g GetGroupID) Prefix() []string { return []string{groupIDKey, groupIDKey2} }

func (g GetGroupID) ShortHelp() string { return "获取群ID" }

func (g GetGroupID) Help() string { return g.ShortHelp() }

func (g GetGroupID) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	out.AddSprintf("群ID: %s", inMessage.GroupID)
	return EventHandler(out)
}
