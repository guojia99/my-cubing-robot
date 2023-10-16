package process

import (
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const NvHaoHaoKey = "女装"

func NvHaoHao(db *gorm.DB, core core.Core, inMessage string) (outMessage string) {
	if !strings.HasPrefix(inMessage, NvHaoHaoKey) {
		return ""
	}
	return "女装浩赶紧女装"
}
