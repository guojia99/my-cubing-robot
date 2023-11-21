package process

import (
	"context"
	"fmt"
	"slices"

	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"
)

const (
	projectClassKey  = "项目分类"
	projectClassKey2 = "project-class"
	projectClassKey3 = "分类"
)

type ProjectClass struct {
}

func (p ProjectClass) Prefix() []string {
	return []string{projectClassKey, projectClassKey2, projectClassKey3}
}

func (p ProjectClass) ShortHelp() string {
	return "获取项目的分类，及下列其他使用的功能所需的{细项}"
}

func (p ProjectClass) Help() string {
	return `项目分类:
1. 项目分类： 获取所有分类
2. 项目分类-{细项}: 获取该分类下的所有项目`
}

func (p ProjectClass) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", projectClassKey, projectClassKey2, projectClassKey3, "-", " ")
	if len(msg) == 0 {
		return EventHandler(out.AddSprintf(p.allClassMsg()))
	}
	return EventHandler(out.AddSprintf(p.classMsg(msg)))
}

func (p ProjectClass) classMsg(in string) string {
	var out string
	for _, val := range projectClass {
		if string(val) == in {
			out += fmt.Sprintf("%s\n", val)

			idx := 1
			for _, pj := range model.AllProjectRoute() {
				if slices.Contains(pj.Class(), string(val)) {
					out += fmt.Sprintf("%d. %s %s\n", idx, pj.Cn(), pj)
					idx++
				}
			}

			return out
		}
	}
	return p.allClassMsg()
}

func (p ProjectClass) allClassMsg() string {
	var out string
	for idx, val := range projectClass {
		out += fmt.Sprintf("%d. %s\n", idx+1, val)
	}
	return out
}

var projectClass = []model.ProjectClass{
	model.ProjectClassWCA,
	model.ProjectClassXCube,
	model.ProjectClassXCubeBF,
	model.ProjectClassXCubeOH,
	model.ProjectClassXCubeFm,
	model.ProjectClassXCubeRelay,
	model.ProjectClassNotCube,
	model.ProjectClassDigit,
	model.ProjectClassSuperHigh,
}
