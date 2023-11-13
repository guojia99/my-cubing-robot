package process

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fogleman/gg"
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/gorm"
)

const NvHaoHaoKey = "女装"

func autoImage() string {
	// 创建一个300x300像素的图像
	rand.Seed(time.Now().UnixNano())
	const width = 300
	const height = 300
	dc := gg.NewContext(width, height)

	// 随机生成图像内容
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// 随机颜色
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			dc.SetRGB(float64(r)/255, float64(g)/255, float64(b)/255)
			dc.SetPixel(x, y)
		}
	}

	out := fmt.Sprintf("/tmp/%d.png", time.Now().UnixNano())
	if err := dc.SavePNG(out); err != nil {
		return ""
	}
	return out
}

func NvHaoHao(db *gorm.DB, core core.Core, inMessage string, qq string) (outMessage string, outImage string) {
	if !strings.HasPrefix(inMessage, NvHaoHaoKey) {
		return
	}
	return "女装浩赶紧女装", ""
	//return "女装浩赶紧女装", autoImage()
}
