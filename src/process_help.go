package src

import "fmt"

const processHelpKey = "*"

var helpMaps = map[string]string{
	processProjectKey:     fmt.Sprintf("选择某个项目的前十排名, `如：%s 333` | %s 三阶", processProjectKey, processProjectKey),
	processPlayerKey:      "选择某个玩家的最佳成绩",
	processContestListKey: "获取最近和某段时间的比赛列表",
	processContestKey:     "获取某场比赛的详细信息和目前排名",
	processSorKey:         "获取排位前十, 可附带如：全项目,趣味,二至五,二至七,异形,全三阶,盲拧 等指令",
	processHelpKey:        "获取使用帮助",
}

func (c *Client) processHelp(msg Message) error {
	out := "--------- 使用帮助 -----------\n"

	for idx, key := range processFnsKey {
		out += fmt.Sprintf("%d、`%s`: %s\n", idx+1, key, helpMaps[key])
	}

	out += "--------- 温馨提示 -----------\n"
	out += "1、请在使用时切勿反复刷新,遇到卡顿时可能是限流导致,机器人响应也会做一定的限制,所以不要反复发送指令\n"
	out += "2、机器人仅提供部分方便的成绩查询, 详细内容请访问官网 http://mycube.club"

	return SendMessage(msg.GroupId, out)
}
