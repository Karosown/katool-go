package web_crawler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/karosown/katool-go/web_crawler/core"
)

// ChallengePolicy 挑战策略接口
// ChallengePolicy challenge policy interface
type ChallengePolicy interface {
	// ShouldHandle 检查是否应该处理该响应
	// ShouldHandle checks if the response should be handled
	ShouldHandle(resp *http.Response) bool

	// Solve 解决挑战，返回Cookies, UserAgent和错误
	// Solve solves the challenge, returning Cookies, UserAgent, and error
	Solve(chrome *core.Contain, url string) ([]*http.Cookie, string, error)
}

// CFSession Cloudflare会话缓存
// CFSession Cloudflare session cache
type CFSession struct {
	Cookies   []*http.Cookie
	UserAgent string
	Expires   time.Time
}

var (
	// cfCache 缓存 Cloudflare 的验证信息 (host -> *CFSession)
	// cfCache caches Cloudflare verification info (host -> *CFSession)
	cfCache     = make(map[string]*CFSession)
	cfCacheLock sync.RWMutex
	// GlobalPolicies 全局策略列表
	// GlobalPolicies global policy list
	GlobalPolicies []ChallengePolicy
)

func init() {
	// 注册默认的Cloudflare策略
	GlobalPolicies = append(GlobalPolicies, &CloudflarePolicy{})
}

// AntiBotConfig 反爬虫配置
// AntiBotConfig anti-bot configuration
type AntiBotConfig struct {
	// EnableCloudflare 是否启用 Cloudflare 绕过
	// EnableCloudflare whether to enable Cloudflare bypass
	EnableCloudflare bool
}

// DefaultAntiBotConfig 返回默认配置
// DefaultAntiBotConfig returns default configuration
func DefaultAntiBotConfig() *AntiBotConfig {
	return &AntiBotConfig{
		EnableCloudflare: false, // 默认不启用启用
	}
}

// AntiBotTransport 自定义Transport，处理各种反爬虫验证
// AntiBotTransport custom Transport to handle various anti-bot verifications
type AntiBotTransport struct {
	Transport http.RoundTripper
	Chrome    *core.Contain
	Policies  []ChallengePolicy
	Config    *AntiBotConfig
}

// RoundTrip 实现RoundTripper接口
// RoundTrip implements RoundTripper interface
func (t *AntiBotTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 应用配置过滤策略
	// Apply config to filter policies
	config := t.Config
	if config == nil {
		config = DefaultAntiBotConfig()
	}

	// 获取主机名，用于缓存操作
	// Get hostname for cache operations
	host := req.URL.Host

	// 尝试注入缓存的Cookie和UA (目前仅针对Cloudflare)
	// Attempt to inject cached Cookies and UA (currently only for Cloudflare)
	// 只有在启用 Cloudflare 时才注入缓存
	// Only inject cache when Cloudflare is enabled
	if config.EnableCloudflare {
		cfCacheLock.RLock()
		if session, ok := cfCache[host]; ok {
			for _, cookie := range session.Cookies {
				req.AddCookie(cookie)
			}
			if session.UserAgent != "" {
				req.Header.Set("User-Agent", session.UserAgent)
			}
		}
		cfCacheLock.RUnlock()
	}

	// 执行原始请求
	// Execute original request
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 遍历所有策略，检查是否需要处理
	// Iterate through all policies, check if handling is needed
	policies := t.Policies
	if len(policies) == 0 {
		policies = GlobalPolicies
	}

	for _, policy := range policies {
		// 如果策略是 CloudflarePolicy 且配置中禁用了 Cloudflare，则跳过
		// If policy is CloudflarePolicy and Cloudflare is disabled in config, skip it
		if _, ok := policy.(*CloudflarePolicy); ok && !config.EnableCloudflare {
			continue
		}

		if policy.ShouldHandle(resp) {
			// 关闭旧的响应体
			resp.Body.Close()

			// 尝试解决验证
			cookies, userAgent, err := policy.Solve(t.Chrome, req.URL.String())
			if err != nil {
				return nil, fmt.Errorf("failed to solve challenge: %v", err)
			}

			// 更新缓存 (如果是Cloudflare策略)
			// Update cache (if it is Cloudflare policy)
			if _, ok := policy.(*CloudflarePolicy); ok {
				cfCacheLock.Lock()
				cfCache[host] = &CFSession{
					Cookies:   cookies,
					UserAgent: userAgent,
					Expires:   time.Now().Add(1 * time.Hour),
				}
				cfCacheLock.Unlock()
			}

			// 更新请求头
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			if userAgent != "" {
				req.Header.Set("User-Agent", userAgent)
			}

			// 重试请求
			return transport.RoundTrip(req)
		}
	}

	return resp, nil
}

