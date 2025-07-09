package render

import (
	"github.com/go-rod/rod"
	"github.com/karosown/katool-go/web_crawler"
)

// 渲染策略和函数定义
// Rendering strategies and function definitions
var (
	// DefaultRestartPolicy 默认重启策略，当文章标题为空时重启
	// DefaultRestartPolicy is the default restart policy that restarts when article title is empty
	DefaultRestartPolicy = func(article web_crawler.Article) bool {
		return article.Title == ""
	}

	// NoRender 不进行渲染，立即停止页面加载
	// NoRender performs no rendering and immediately stops page loading
	NoRender = func(page *rod.Page) {
		page.MustStopLoading()
	}

	// Render 进行完整渲染，等待页面加载完成
	// Render performs full rendering and waits for page loading to complete
	Render = func(page *rod.Page) {
		page.MustWaitLoad()
	}
)
