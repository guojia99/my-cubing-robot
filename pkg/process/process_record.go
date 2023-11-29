package process

import (
	"context"

	coreModel "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/utils"
)

const (
	recordKey  = "记录"
	recordKey2 = "record"
)

type Record struct {
}

func (r Record) CheckPrefix(in string) bool {
	return false
}

func (r Record) Prefix() []string { return []string{recordKey, recordKey2} }

func (r Record) ShortHelp() string {
	return "获取赛季记录"
}

func (r Record) Help() string {
	return `记录
1. 记录： 直接获取所有项目最佳成绩
2. 记录-{比赛ID}: 获取某场比赛打破的记录`
}

func (r Record) Do(ctx context.Context, db *gorm.DB, core coreModel.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	msg := ReplaceAll(inMessage.Content, "", recordKey2, recordKey, "-")
	if len(msg) == 0 {
		return r.sendBestRecord(ctx, db, core, inMessage, EventHandler)
	}
	return r.sendContestRecord(ctx, db, core, inMessage, EventHandler)
}

func (r Record) sendContestRecord(ctx context.Context, db *gorm.DB, core coreModel.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	numbers := utils.GetNumbers(inMessage.Content)
	if len(numbers) == 0 {
		return nil
	}

	var cache = make(map[model.Project][]coreModel.RecordMessage)
	var contest model.Contest
	for _, val := range core.GetContestRecord(uint(numbers[0])) {
		contest = val.Contest
		if _, ok := cache[val.Score.Project]; !ok {
			cache[val.Score.Project] = make([]coreModel.RecordMessage, 0)
		}
		cache[val.Score.Project] = append(cache[val.Score.Project], val)
	}

	if len(cache) == 0 {
		out.AddSprintf("该比赛不存在、未产生记录或未结束")
		return EventHandler(out)
	}

	out.AddSprintf("%s 比赛记录\n", contest.Name)
	for _, pj := range model.AllProjectRoute() {
		records, ok := cache[pj]
		if !ok {
			continue
		}

		var best *coreModel.RecordMessage
		var avg *coreModel.RecordMessage

		for _, record := range records {
			if record.Record.RType == model.RecordBySingle {
				best = &record
			} else {
				avg = &record
			}
		}

		if best != nil && avg != nil {
			out.AddSprintf(
				"%s %s %s | %s %s\n", best.Score.Project.Cn(), utils.TimeParser(best.Score, false), best.Score.PlayerName,
				avg.Score.PlayerName, utils.TimeParser(avg.Score, true),
			)
			continue
		}

		if best != nil {
			out.AddSprintf("%s %s %s\n", best.Score.Project.Cn(), utils.TimeParser(best.Score, false), best.Score.PlayerName)
		}
		if avg != nil {
			out.AddSprintf("%s      | %s %s\n", avg.Score.Project.Cn(), utils.TimeParser(avg.Score, true), avg.Score.PlayerName)
		}
	}
	return EventHandler(out)
}

func (r Record) sendBestRecord(ctx context.Context, db *gorm.DB, core coreModel.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()

	bestM, avgM := core.GetAllProjectBestScores()

	out.AddSprintf("最佳记录\n")

	for _, class := range projectClass {
		set := false
		for _, pj := range getClassProjects(class) {
			best, ok := bestM[pj]
			if !ok {
				continue
			}
			if !set {
				out.AddSprintf("----------- %s -----------\n", class)
				set = true
			}

			out.AddSprintf("%s %s %s", pj.Cn()+":", utils.TimeParser(best, false), best.PlayerName)

			avg, ok := avgM[pj]
			if ok {
				out.AddSprintf("| %s %s", utils.TimeParser(avg, true), avg.PlayerName)
			}
			out.AddSprintf("\n")
		}
	}
	out.AddSprintf("--------------------------\n详情查看 https://mycube.club/statistics/best\n")
	return EventHandler(out)
}
