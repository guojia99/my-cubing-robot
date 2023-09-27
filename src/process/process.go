package process

import (
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

type (
	ProcessFn func(db *gorm.DB, core core.Core, inMessage string) (outMessage string)
)

var ProcessDict = []ProcessFn{
	Rank,
	Sor,
	Contest,
	PK,
	Help,

	// player查询永久垫底
	Player,
}
