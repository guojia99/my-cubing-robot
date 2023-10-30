package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/guojia99/my-cubing-core/model"
)

func TimeParser(score model.Score, isAvg bool) string {

	if score.Project.RouteType() == model.RouteTypeRepeatedly {
		return fmt.Sprintf("%2.0f / %2.0f %s", score.Result1, score.Result2, TimeParser(
			model.Score{Best: score.Result3, Project: model.Cube333, Avg: score.Result3}, isAvg,
		))
	}

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
	if in <= model.DNF {
		return "DNF"
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

func ParserTimeToSeconds(t string) float64 {
	if t == "DNF" || strings.ContainsAny(t, "dD") {
		return model.DNF
	}
	if t == "DNS" || strings.Contains(t, "s") {
		return model.DNS
	}
	// 解析纯秒数格式
	if regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(t) {
		seconds, _ := strconv.ParseFloat(t, 64)
		return seconds
	}

	// 解析分+秒格式
	if regexp.MustCompile(`^\d{1,3}:\d{1,3}(\.\d+)?$`).MatchString(t) {
		parts := strings.Split(t, ":")
		minutes, _ := strconv.ParseFloat(parts[0], 64)
		seconds, _ := strconv.ParseFloat(parts[1], 64)
		return minutes*60 + seconds
	}

	// 解析时+分+秒格式
	if regexp.MustCompile(`^\d{1,3}:\d{1,3}:\d{1,3}(\.\d+)?$`).MatchString(t) {
		parts := strings.Split(t, ":")
		hours, _ := strconv.ParseFloat(parts[0], 64)
		minutes, _ := strconv.ParseFloat(parts[1], 64)
		seconds, _ := strconv.ParseFloat(parts[2], 64)
		return hours*3600 + minutes*60 + seconds
	}

	return model.DNF
}
