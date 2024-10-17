package bf

import "unicode/utf8"

func splitByFixedLength(s string, length int) []string {
	var result []string
	for len(s) > 0 {
		// 如果剩余的字符串长度小于或等于切割长度，直接添加到结果中并退出循环
		if utf8.RuneCountInString(s) <= length {
			result = append(result, s)
			break
		}
		// 否则，添加指定长度的子串到结果中
		substr := string([]rune(s)[:length])
		result = append(result, substr)
		// 切去已经处理的部分
		s = string([]rune(s)[length:])
	}
	return result
}
