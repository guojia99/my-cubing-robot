package utils

import "strings"

func HasPrefixList(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
