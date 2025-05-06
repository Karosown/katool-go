package test

import (
	"testing"
	"unicode"

	"github.com/karosown/katool-go/words"
)

func TestLanguate(t *testing.T) {
	println(words.ContainsLanguage("卧槽", unicode.Han))
	println(words.ContainsLanguage("卧槽123", unicode.Han))
	println(words.ContainsLanguage("123qweqwe", unicode.Han))
	println(words.ContainsLanguage("123qweqwe", unicode.Number, unicode.Letter, unicode.Han))
	println(words.ContainsLanguage("qweqwe", unicode.Letter))
	println(words.ContainsLanguage("123", unicode.Number))
	println("------")
	println(words.OnlyLanguage("卧槽", unicode.Han))
	println(words.OnlyLanguage("卧槽12313", unicode.Han))
	println(words.OnlyLanguage("卧槽12313", unicode.Han, unicode.Number))
	println("------")
	println(words.SubString("```json\n"+
		"{测试}\n"+
		"```", "```json", "```"))
}
