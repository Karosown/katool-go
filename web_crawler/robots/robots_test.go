package robots

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/karosown/katool-go/web_crawler"
)

// 示例 robots.txt 内容
// Example robots.txt content
const exampleRobotsTxt = `
# robots.txt for example.com
User-agent: *
Disallow: /admin/
Disallow: /private/
Allow: /public/
Crawl-delay: 10

User-agent: Googlebot
Allow: /
Crawl-delay: 1

Sitemap: https://example.com/sitemap.xml
Sitemap: https://example.com/sitemap-news.xml
`

func TestParseRobotsTxt(t *testing.T) {
	robots := NewRobotsTxt()
	reader := strings.NewReader(exampleRobotsTxt)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// 检查 User-Agent 组
	// Check User-Agent groups
	if len(robots.UserAgents) != 2 {
		t.Errorf("Expected 2 user agent groups, got %d", len(robots.UserAgents))
	}

	// 检查 "*" 用户代理组
	// Check "*" user agent group
	starGroup, ok := robots.UserAgents["*"]
	if !ok {
		t.Fatal("Expected '*' user agent group not found")
	}

	if len(starGroup.Disallow) != 2 {
		t.Errorf("Expected 2 disallowed paths, got %d", len(starGroup.Disallow))
	}

	if starGroup.CrawlDelay != 10 {
		t.Errorf("Expected crawl delay 10, got %d", starGroup.CrawlDelay)
	}

	// 检查 Googlebot 用户代理组
	// Check Googlebot user agent group
	googleGroup, ok := robots.UserAgents["Googlebot"]
	if !ok {
		t.Fatal("Expected 'Googlebot' user agent group not found")
	}

	if googleGroup.CrawlDelay != 1 {
		t.Errorf("Expected crawl delay 1, got %d", googleGroup.CrawlDelay)
	}

	// 检查 Sitemap
	// Check Sitemap
	if len(robots.Sitemaps) != 2 {
		t.Errorf("Expected 2 sitemaps, got %d", len(robots.Sitemaps))
	}
}

func TestIsAllowed(t *testing.T) {
	robots := NewRobotsTxt()
	reader := strings.NewReader(exampleRobotsTxt)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// 测试禁止的路径
	// Test disallowed paths
	testCases := []struct {
		url       string
		userAgent string
		allowed   bool
		desc      string
	}{
		{"https://example.com/admin/", "*", false, "admin path should be disallowed"},
		{"https://example.com/private/", "*", false, "private path should be disallowed"},
		{"https://example.com/public/", "*", true, "public path should be allowed"},
		{"https://example.com/", "*", true, "root path should be allowed"},
		{"https://example.com/admin/", "Googlebot", true, "Googlebot should be allowed"},
		{"https://example.com/private/", "Googlebot", true, "Googlebot should be allowed"},
	}

	for _, tc := range testCases {
		result := robots.IsAllowed(tc.url, tc.userAgent)
		if result != tc.allowed {
			t.Errorf("%s: Expected %v, got %v for URL %s with User-Agent %s",
				tc.desc, tc.allowed, result, tc.url, tc.userAgent)
		}
	}
}

func TestGetCrawlDelay(t *testing.T) {
	robots := NewRobotsTxt()
	reader := strings.NewReader(exampleRobotsTxt)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// 测试获取爬取延迟
	// Test getting crawl delay
	delay := robots.GetCrawlDelay("*")
	if delay != 10 {
		t.Errorf("Expected crawl delay 10 for '*', got %d", delay)
	}

	delay = robots.GetCrawlDelay("Googlebot")
	if delay != 1 {
		t.Errorf("Expected crawl delay 1 for 'Googlebot', got %d", delay)
	}

	delay = robots.GetCrawlDelay("UnknownBot")
	if delay != 10 {
		t.Errorf("Expected crawl delay 10 for 'UnknownBot' (should match '*'), got %d", delay)
	}
}

func TestGetSitemaps(t *testing.T) {
	robots := NewRobotsTxt()
	reader := strings.NewReader(exampleRobotsTxt)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	sitemaps := robots.GetSitemaps()
	if len(sitemaps) != 2 {
		t.Errorf("Expected 2 sitemaps, got %d", len(sitemaps))
	}

	expectedSitemaps := []string{
		"https://example.com/sitemap.xml",
		"https://example.com/sitemap-news.xml",
	}

	for i, expected := range expectedSitemaps {
		if i >= len(sitemaps) || sitemaps[i] != expected {
			t.Errorf("Expected sitemap %s, got %s", expected, sitemaps[i])
		}
	}
}

