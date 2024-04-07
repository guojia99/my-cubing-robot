package process

import (
	"fmt"
	"testing"
)

func TestPreEnter_CheckPrefix(t *testing.T) {
	in := "三阶 D D D D D"

	c := &PreEnter{}

	fmt.Println(c.CheckPrefix(in))
}
