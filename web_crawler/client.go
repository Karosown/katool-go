package web_crawler

import (
	"github.com/karosown/katool-go/web_crawler/core"
)

// Client 网络爬虫客户端
// Client Web crawler client
type Client struct {
	Chrome   *core.Contain
	Policies []ChallengePolicy
	Config   *AntiBotConfig
}

// DefaultClient 默认客户端实例
// DefaultClient default client instance
var DefaultClient = &Client{
	Policies: GlobalPolicies,
	Config:   DefaultAntiBotConfig(),
}

// NewClient 创建新的客户端
// NewClient creates a new client
func NewClient(chrome *core.Contain, policies ...ChallengePolicy) *Client {
	// 如果没有提供策略，则使用默认的全局策略
	// If no policies provided, use default global policies
	p := policies
	if len(p) == 0 {
		p = GlobalPolicies
	}
	return &Client{
		Chrome:   chrome,
		Policies: p,
		Config:   DefaultAntiBotConfig(),
	}
}

// NewClientWithConfig 创建新的客户端（带配置）
// NewClientWithConfig creates a new client (with config)
func NewClientWithConfig(chrome *core.Contain, config *AntiBotConfig, policies ...ChallengePolicy) *Client {
	// 如果没有提供策略，则使用默认的全局策略
	// If no policies provided, use default global policies
	p := policies
	if len(p) == 0 {
		p = GlobalPolicies
	}
	if config == nil {
		config = DefaultAntiBotConfig()
	}
	return &Client{
		Chrome:   chrome,
		Policies: p,
		Config:   config,
	}
}

// AddPolicy 添加验证策略
// AddPolicy adds challenge policy
func (c *Client) AddPolicy(policy ChallengePolicy) {
	c.Policies = append(c.Policies, policy)
}

// getChrome 获取Chrome实例，优先使用Client内部实例，否则使用全局WebChrome
// getChrome gets the Chrome instance, prioritizing the internal instance, otherwise using global WebChrome
func (c *Client) getChrome() *core.Contain {
	if c != nil && c.Chrome != nil {
		return c.Chrome
	}
	return WebChrome
}

// GetChrome 获取Chrome实例（公开方法）
// GetChrome gets the Chrome instance (public method)
func (c *Client) GetChrome() *core.Contain {
	return c.getChrome()
}

// GetConfig 获取配置（公开方法）
// GetConfig gets the config (public method)
func (c *Client) GetConfig() *AntiBotConfig {
	if c != nil && c.Config != nil {
		return c.Config
	}
	return DefaultAntiBotConfig()
}

// GetPolicies 获取策略列表（公开方法）
// GetPolicies gets the policies list (public method)
func (c *Client) GetPolicies() []ChallengePolicy {
	if c != nil && len(c.Policies) > 0 {
		return c.Policies
	}
	return GlobalPolicies
}