func TestMatchesPath(t *testing.T) {
	robots := NewRobotsTxt()

	testCases := []struct {
		path    string
		pattern string
		matches bool
		desc    string
	}{
		{"/admin/", "/admin/", true, "exact match"},
		{"/admin/test", "/admin/", true, "path starts with pattern"},
		{"/public/", "/admin/", false, "path doesn't match"},
		{"/", "/", true, "root path matches root pattern"},
		{"/any/path", "/", true, "any path matches root pattern"},
		{"/test/*/path", "/test/*/path", true, "wildcard pattern"},
	}

	for _, tc := range testCases {
		result := robots.matchesPathPattern(tc.path, tc.pattern)
		if result != tc.matches {
			t.Errorf("%s: Expected %v, got %v for path '%s' and pattern '%s'",
				tc.desc, tc.matches, result, tc.path, tc.pattern)
		}
	}
}

// ExampleFetchFromURL 演示如何从 URL 获取 robots.txt
// ExampleFetchFromURL demonstrates how to fetch robots.txt from URL
func ExampleFetchFromURL() {
	// 注意：这是一个示例，实际使用时需要有效的 URL
	// Note: This is an example, actual usage requires a valid URL
	robots, err := FetchFromURL("https://example.com", 10*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 检查 URL 是否被允许
	// Check if URL is allowed
	allowed := robots.IsAllowed("https://example.com/admin/", "MyBot/1.0")
	fmt.Printf("Allowed: %v\n", allowed)

	// 获取 sitemap
	// Get sitemaps
	sitemaps := robots.GetSitemaps()
	fmt.Printf("Sitemaps: %v\n", sitemaps)
}

// ExampleFetchFromClient 演示如何使用 Client 获取 robots.txt
// ExampleFetchFromClient demonstrates how to fetch robots.txt using Client
func ExampleFetchFromClient() {
	// 创建客户端
	// Create client
	client := web_crawler.NewClient(nil)

	// 获取 robots.txt
	// Fetch robots.txt
	robots, err := FetchFromClient(client, "https://example.com", 10*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 检查 URL 是否被允许
	// Check if URL is allowed
	allowed := robots.IsAllowed("https://example.com/test/", "MyBot/1.0")
	fmt.Printf("Allowed: %v\n", allowed)

	// 获取爬取延迟
	// Get crawl delay
	delay := robots.GetCrawlDelay("MyBot/1.0")
	fmt.Printf("Crawl delay: %d seconds\n", delay)
}

// TestCustomPath 测试自定义路径功能
// TestCustomPath tests custom path functionality
func TestCustomPath(t *testing.T) {
	// 这个测试主要验证函数签名是否正确，实际测试需要真实的服务器
	// This test mainly verifies that function signatures are correct, actual testing requires a real server
	_ = FetchFromURL
	_ = FetchFromClient
	_ = FetchFromURLWithAutoDetect
	_ = FetchFromClientWithAutoDetect
	_ = FetchFromDirectURL
	_ = FetchFromClientDirectURL
}

// ExampleFetchFromURL_customPath 演示如何使用自定义路径获取 robots.txt
// ExampleFetchFromURL_customPath demonstrates how to fetch robots.txt with custom path
func ExampleFetchFromURL_customPath() {
	// 使用自定义路径
	// Use custom path
	robots, err := FetchFromURL("https://example.com", 10*time.Second, "/custom/robots.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	allowed := robots.IsAllowed("https://example.com/test/", "MyBot/1.0")
	fmt.Printf("Allowed: %v\n", allowed)
}

// ExampleFetchFromURLWithAutoDetect 演示自动检测多个路径
// ExampleFetchFromURLWithAutoDetect demonstrates auto-detecting multiple paths
func ExampleFetchFromURLWithAutoDetect() {
	// 自动尝试多个常见路径
	// Automatically try multiple common paths
	robots, err := FetchFromURLWithAutoDetect("https://example.com", 10*time.Second, "/custom/robots")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	allowed := robots.IsAllowed("https://example.com/test/", "MyBot/1.0")
	fmt.Printf("Allowed: %v\n", allowed)
}

// ExampleFetchFromDirectURL 演示直接使用完整 URL
// ExampleFetchFromDirectURL demonstrates using full URL directly
func ExampleFetchFromDirectURL() {
	// 直接使用完整 URL
	// Use full URL directly
	robots, err := FetchFromDirectURL("https://example.com/custom/path/robots.txt", 10*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	allowed := robots.IsAllowed("https://example.com/test/", "MyBot/1.0")
	fmt.Printf("Allowed: %v\n", allowed)
}

func TestRobots(t *testing.T) {
	url, err := FetchFromDirectURL("https://www.landiannews.com/robots.txt", 10*time.Second)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	allowed := url.IsAllowed("https://www.landiannews.com/123123213213?post_type=topic", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	t.Logf("Allowed: %v\n", allowed)
}

// TestQueryStringMatching 测试查询字符串匹配
// TestQueryStringMatching tests query string matching
func TestQueryStringMatching(t *testing.T) {
	robotsTxtContent := `
User-agent: *
Disallow: /my?project_id=*
Disallow: /login?ok_url=*
Disallow: /loginWithPwd?ok_url=*
Disallow: /loginWithPw?ok_url=*
Disallow: /users
Disallow: /search/
`

	robots := NewRobotsTxt()
	reader := strings.NewReader(robotsTxtContent)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	testCases := []struct {
		url       string
		userAgent string
		allowed   bool
		desc      string
	}{
		{"https://example.com/my?project_id=123", "*", false, "should match /my?project_id=*"},
		{"https://example.com/my?project_id=456&other=value", "*", false, "should match /my?project_id=* with additional params"},
		{"https://example.com/login?ok_url=test", "*", false, "should match /login?ok_url=*"},
		{"https://example.com/loginWithPwd?ok_url=test123", "*", false, "should match /loginWithPwd?ok_url=*"},
		{"https://example.com/loginWithPw?ok_url=test", "*", false, "should match /loginWithPw?ok_url=*"},
		{"https://example.com/loginWithPw?q=value&ok_url=123", "*", false, "should match /loginWithPw?ok_url=* (ok_url param exists)"},
		{"https://example.com/loginWithPw?ok_url=123&q=value", "*", false, "should match /loginWithPw?ok_url=* (ok_url param at start)"},
		{"https://example.com/loginWithPw?q=qwe1232131ok_url=123&q=qwe123213", "*", false, "should match /loginWithPw?ok_url=* (ok_url param in middle)"},
		{"https://example.com/loginWithPw?q=value", "*", true, "should not match (no ok_url param)"},
		{"https://example.com/users", "*", false, "should match /users"},
		{"https://example.com/users/123", "*", false, "should match /users (prefix match)"},
		{"https://example.com/search/", "*", false, "should match /search/"},
		{"https://example.com/search/query", "*", false, "should match /search/ (prefix match)"},
		{"https://example.com/my", "*", true, "should not match (no query string)"},
		{"https://example.com/login", "*", true, "should not match (no query string)"},
	}

	for _, tc := range testCases {
		result := robots.IsAllowed(tc.url, tc.userAgent)
		if result != tc.allowed {
			t.Errorf("%s: Expected %v, got %v for URL %s",
				tc.desc, tc.allowed, result, tc.url)
		}
	}
}

// TestQueryParameterParsing 测试查询参数解析
// TestQueryParameterParsing tests query parameter parsing
func TestQueryParameterParsing(t *testing.T) {
	robots := NewRobotsTxt()

	testCases := []struct {
		query    string
		expected map[string]string
		desc     string
	}{
		{
			"ok_url=123",
			map[string]string{"ok_url": "123"},
			"single parameter",
		},
		{
			"ok_url=123&q=value",
			map[string]string{"ok_url": "123", "q": "value"},
			"multiple parameters",
		},
		{
			"q=value&ok_url=123",
			map[string]string{"q": "value", "ok_url": "123"},
			"parameters in different order",
		},
		{
			"q=qwe1232131ok_url=123&q=qwe123213",
			map[string]string{"q": "qwe1232131ok_url=123"},
			"complex query string (first q value contains ok_url=123)",
		},
		{
			"ok_url=*",
			map[string]string{"ok_url": "*"},
			"wildcard value",
		},
	}

	for _, tc := range testCases {
		result := robots.parseQueryString(tc.query)
		if len(result) != len(tc.expected) {
			t.Errorf("%s: Expected %d params, got %d", tc.desc, len(tc.expected), len(result))
			continue
		}
		for key, expectedValue := range tc.expected {
			if actualValue, exists := result[key]; !exists {
				t.Errorf("%s: Expected param %s, but not found", tc.desc, key)
			} else if actualValue != expectedValue {
				t.Errorf("%s: Expected param %s=%s, got %s", tc.desc, key, expectedValue, actualValue)
			}
		}
	}
}

// TestWildcardMatching 测试通配符匹配
// TestWildcardMatching tests wildcard matching
func TestWildcardMatching(t *testing.T) {
	robots := NewRobotsTxt()

	testCases := []struct {
		path    string
		pattern string
		matches bool
		desc    string
	}{
		{"/admin/", "/admin/", true, "exact match"},
		{"/admin/test", "/admin/", true, "prefix match"},
		{"/public/", "/admin/", false, "no match"},
		{"/", "/", true, "root path matches root pattern"},
		{"/any/path", "/", true, "any path matches root pattern"},
		{"/test/path", "/test/*", true, "wildcard at end"},
		{"/test/123/path", "/test/*/path", true, "wildcard in middle"},
		{"/test/abc/def/path", "/test/*/path", true, "wildcard matches multiple segments"},
		{"/test/path", "/test/*/path", false, "wildcard requires segment"},
		{"/anything", "*", true, "single * matches all"},
		{"/test/*/path", "/test/*/path", true, "wildcard in pattern matches wildcard in path"},
	}

	for _, tc := range testCases {
		result := robots.matchesPathPattern(tc.path, tc.pattern)
		if result != tc.matches {
			t.Errorf("%s: Expected %v, got %v for path '%s' and pattern '%s'",
				tc.desc, tc.matches, result, tc.path, tc.pattern)
		}
	}
}

// TestDomainMatching 测试域名匹配
// TestDomainMatching tests domain matching
func TestDomainMatching(t *testing.T) {
	robotsTxtContent := `
User-agent: *
Disallow: /admin/
Disallow: /private/
`

	robots := NewRobotsTxt()
	reader := strings.NewReader(robotsTxtContent)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// 设置源 URL
	// Set source URL
	robots.SourceURL = "https://example.com/robots.txt"

	testCases := []struct {
		url       string
		userAgent string
		allowed   bool
		desc      string
	}{
		// 相同域名应该应用规则
		// Same domain should apply rules
		{"https://example.com/admin/", "*", false, "same domain, should be disallowed"},
		{"https://example.com/private/", "*", false, "same domain, should be disallowed"},
		{"https://example.com/public/", "*", true, "same domain, should be allowed"},

		// 不同域名应该默认允许（robots.txt 不适用于其他域名）
		// Different domain should default to allowed (robots.txt doesn't apply to other domains)
		{"https://other.com/admin/", "*", true, "different domain, should be allowed"},
		{"https://other.com/private/", "*", true, "different domain, should be allowed"},

		// 不同协议应该默认允许
		// Different protocol should default to allowed
		{"http://example.com/admin/", "*", true, "different protocol, should be allowed"},

		// 不同端口应该默认允许
		// Different port should default to allowed
		{"https://example.com:8080/admin/", "*", true, "different port, should be allowed"},

		// 子域名应该默认允许（严格来说，robots.txt 不适用于子域名）
		// Subdomain should default to allowed (strictly, robots.txt doesn't apply to subdomains)
		{"https://www.example.com/admin/", "*", true, "subdomain, should be allowed"},
		{"https://m.example.com/admin/", "*", true, "subdomain, should be allowed"},
	}

	for _, tc := range testCases {
		result := robots.IsAllowed(tc.url, tc.userAgent)
		if result != tc.allowed {
			t.Errorf("%s: Expected %v, got %v for URL %s",
				tc.desc, tc.allowed, result, tc.url)
		}
	}
}

// TestDomainMatchingWithoutSourceURL 测试没有源 URL 的情况
// TestDomainMatchingWithoutSourceURL tests case without source URL
func TestDomainMatchingWithoutSourceURL(t *testing.T) {
	robotsTxtContent := `
User-agent: *
Disallow: /admin/
`

	robots := NewRobotsTxt()
	reader := strings.NewReader(robotsTxtContent)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	// 不设置源 URL（手动解析的情况）
	// Don't set source URL (manual parsing case)
	// 这种情况下应该应用规则（因为没有域名信息，无法验证）
	// In this case, rules should apply (cannot verify without domain info)
	result := robots.IsAllowed("https://example.com/admin/", "*")
	if result != false {
		t.Errorf("Expected false, got %v (should apply rules when no source URL)", result)
	}
}

// TestSpecificity 测试规则特异性（更具体的规则优先）
// TestSpecificity tests rule specificity (more specific rules take precedence)
func TestSpecificity(t *testing.T) {
	robotsTxtContent := `
User-agent: *
Disallow: /api/
Allow: /api/public/
`

	robots := NewRobotsTxt()
	reader := strings.NewReader(robotsTxtContent)
	err := robots.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse robots.txt: %v", err)
	}

	testCases := []struct {
		url       string
		userAgent string
		allowed   bool
		desc      string
	}{
		{"https://example.com/api/", "*", false, "should be disallowed by /api/"},
		{"https://example.com/api/private/", "*", false, "should be disallowed by /api/"},
		{"https://example.com/api/public/", "*", true, "should be allowed by more specific /api/public/"},
		{"https://example.com/api/public/data", "*", true, "should be allowed by more specific /api/public/"},
	}

	for _, tc := range testCases {
		result := robots.IsAllowed(tc.url, tc.userAgent)
		if result != tc.allowed {
			t.Errorf("%s: Expected %v, got %v for URL %s",
				tc.desc, tc.allowed, result, tc.url)
		}
	}
}
