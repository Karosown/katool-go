package web_crawler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/karosown/katool-go/web_crawler/core"
)

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
)

// CloudscraperTransport 自定义Transport，处理Cloudflare验证
// CloudscraperTransport custom Transport to handle Cloudflare verification
type CloudscraperTransport struct {
	Transport http.RoundTripper
	Chrome    *core.Contain
}

// RoundTrip 实现RoundTripper接口
// RoundTrip implements RoundTripper interface
func (t *CloudscraperTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 尝试注入缓存的Cookie和UA
	// Attempt to inject cached Cookies and UA
	host := req.URL.Host
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

	// 检查是否遇到Cloudflare验证
	// Check if Cloudflare verification is encountered
	if IsCloudflareChallenge(resp) {
		// 关闭旧的响应体
		resp.Body.Close()

		// 尝试解决验证
		// Attempt to solve verification
		cookies, userAgent, err := SolveChallenge(t.Chrome, req.URL.String())
		if err != nil {
			return nil, fmt.Errorf("failed to solve Cloudflare challenge: %v", err)
		}

		// 更新缓存
		// Update cache
		cfCacheLock.Lock()
		cfCache[host] = &CFSession{
			Cookies:   cookies,
			UserAgent: userAgent,
			// 假设有效期1小时，或者从cookie中获取最小过期时间
			// Assume 1 hour validity, or get min expiration from cookies
			Expires: time.Now().Add(1 * time.Hour),
		}
		cfCacheLock.Unlock()

		// 更新请求头
		// Update request headers
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		if userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}

		// 重试请求
		// Retry request
		return transport.RoundTrip(req)
	}

	return resp, nil
}

// IsCloudflareChallenge 检查响应是否为Cloudflare验证页面
// IsCloudflareChallenge checks if the response is a Cloudflare verification page
func IsCloudflareChallenge(resp *http.Response) bool {
	// 检查状态码
	// Check status code
	if resp.StatusCode != 503 && resp.StatusCode != 403 {
		return false
	}

	// 检查Server头
	// Check Server header
	server := resp.Header.Get("Server")
	if !strings.EqualFold(server, "cloudflare") && !strings.EqualFold(server, "cloudflare-nginx") {
		return false
	}

	// 进一步检查（可选，读取body可能会消耗流，这里主要依赖头和状态码）
	// Further checks (optional, reading body might consume stream, mainly relying on headers and status code here)
	// 如果需要更精确，可以Peek body，但net/http的Response body是ReadCloser

	return true
}

// SolveChallenge 使用WebChrome解决Cloudflare验证
// SolveChallenge solves Cloudflare verification using WebChrome
func SolveChallenge(chrome *core.Contain, url string) ([]*http.Cookie, string, error) {
	if chrome == nil {
		return nil, "", fmt.Errorf("WebChrome not initialized")
	}

	// 使用Stealth模式创建页面
	// Create page using Stealth mode
	page := chrome.PageWithStealth(url)
	defer page.Close()

	// 移除 WaitLoad，直接轮询 Cookie，一旦 Cloudflare 写入 Cookie 立即返回，无需等待重定向后的页面加载
	// Remove WaitLoad, poll Cookie directly. Return immediately once Cloudflare writes Cookie, no need to wait for redirected page load.
	// page.MustWaitLoad()

	// 使用Context控制超时，缩短轮询间隔以提高效率
	// Use Context to control timeout, shorten polling interval to improve efficiency
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 每100ms检查一次，追求极致响应速度
	// Check every 100ms, pursuing extreme response speed
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, "", fmt.Errorf("timeout waiting for Cloudflare challenge")
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
				// 获取最终的Cookies和User-Agent
				// Get final Cookies and User-Agent
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

// NewCloudscraperClient 创建一个支持Cloudflare bypass的HTTP Client
// NewCloudscraperClient creates an HTTP Client supporting Cloudflare bypass
func NewCloudscraperClient(chrome *core.Contain, timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &CloudscraperTransport{
			Transport: http.DefaultTransport,
			Chrome:    chrome,
		},
	}
}
