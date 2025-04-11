package render

import (
	"github.com/go-rod/rod"
	"github.com/karosown/katool-go/web_crawler"
)

var (
	DefaultRestartPolicy = func(article web_crawler.Article) bool {
		return article.Title == ""
	}
	NoRender = func(page *rod.Page) {
		page.MustStopLoading()
	}
	Render = func(page *rod.Page) {
		page.MustWaitLoad()
	}
)
