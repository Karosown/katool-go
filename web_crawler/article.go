package web_crawler

import (
	"bytes"
	"codeberg.org/readeck/go-readability"
	nurl "net/url"
	"time"

	"github.com/go-rod/rod"
)

// GetArticleWithURL 通过URL获取文章内容
// GetArticleWithURL gets article content by URL
func GetArticleWithURL(url string, requestModifiers ...RequestWith) Article {
	return DefaultClient.GetArticleWithURL(url, requestModifiers...)
}

// GetArticleWithURL 通过URL获取文章内容
// GetArticleWithURL gets article content by URL
func (c *Client) GetArticleWithURL(url string, requestModifiers ...RequestWith) Article {
	article, err := c.fromURL(url, 30*time.Second, requestModifiers...)
	return Article{article, WebReaderError{err}}
}

// GetArticleWithChrome 使用Chrome浏览器获取文章内容（支持JavaScript渲染）
// GetArticleWithChrome gets article content using Chrome browser (supports JavaScript rendering)
func GetArticleWithChrome(url string, rendorFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) Article {
	return DefaultClient.GetArticleWithChrome(url, rendorFunc, restartCondition, i...)
}

// GetArticleWithChrome 使用Chrome浏览器获取文章内容（支持JavaScript渲染）
// GetArticleWithChrome gets article content using Chrome browser (supports JavaScript rendering)
func (c *Client) GetArticleWithChrome(url string, rendorFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) Article {
	code := GetArticleWithSourceCode(c.ReadSourceCode(url, "", rendorFunc), url)
	if i != nil && i[0] <= 5 && restartCondition != nil && restartCondition(code) {
		if c.getChrome() != nil {
			c.getChrome().ReStart()
		}
		return c.GetArticleWithChrome(url, rendorFunc, restartCondition, i[0]+1)
	}
	return code
}

// GetArticleWithSourceCode 从源代码获取文章内容
// GetArticleWithSourceCode gets article content from source code
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
