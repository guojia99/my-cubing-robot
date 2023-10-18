package process

import (
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

type (
	ProcessFn func(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string)
)

var ProcessDict = []ProcessFn{
	Rank,
	Sor,
	Contest,
	AddPreScore,
	PK,
	Help,
	NvHaoHao,
	Player,
	Projects,
}
