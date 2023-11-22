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

// CutMsgWithFields
// 切割格式： ${header}[${title}] ${values...}
func CutMsgWithFields(input string, cut string) (header string, title string, values []string) {

	// 这里是为了强制让]后面有空格
	input = strings.ReplaceAll(input, "]", "] ")

	// ${header}[${title}]
	// ${values...}
	parts := strings.Fields(input)

	if len(parts) >= 1 {
		headers := strings.SplitN(parts[0], "[", -1)
		if len(headers) >= 1 {
			header = headers[0]
		}
		if len(headers) >= 2 {
			title = headers[1][:len(headers[1])-1]
		}
	}

	if len(parts) >= 2 {
		otherValue := strings.Join(parts[1:], " ")
		values = strings.Split(otherValue, cut)
	}

	return header, title, values
}
