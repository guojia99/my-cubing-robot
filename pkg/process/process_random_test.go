package process

import (
	"fmt"
	"testing"
)

func Test_randomWithMsg(t *testing.T) {

	var msgs = []string{
		"",
		"3bf edge",
		"3bf edge*3",
		"3bf corner*3",
		"BCEFHIJKLMNPQRSTWXYZ ",
		"3bf xedge*3",
		"*3",
	}
	for _, msg := range msgs {
		t.Run(
			msg, func(t *testing.T) {
				fmt.Println(randomWithMsg(msg))
			},
		)
	}

}
