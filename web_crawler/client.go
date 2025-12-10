package web_crawler

import (
	"github.com/karosown/katool-go/web_crawler/core"
)

// Client 网络爬虫客户端
// Client Web crawler client
type Client struct {
	Chrome *core.Contain
}

// DefaultClient 默认客户端实例
// DefaultClient default client instance
var DefaultClient = &Client{}

// NewClient 创建新的客户端
// NewClient creates a new client
func NewClient(chrome *core.Contain) *Client {
	return &Client{
		Chrome: chrome,
	}
}

// getChrome 获取Chrome实例，优先使用Client内部实例，否则使用全局WebChrome
// getChrome gets the Chrome instance, prioritizing the internal instance, otherwise using global WebChrome
func (c *Client) getChrome() *core.Contain {
	if c != nil && c.Chrome != nil {
		return c.Chrome
	}
	return WebChrome
}
