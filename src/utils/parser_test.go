package utils

import (
	"testing"
)

func Test_formatData(t *testing.T) {
	data := `yin 4.36 | 鱿鱼 5.82
鱿鱼 4.61 | wjz 5.87
蛋糕 5.01 | 711 6.30
711 5.03 | 小叮当uhc 6.57
wjz 5.32 | 江海 6.63
小叮当uhc 5.35 | 蛋糕 6.70
几何 5.45 | 几何 6.93
江海 6.16 | vv 7.49
Tw1light 6.30 | 孤烟往事 7.73
串串香 6.38 | czj 7.92`

	formatData(data)
}
