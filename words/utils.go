package words

import (
	"strings"
	"unicode"
)

// SubString 从字符串中提取子串（从start到end之间的内容）
// SubString extracts substring from string (content between start and end)
func SubString(s, start, end string) string {
	fi := strings.Index(s, start)
	if fi == -1 {
		fi = 0
	}
	li := strings.LastIndex(s, end)
	if li == -1 {
		li = len(s)
	}
	return s[fi+len(start) : li]
}

// ContainsLanguage 检查字符串是否包含指定语言字符
// ContainsLanguage checks if string contains characters of specified languages
func ContainsLanguage(s string, languages ...*unicode.RangeTable) bool {
	for _, v := range s {
		for _, language := range languages {
			is := unicode.Is(language, v)
			if is {
				return true
			}
		}
	}
	return false
}

// OnlyLanguage 检查字符串是否只包含指定语言字符
// OnlyLanguage checks if string contains only characters of specified languages
func OnlyLanguage(s string, languages ...*unicode.RangeTable) bool {
	for _, v := range s {
		is := false
		for _, language := range languages {
			is = is || unicode.Is(language, v)
		}
		if !is {
			return false
		}
	}
	return true
}

// CaseShift 对已有字母进行大小写互转，大写转换为小写，小写转换为大写
// CaseShift toggles case for existing letters, uppercase to lowercase and vice versa
func CaseShift(str string) string {
	res := ""
	for _, item := range str {
		res += string(item ^ 32)
	}
	return res
}
