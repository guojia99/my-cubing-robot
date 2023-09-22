package src

import (
	"fmt"
	"testing"

	"github.com/guojia99/my-cubing/src/core/model"
)

func TestWcaRoute(t *testing.T) {
	fmt.Println(wcaRoute)
	model.WCAProjectRoute()
	fmt.Println(countChineseCharacters(model.Cube333.Cn()))
	fmt.Println(countChineseCharacters(model.CubeSq1.Cn()))
}
