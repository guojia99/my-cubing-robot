package process

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

var _ Process = &Random{}

type randomValue struct {
	value  []string
	num    int
	repeat bool
}

// shuffledCopy 泛型函数，用于复制并打乱切片
// allowDuplicates 参数控制是否允许重复
// 如果允许重复，将随机选择-2, -1, 0, 1, 2作为长度变化量
func shuffledCopy[T any](slice []T, allowDuplicates bool) []T {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 初始化新切片长度为原切片长度
	newLen := len(slice)

	if allowDuplicates {
		// 从 -2, -1, 0, 1, 2 中随机选择一个作为长度变化量
		deltaOptions := []int{-2, -1, 0, 1, 2}
		delta := deltaOptions[rand.Intn(len(deltaOptions))]

		// 计算新的切片长度，确保不小于 0
		newLen += delta
		if newLen < 0 {
			newLen = 0
		}
	}

	// 创建新切片
	newSlice := make([]T, newLen)

	if allowDuplicates {
		// 允许重复时，从原切片中随机选择元素填充新切片
		for i := 0; i < newLen; i++ {
			newSlice[i] = slice[rand.Intn(len(slice))]
		}
		return newSlice
	}
	// 不允许重复时，直接复制原切片并打乱顺序
	copy(newSlice, slice)

	// 打乱新切片
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(newSlice), func(i, j int) {
			newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
		},
	)
	return newSlice
}

var randomKeys = map[string]randomValue{
	"edge": {
		value: []string{"CD", "EF", "GH", "IJ", "KL", "MN", "OP", "QR", "ST", "WX", "YZ"},
		num:   2, repeat: false,
	},
	"corner": {
		value:  []string{"ABC", "DEF", "GIH", "MNW", "OPQ", "RST", "XYZ"},
		num:    2,
		repeat: false,
	},
	"xedge": {
		value: []string{"CD", "EF", "GH", "IJ", "KL", "MN", "OP", "QR", "ST", "WX", "YZ"},
		num:   2, repeat: true,
	},
	"xcorner": {
		value:  []string{"ABC", "DEF", "GIH", "MNW", "OPQ", "RST", "XYZ"},
		num:    2,
		repeat: true,
	},
	"default": {
		value: []string{
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},
		num:    2,
		repeat: false,
	},
}

const (
	randomKey1 = "随机"
	randomKey2 = "random"
)

type Random struct {
	once sync.Once

	keyMap map[string]Process
}

func (c *Random) CheckPrefix(in string) bool {
	return false
}

func (c *Random) Prefix() []string { return []string{randomKey1, randomKey2} }

func (c *Random) Do(ctx context.Context, db *gorm.DB, core core.Core, inMessage InMessage, EventHandler SendEventHandler) error {
	out := inMessage.CopyOut()
	msg := ReplaceAll(inMessage.Content, "", randomKey1, randomKey2)

	lists := randomWithMsg(msg)

	if len(lists) == 1 {
		return EventHandler(out.AddSprintf("%s", lists[0]))
	}

	for i, l := range lists {
		out.AddSprintf("%d. %s\n", i+1, l)
	}
	return EventHandler(out)
}

func randomWithMsg(msg string) []string {
	msg = strings.ReplaceAll(msg, " ", "")

	//if len(msg) == 0 || len(strings.ReplaceAll(msg, " ", "")) == 0 {
	// 	正常输出
	//}
	var val randomValue
	if len(msg) == 0 || len(strings.ReplaceAll(msg, " ", "")) == 0 || strings.Index(msg, "*") == 0 {
		val = randomKeys["default"]
	} else if strings.Contains(msg, "3bf") {
		for _, k := range []string{
			"xedge", "xcorner", "edge", "corner",
		} {
			if strings.Contains(msg, k) {
				val = randomKeys[k]
				break
			}
		}
	} else {
		val = randomValue{
			value:  strings.Split(msg, ""),
			num:    2,
			repeat: false,
		}
	}

	var num = 1
	if strings.Contains(msg, "*") {
		numStrIdx := strings.Index(msg, "*")
		var err error
		if num, err = strconv.Atoi(msg[numStrIdx+1:]); err != nil {
			num = 1
		}
	}
	if num > 100 {
		num = 100
	}

	var outs []string
	for i := 0; i < num; i++ {
		s := shuffledCopy(val.value, val.repeat)
		var l = ""
		for x, v := range s {
			if x%val.num == 0 {
				l += " "
			}
			if len(v) <= 1 {
				l += v
				continue
			}

			xx := shuffledCopy(strings.Split(v, ""), false)
			l += xx[0]
		}
		outs = append(outs, l)
	}

	return outs
}

func (c *Random) ShortHelp() string {
	return "获取帮助信息, 帮助-{指令}可获取详细帮助"
}

func (c *Random) Help() string {
	return `
1. 输出随机的字母组合
2. 随机 3bf {edge | corner | xedge | xcorner} {*n}
`
}
