package bf

import (
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type RandomCube struct {
	// 333
	Corner *orderedmap.OrderedMap[string, string] // 角
	Edge   *orderedmap.OrderedMap[string, string] // 边

	// 444 -> Corner + Wing + XCenter
	Wing    *orderedmap.OrderedMap[string, string] // 翼棱
	XCenter *orderedmap.OrderedMap[string, string] // 角心

	// 555 -> 333 + 444
	ECenter *orderedmap.OrderedMap[string, string] // 边心

	// 字典
	CornerDict  string
	EdgeDict    string
	WingDict    string
	XCenterDict string
	ECenterDict string
}

func (r *RandomCube) Init() {
	r.CornerDict = "ABC DEF GHI --- WMN OPQ RST XYZ"
	r.EdgeDict = "-- CD EF GH IJ KL MN OP QR ST WX YZ"
	r.WingDict = "-B CD EF GH IJ KL MN OP QR ST WX YZ"
	r.XCenterDict = "ABC DEF GHI -KL WMN OPQ RST XYZ"
	r.ECenterDict = "-B CD EF GH IJ KL MN OP QR ST WX YZ"

	r.Corner = r.setMapWithDict(r.CornerDict, 3)   // 三个一组
	r.Edge = r.setMapWithDict(r.EdgeDict, 2)       // 两个字符一组
	r.Wing = r.setMapWithDict(r.WingDict, 1)       // 独立一个一组
	r.XCenter = r.setMapWithDict(r.XCenterDict, 1) // 独立一个一组
	r.ECenter = r.setMapWithDict(r.ECenterDict, 1) // 独立一个一组
}

func (r *RandomCube) setMapWithDict(dict string, num int) *orderedmap.OrderedMap[string, string] {
	var out = orderedmap.New[string, string]()
	dict = strings.ReplaceAll(dict, " ", "")
	for _, val := range splitByFixedLength(dict, num) {
		out.Set(val, val)
	}
	return out
}
