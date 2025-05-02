package words

import (
	"strings"
	"unicode"
)

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
