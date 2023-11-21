package process

import (
	"context"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const (
	pkKey  = "PK"
	pkKey2 = "pk"
)

type PK struct {
}

func (P PK) Prefix() []string { return []string{pkKey, pkKey2} }

func (P PK) ShortHelp() string {
	return "对某两个选手的成绩进行PK, PK-{细项} {选手1:ID/名称}vs{选手2:ID/名称}"
}

func (P PK) Help() string {
	return `PK
1. PK: PK {选手1: ID/名称} vs {选手2：ID/名称}
2. PK单类型PK: PK-{细项} {选手1: ID/名称} vs {选手2：ID/名称}
3. 指定项目PK: 指定项目PK， PK[项目] {选手1: ID/名称} vs {选手2：ID/名称}`
}

func (P PK) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	//out := inMessage.CopyOut()
	//
	//msg := ReplaceAll(inMessage.Content, "", pkKey2, pkKey, "-")
	//msg = ReplaceAll(msg, ",", "，", ".")
	//msg = ReplaceAll(msg, "vs", "VS", "Vs", "vS")
	//
	//// 在首个空格处切割
	//data := strings.SplitN(msg, " ", 1)
	//if len(data) != 2 {
	//	return EventHandler(out.AddSprintf(P.Help()))
	//}
	//header := data[0]
	//footer := data[1]
	//
	//// 获取两个角色
	//strings.Split(footer, "vs")
	//
	//// pk[333,222,444] 1 vs 2
	//if strings.Contains(msg, "[") || strings.Contains(msg, "【") {
	//
	//}
	//
	//// pk-WCA项目 1 vs 2
	//if len(msg) > 0 {
	//
	//}
	//
	//// pk 1 vs 2

	return nil
}

func (P PK) withProjectPK(p1, p2 model.Player, pjs []model.Project) string {
	var out string
	return out
}
