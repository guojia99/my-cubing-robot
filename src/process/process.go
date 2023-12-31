package process

import (
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

type (
	ProcessFn func(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string, image string)
)

var ProcessDict = []ProcessFn{
	Rank,
	Sor,
	Contest,
	AddPreScore,
	GetNotPlayerProject,
	PK,
	NvHaoHao,
	Player,
	Projects,
	Help,
}