// CloudflarePolicy Cloudflare策略实现
type CloudflarePolicy struct{}

func (p *CloudflarePolicy) ShouldHandle(resp *http.Response) bool {
	// 检查状态码
	if resp.StatusCode != 503 && resp.StatusCode != 403 {
		return false
	}
	// 检查Server头
	server := resp.Header.Get("Server")
	return strings.EqualFold(server, "cloudflare") || strings.EqualFold(server, "cloudflare-nginx")
}

func (p *CloudflarePolicy) Solve(chrome *core.Contain, url string) ([]*http.Cookie, string, error) {
	return SolveChallenge(chrome, url)
}

// SolveChallenge 使用WebChrome解决验证（通用逻辑，保留给CF使用）
func SolveChallenge(chrome *core.Contain, url string) ([]*http.Cookie, string, error) {
	if chrome == nil {
		return nil, "", fmt.Errorf("WebChrome not initialized")
	}

	page := chrome.PageWithStealth(url)
	defer page.Close()

	// 并发尝试点击验证框
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				attemptClickChallenge(page)
			}
		}
	}()

	// 每100ms检查一次Cookie
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, "", fmt.Errorf("timeout waiting for challenge")
		case <-ticker.C:
			cookies, err := page.Cookies([]string{url})
			if err != nil {
				continue
			}

			found := false
			for _, c := range cookies {
				if c.Name == "cf_clearance" {
					found = true
					break
				}
			}

			if found {
				netCookies := cookies
				httpCookies := make([]*http.Cookie, len(netCookies))
				for i, c := range netCookies {
					httpCookies[i] = &http.Cookie{
						Name:    c.Name,
						Value:   c.Value,
						Domain:  c.Domain,
						Path:    c.Path,
						Expires: c.Expires.Time(),
					}
				}

				ua, err := page.Eval("() => navigator.userAgent")
				if err != nil {
					return nil, "", err
				}

				return httpCookies, ua.Value.String(), nil
			}
		}
	}
}

// attemptClickChallenge 尝试寻找并点击Cloudflare验证框
func attemptClickChallenge(page *rod.Page) {
	iframes, err := page.Elements("iframe")
	if err == nil {
		for _, frame := range iframes {
			framePage, err := frame.Frame()
			if err != nil {
				continue
			}
			if el, err := framePage.Element(`input[type="checkbox"]`); err == nil {
				box, err := el.Shape()
				if err == nil {
					center := box.OnePointInside()
					page.Mouse.MoveTo(*center)
					time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
					page.Mouse.Click(proto.InputMouseButtonLeft, 1)
				}
				return
			}
		}
	}
	if el, err := page.Element("#challenge-stage input[type='checkbox']"); err == nil {
		el.Click(proto.InputMouseButtonLeft, 1)
	}
}

// NewCloudscraperClient 创建一个支持反爬虫Bypass的HTTP Client
// NewCloudscraperClient creates an HTTP Client supporting anti-bot bypass
// Deprecated: use NewAntiBotClient instead
func NewCloudscraperClient(chrome *core.Contain, timeout time.Duration) *http.Client {
	return NewAntiBotClient(chrome, timeout)
}

// NewAntiBotClient 创建一个支持反爬虫Bypass的HTTP Client
// NewAntiBotClient creates an HTTP Client supporting anti-bot bypass
func NewAntiBotClient(chrome *core.Contain, timeout time.Duration, policies ...ChallengePolicy) *http.Client {
	return NewAntiBotClientWithConfig(chrome, timeout, DefaultAntiBotConfig(), policies...)
}

// NewAntiBotClientWithConfig 创建一个支持反爬虫Bypass的HTTP Client（带配置）
// NewAntiBotClientWithConfig creates an HTTP Client supporting anti-bot bypass (with config)
func NewAntiBotClientWithConfig(chrome *core.Contain, timeout time.Duration, config *AntiBotConfig, policies ...ChallengePolicy) *http.Client {
	if config == nil {
		config = DefaultAntiBotConfig()
	}
	return &http.Client{
		Timeout: timeout,
		Transport: &AntiBotTransport{
			Transport: http.DefaultTransport,
			Chrome:    chrome,
			Policies:  policies,
			Config:    config,
		},
	}
}
