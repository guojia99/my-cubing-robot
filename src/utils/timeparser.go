package utils

import (
	"fmt"
	"math"

	"github.com/guojia99/my-cubing/src/core/model"
)

func TimeParser(score model.Score, isAvg bool) string {

	if score.Project == model.Cube333FM {
		if isAvg {
			return fmt.Sprintf("%2.2f", score.Avg)
		}
		return fmt.Sprintf("%d", int(score.Best))
	}

	in := score.Best
	if isAvg {
		in = score.Avg
	}
	if score.Project.RouteType() == model.RouteTypeRepeatedly {
		in = score.Result3
	}

	if in < 60 {
		return fmt.Sprintf("%0.2f", in)
	}
	m := int(math.Floor(in) / 60)
	s := in - float64(m*60)

	ss := fmt.Sprintf("%0.2f", s)
	if s < 10 {
		ss = fmt.Sprintf("0%0.2f", s)
	}

	return fmt.Sprintf("%d:%s", m, ss)
}
