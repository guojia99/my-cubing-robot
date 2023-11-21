package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func CountChineseCharacters(text string) int {
	count := 0
	for _, char := range text {
		if utf8.RuneLen(char) >= 1 { // 判断字符是否占用多个字节（多个码点）
			count++
		}
	}
	return count
}
func CountChineseCharacters2(text string) int {
	count := 0
	for _, char := range text {
		if utf8.RuneLen(char) >= 2 { // 判断字符是否占用多个字节（多个码点）
			count += 2
		} else {
			count += 1
		}
	}
	return count
}

func TB(in any, n int) string {
	n *= 2
	msg := fmt.Sprintf("%+v", in)

	kn := CountChineseCharacters2(msg)
	if kn >= n {
		return msg
	}

	msg += strings.Repeat(" ", n-kn)
	return msg
}
