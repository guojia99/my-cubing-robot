package utils

import (
	"fmt"
	"strings"
)

func FormatData(data string) {
	lines := strings.Split(data, "\n")

	maxNameLength := 0
	maxScoreLength := 0

	for _, line := range lines {
		parts := strings.Split(line, " | ")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		score := parts[1]

		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
		if len(score) > maxScoreLength {
			maxScoreLength = len(score)
		}
	}

	formatString := fmt.Sprintf("%%-%ds | %%%ds\n", maxNameLength, maxScoreLength)

	for _, line := range lines {
		parts := strings.Split(line, " | ")
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		score := parts[1]

		formattedLine := fmt.Sprintf(formatString, name, score)
		fmt.Println(formattedLine)
	}
}
