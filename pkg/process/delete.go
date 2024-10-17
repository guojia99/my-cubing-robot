package process

import (
	"context"
	"github.com/guojia99/my-cubing-core/model"
	"github.com/guojia99/my_cubing_robot/pkg/utils"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const (
	deleteKey  = "x删除"
	deleteKey2 = "x_delete"
)

type XDelete struct {
}

func (g XDelete) IsGroup() bool {
	return true
}

func (g XDelete) CheckPrefix(in string) bool { return false }

func (g XDelete) Prefix() []string { return []string{deleteKey, deleteKey2} }

func (g XDelete) ShortHelp() string { return "删除ID" }

func (g XDelete) Help() string { return g.ShortHelp() }

func (g XDelete) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	var usr model.PlayerUser
	if err := db.First(&usr, "qq_bot_uni_id = ?", inMessage.UserID).Error; err != nil {
		return nil
	}
	out := inMessage.CopyOut()
	if usr.PlayerID != 6 {
		return EventHandler(out.AddSprintf("无权限"))
	}
	msg := ReplaceAll(inMessage.Content, "", deleteKey, deleteKey2, " ", "-")
	if len(msg) > 1 && msg[0] == '-' {
		msg = msg[1:]
	}
	var player model.Player
	number := utils.GetNumbers(msg)
	if len(number) > 0 && number[0] > 0 {
		id := int(number[0])
		if err := db.Where("id = ?", id).First(&player).Error; err == nil {
			return EventHandler(out.AddSprintf("查询不到玩家"))
		}
	}

	return EventHandler(out.AddSprintf("%s 删除成功", player.Name))
}
