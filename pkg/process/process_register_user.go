package process

import (
	"context"
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	registerUserKey = "登记"

	unregisterUserKey = "注销"
)

var _ Process = &RegisterUser{}

type RegisterUser struct {
}

func (r *RegisterUser) Prefix() []string { return []string{registerUserKey} }

func (r *RegisterUser) ShortHelp() string {
	return "依据你的QQ频道或群聊进行编码的登记, 录入功能将不需要写入对应的昵称"
}

func (r *RegisterUser) Help() string {
	return `登记
1. 登记 {准确的昵称 | ID} : 进行绑定操作
2. 登记-注销: 进行解绑操作`
}

func (r *RegisterUser) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", registerUserKey, "-", " ")
	if len(msg) == 0 {
		return EventHandler(out.AddSprintf("不允许使用空ID/空昵称"))
	}
	msg = strings.TrimSpace(msg)

	// 解绑
	if unRegister := strings.Contains(msg, unregisterUserKey); unRegister {
		var user model.PlayerUser
		if err := db.Where("qq_bot_uni_id = ?", inMessage.UserID).First(&user).Error; err != nil {
			return EventHandler(out.AddSprintf("`%s` 未登记，请联系浩浩登记", inMessage.UserID))
		}
		if user.QQBotUniID == inMessage.UserID {
			user.QQBotUniID = ""
		}
		db.Save(user)
		return EventHandler(out.AddSprintf("解绑成功"))
	}

	// 绑定
	var checkUser model.PlayerUser
	if err := db.Where("qq_bot_uni_id = ?", inMessage.UserID).First(&checkUser).Error; err == nil {
		var p model.Player
		db.Where("id = ?", checkUser.PlayerID).First(&p)
		return EventHandler(out.AddSprintf("该ID `%s`, 已经被`%s`选手绑定", inMessage.UserID, p.Name))
	}

	var player model.Player
	numbers := utils.GetNumbers(msg)
	if err := db.Where("name = ?", msg).First(&player); err != nil {
		if len(numbers) > 0 {
			db.Where("id = ?", int(numbers[0])).First(&player)
		}
	}

	if player.ID == 0 {
		return EventHandler(out.AddSprintf("查询不到选手 `%s`", msg))
	}

	var user model.PlayerUser
	if err := db.Where("player_id = ?", player.ID).First(&user).Error; err != nil {
		return EventHandler(out.AddSprintf("查询不到选手的用户信息, 请联系浩浩登记qq后重试"))
	}
	if user.QQBotUniID != "" {
		return EventHandler(out.AddSprintf("`%s`选手已绑定 ID `%s`, 请联系该ID用户进行解绑后重试", player.Name, user.QQBotUniID))
	}

	user.QQBotUniID = inMessage.UserID
	db.Save(user)
	return EventHandler(out.AddSprintf("绑定成功"))
}
