package process

import (
	"strings"

	"github.com/guojia99/my-cubing/src/core"
	"gorm.io/gorm"
)

const ContestKey = "contest"
const ContestKey2 = "contest-"

const ContestSubKeySor = "-sor"
const ContestSubKeyRank = "-rank"
const ContestSubKeyRecord = "-record"

func Contest(db *gorm.DB, core core.Core, inMessage string) (outMessage string) {

	if strings.Contains(inMessage, ContestKey) || strings.Contains(inMessage, ContestKey2) {
		return "暂未实现"
	}

	return ""
}