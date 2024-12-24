package web_crawler

import (
	"bytes"
	nurl "net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-shiori/go-readability"
)

func GetArticleWithURL(url string, requestModifiers ...RequestWith) Article {
	article, err := fromURL(url, 30*time.Second, requestModifiers...)
	return Article{article, WebReaderError{err}}
}
func GetArticleWithChrome(url string, rendorFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) Article {
	code := GetArticleWithSourceCode(ReadSourceCode(url, "", rendorFunc), url)
	if i != nil && i[0] <= 5 && restartCondition != nil && restartCondition(code) {
		WebChrome.ReStart()
		return GetArticleWithChrome(url, rendorFunc, restartCondition, i[0]+1)
	}
	return code
}
func GetArticleWithSourceCode(code SourceCode, url string) Article {
	if code.IsErr() {
		return Article{
			Article:        readability.Article{},
			WebReaderError: code.WebReaderError,
		}
	}
	parse, err := nurl.Parse(url)
	article, err := readability.FromReader(bytes.NewBufferString(code.WebReaderString.String()), parse)
	return Article{article, WebReaderError{err}}
}
