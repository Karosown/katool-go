package words

import (
	"unicode"
)

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
