package process

import (
	"errors"
	"strings"

	"k8s.io/klog"
)

func PrefixMap(process ...Process) map[string]Process {
	var out = make(map[string]Process)
	for _, p := range process {
		for _, pf := range p.Prefix() {
			out[pf] = p
			klog.Infof("key %s", pf)
		}
	}

	return out
}

func findSubSeq(input string) []string {
	var cache [][]rune
	msg := []rune(input)

	for i := 0; i < len(msg)+1; i++ {
		cache = append(cache, msg[:i])
		if i > MaxKeyLength {
			break
		}
	}

	var out []string
	for _, val := range cache {
		out = append(out, string(val))
	}
	return out
}

func CheckPrefix(msg string, processMap map[string]Process) (Process, error) {
	for _, key := range findSubSeq(msg) {
		if p, ok := processMap[key]; ok {
			return p, nil
		}
	}
	return nil, errors.New("命令未找到")
}

func ReplaceAll(s, new string, old ...string) string {
	for _, o := range old {
		s = strings.ReplaceAll(s, o, new)
	}
	return s
}
