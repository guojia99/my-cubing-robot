package process

import (
	"strings"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const HelpKey = "help"
const HelpKey2 = "帮助"

func Help(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string) {

	if !(len(inMessage) == 0 || strings.Contains(inMessage, HelpKey) || strings.Contains(inMessage, HelpKey2)) {
		return ""
	}

	outMessage = `--------- 使用帮助 -----------
0、查询所用的单位均是小写英文字符， 查询时尽可能不要包含无用字符
1、选手查询: "*选手 {选手名称}"
2、单项目查询: "*rank-{项目名}"
3、排位分查询: "*sor-{排位项}" 或者 "*sor{排位项}" 
排位项目: 全项目,wca,趣味,xcube,二至五,wca2345,二至七,wca234567,异形,wca_alien,全三阶,wca333,盲拧,wca_bf
4、PK: “*PK {选手1} vs {选手2}”
5、录入：“*录入 {项目} {成绩1},{成绩2},{成绩3},{成绩4},{成绩5}”
6、未参与： “*未参与” 查看本期未参与的群赛项目， 可附加wca 或趣味两个参数， 如“*未参与wca” 
--------- 温馨提示 -----------
1、请在使用时切勿反复刷新,遇到卡顿时可能是限流导致,机器人响应也会做一定的限制,所以不要反复发送指令
2、机器人仅提供部分方便的成绩查询, 详细内容请访问官网 http://www.mycube.club
3、看到这条消息，请催浩浩女装
`

	return outMessage
}

// 5、比赛查询: "*contest-{比赛名称}", 如果添加比赛名称则输出比赛列表
//6、比赛详细内容查询: "*contest-{比赛名称}-sor-{排位项}", "*contest-{比赛名称}-rank-{项目名}", "*contest-{比赛名称}-record"
