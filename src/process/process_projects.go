package process

import (
	"fmt"
	"strings"

	coreModel "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const ProjectsKey = "项目列表"

func Projects(db *gorm.DB, core coreModel.Core, inMessage string, qq string) (outMessage string) {
	//	http://www.mycube.club/projects

	if !strings.HasPrefix(inMessage, ProjectsKey) {
		return ""
	}

	var out = "项目列表\n"
	for idx, val := range model.AllProjectRoute() {
		out += fmt.Sprintf("%d. %s %s\n", idx, val, val.Cn())
		if idx > 10 {
			break
		}
	}
	out += "\n..."
	out += "详细请查询http://www.mycube.club/projects"
	return out
}
