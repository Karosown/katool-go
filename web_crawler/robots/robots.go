package robots

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"github.com/karosown/katool-go/web_crawler"
)

// RobotsTxt robots.txt 解析器结构体
// RobotsTxt represents a robots.txt parser structure
type RobotsTxt struct {
	UserAgents map[string]*UserAgentGroup // 用户代理组映射 / User agent group mapping
	Sitemaps   []string                   // Sitemap 列表 / Sitemap list
	SourceURL  string                     // robots.txt 的源 URL（用于域名验证）/ Source URL of robots.txt (for domain validation)
	error      error                      // 错误信息 / Error information
}

// UserAgentGroup 用户代理组，包含允许和禁止的路径
// UserAgentGroup represents a user agent group with allowed and disallowed paths
type UserAgentGroup struct {
	UserAgent  string   // 用户代理 / User agent
	Allow      []string // 允许的路径 / Allowed paths
	Disallow   []string // 禁止的路径 / Disallowed paths
	CrawlDelay int      // 爬取延迟（秒）/ Crawl delay (seconds)
}

// IsErr 检查是否有错误
// IsErr checks if there is an error
func (r *RobotsTxt) IsErr() bool {
	return r.error != nil
}

// Error 返回错误信息
// Error returns error information
func (r *RobotsTxt) Error() error {
	return r.error
}

// SetErr 设置错误信息
// SetErr sets error information
func (r *RobotsTxt) SetErr(err error) *RobotsTxt {
	r.error = err
	return r
}

// NewRobotsTxt 创建新的 robots.txt 解析器
// NewRobotsTxt creates a new robots.txt parser
func NewRobotsTxt() *RobotsTxt {
	return &RobotsTxt{
		UserAgents: make(map[string]*UserAgentGroup),
		Sitemaps:   make([]string, 0),
		SourceURL:  "",
	}
}

// FetchFromURL 从指定 URL 获取并解析 robots.txt
// 如果 robotsPath 为空，默认使用 /robots.txt
// FetchFromURL fetches and parses robots.txt from the specified URL
// If robotsPath is empty, defaults to /robots.txt
func FetchFromURL(baseURL string, timeout time.Duration, robotsPath ...string) (*RobotsTxt, error) {
	parsedURL, err := nurl.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// 确定 robots.txt 路径
	// Determine robots.txt path
	path := "/robots.txt"
	if len(robotsPath) > 0 && robotsPath[0] != "" {
		path = robotsPath[0]
		// 如果路径不是以 / 开头且不是完整 URL，则添加 /
		// If path doesn't start with / and is not a full URL, add /
		if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
			path = "/" + path
		}
	}

	// 构建 robots.txt URL
	// Build robots.txt URL
	var robotsURL string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// 如果已经是完整 URL，直接使用
		// If it's already a full URL, use it directly
		robotsURL = path
	} else {
		robotsURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)
	}

	return fetchRobotsTxt(robotsURL, nil, timeout)
}

// FetchFromURLWithAutoDetect 自动尝试多个常见路径获取 robots.txt
// FetchFromURLWithAutoDetect automatically tries multiple common paths to fetch robots.txt
func FetchFromURLWithAutoDetect(baseURL string, timeout time.Duration, customPaths ...string) (*RobotsTxt, error) {
	parsedURL, err := nurl.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// 默认尝试的路径列表
	// Default list of paths to try
	defaultPaths := []string{
		"/robots.txt",
		"/robots",
		"/ROBOTS.TXT",
		"/Robots.txt",
	}

	// 合并自定义路径
	// Merge custom paths
	paths := append(defaultPaths, customPaths...)

	// 尝试每个路径
	// Try each path
	var lastErr error
	for _, path := range paths {
		robotsURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)
		robots, err := fetchRobotsTxt(robotsURL, nil, timeout)
		if err == nil {
			return robots, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("failed to fetch robots.txt from any path, last error: %v", lastErr)
}

// FetchFromDirectURL 直接从完整 URL 获取并解析 robots.txt
// FetchFromDirectURL fetches and parses robots.txt directly from a full URL
func FetchFromDirectURL(robotsURL string, timeout time.Duration) (*RobotsTxt, error) {
	return fetchRobotsTxt(robotsURL, nil, timeout)
}

// FetchFromClient 使用指定的 Client 从 URL 获取并解析 robots.txt
// 如果 robotsPath 为空，默认使用 /robots.txt
// FetchFromClient fetches and parses robots.txt using the specified Client
// If robotsPath is empty, defaults to /robots.txt
func FetchFromClient(client *web_crawler.Client, baseURL string, timeout time.Duration, robotsPath ...string) (*RobotsTxt, error) {
	parsedURL, err := nurl.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// 确定 robots.txt 路径
	// Determine robots.txt path
	path := "/robots.txt"
	if len(robotsPath) > 0 && robotsPath[0] != "" {
		path = robotsPath[0]
		// 如果路径不是以 / 开头且不是完整 URL，则添加 /
		// If path doesn't start with / and is not a full URL, add /
		if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
			path = "/" + path
		}
	}

	// 构建 robots.txt URL
	// Build robots.txt URL
	var robotsURL string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// 如果已经是完整 URL，直接使用
		// If it's already a full URL, use it directly
		robotsURL = path
	} else {
		robotsURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)
	}

	return fetchRobotsTxt(robotsURL, client, timeout)
}

// FetchFromClientWithAutoDetect 使用指定的 Client 自动尝试多个常见路径获取 robots.txt
// FetchFromClientWithAutoDetect uses the specified Client to automatically try multiple common paths to fetch robots.txt
func FetchFromClientWithAutoDetect(client *web_crawler.Client, baseURL string, timeout time.Duration, customPaths ...string) (*RobotsTxt, error) {
	parsedURL, err := nurl.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// 默认尝试的路径列表
	// Default list of paths to try
	defaultPaths := []string{
		"/robots.txt",
		"/robots",
		"/ROBOTS.TXT",
		"/Robots.txt",
	}

	// 合并自定义路径
	// Merge custom paths
	paths := append(defaultPaths, customPaths...)

	// 尝试每个路径
	// Try each path
	var lastErr error
	for _, path := range paths {
		robotsURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)
		robots, err := fetchRobotsTxt(robotsURL, client, timeout)
		if err == nil {
			return robots, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("failed to fetch robots.txt from any path, last error: %v", lastErr)
}

// FetchFromClientDirectURL 使用指定的 Client 直接从完整 URL 获取并解析 robots.txt
// FetchFromClientDirectURL uses the specified Client to fetch and parse robots.txt directly from a full URL
func FetchFromClientDirectURL(client *web_crawler.Client, robotsURL string, timeout time.Duration) (*RobotsTxt, error) {
	return fetchRobotsTxt(robotsURL, client, timeout)
}

// fetchRobotsTxt 内部函数，用于实际获取 robots.txt
// fetchRobotsTxt is an internal function for actually fetching robots.txt
func fetchRobotsTxt(robotsURL string, client *web_crawler.Client, timeout time.Duration) (*RobotsTxt, error) {
	var httpClient *http.Client

	if client != nil {
		// 使用指定的 Client
		// Use specified Client
		httpClient = web_crawler.NewAntiBotClientWithConfig(
			client.GetChrome(),
			timeout,
			client.GetConfig(),
			client.GetPolicies()...,
		)
	} else {
		// 使用默认客户端
		// Use default client
		httpClient = web_crawler.NewAntiBotClient(nil, timeout)
	}

	req, err := http.NewRequest("GET", robotsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置 User-Agent
	// Set User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RobotsTxtBot/1.0)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch robots.txt: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("robots.txt returned status code: %d", resp.StatusCode)
	}

	// 解析 robots.txt
	// Parse robots.txt
	robots := NewRobotsTxt()
	robots.SourceURL = robotsURL // 保存源 URL / Save source URL
	err = robots.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse robots.txt: %v", err)
	}

	return robots, nil
}

// Parse 解析 robots.txt 内容
// Parse parses robots.txt content
func (r *RobotsTxt) Parse(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	var currentGroup *UserAgentGroup

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析指令
		// Parse directives
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		directive := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch directive {
		case "user-agent":
			// 创建新的用户代理组
			// Create new user agent group
			currentGroup = &UserAgentGroup{
				UserAgent:  value,
				Allow:      make([]string, 0),
				Disallow:   make([]string, 0),
				CrawlDelay: 0,
			}
			r.UserAgents[value] = currentGroup

		case "disallow":
			if currentGroup != nil {
				currentGroup.Disallow = append(currentGroup.Disallow, value)
			}

		case "allow":
			if currentGroup != nil {
				currentGroup.Allow = append(currentGroup.Allow, value)
			}

		case "crawl-delay":
			if currentGroup != nil {
				var delay int
				fmt.Sscanf(value, "%d", &delay)
				currentGroup.CrawlDelay = delay
			}

		case "sitemap":
			r.Sitemaps = append(r.Sitemaps, value)
		}
	}

	return scanner.Err()
}

// IsAllowed 检查指定 URL 是否被允许访问（针对特定 User-Agent）
// IsAllowed checks if the specified URL is allowed to be accessed (for a specific User-Agent)
func (r *RobotsTxt) IsAllowed(url, userAgent string) bool {
	parsedURL, err := nurl.Parse(url)
	if err != nil {
		return true // 如果 URL 解析失败，默认允许 / If URL parsing fails, default to allowed
	}

	// 检查域名是否匹配（robots.txt 只适用于其所在域名）
	// Check if domain matches (robots.txt only applies to its own domain)
	if !r.isSameDomain(parsedURL) {
		// 如果域名不匹配，返回 true（默认允许），因为该 robots.txt 不适用于其他域名
		// If domain doesn't match, return true (default allowed), because this robots.txt doesn't apply to other domains
		return true
	}

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	// 构建完整的路径（包含查询字符串，如果规则中包含查询字符串）
	// Build full path (including query string if rule contains query string)
	fullPath := path
	if parsedURL.RawQuery != "" {
		fullPath = path + "?" + parsedURL.RawQuery
	}

	// 获取匹配的用户代理组
	// Get matching user agent group
	group := r.getMatchingGroup(userAgent)
	if group == nil {
		return true // 如果没有匹配的规则，默认允许 / If no matching rules, default to allowed
	}

	// 找到所有匹配的规则，并确定最具体的规则
	// Find all matching rules and determine the most specific rule
	var bestMatch *ruleMatch

	// 检查所有 Disallow 规则
	// Check all Disallow rules
	for _, disallowPath := range group.Disallow {
		if match := r.matchRule(fullPath, path, disallowPath); match != nil {
			if bestMatch == nil || match.specificity > bestMatch.specificity {
				bestMatch = &ruleMatch{
					path:        disallowPath,
					allowed:     false,
					specificity: match.specificity,
				}
			} else if match.specificity == bestMatch.specificity && !bestMatch.allowed {
				// 如果特异性相同，Disallow 优先于 Allow
				// If specificity is the same, Disallow takes precedence over Allow
				bestMatch = &ruleMatch{
					path:        disallowPath,
					allowed:     false,
					specificity: match.specificity,
				}
			}
		}
	}

	// 检查所有 Allow 规则
	// Check all Allow rules
	for _, allowPath := range group.Allow {
		if match := r.matchRule(fullPath, path, allowPath); match != nil {
			if bestMatch == nil || match.specificity > bestMatch.specificity {
				bestMatch = &ruleMatch{
					path:        allowPath,
					allowed:     true,
					specificity: match.specificity,
				}
			} else if match.specificity == bestMatch.specificity && bestMatch.allowed {
				// 如果特异性相同且都是 Allow，保留更具体的
				// If specificity is the same and both are Allow, keep the more specific one
				bestMatch = &ruleMatch{
					path:        allowPath,
					allowed:     true,
					specificity: match.specificity,
				}
			}
		}
	}

	// 如果没有匹配的规则，默认允许
	// If no matching rules, default to allowed
	if bestMatch == nil {
		return true
	}

	return bestMatch.allowed
}

// ruleMatch 规则匹配结果
// ruleMatch represents a rule match result
type ruleMatch struct {
	path        string
	allowed     bool
	specificity int // 特异性，值越大越具体 / Specificity, higher value means more specific
}

// matchRule 检查路径是否匹配规则，返回匹配信息和特异性
// matchRule checks if path matches rule, returns match info and specificity
func (r *RobotsTxt) matchRule(fullPath, path, pattern string) *ruleMatch {
	if pattern == "" {
		return nil
	}

	// 如果模式是 "/"，匹配所有路径
	// If pattern is "/", match all paths
	if pattern == "/" {
		return &ruleMatch{
			path:        pattern,
			allowed:     false,
			specificity: 1,
		}
	}

	// 检查规则是否包含查询字符串
	// Check if rule contains query string
	patternHasQuery := strings.Contains(pattern, "?")
	var patternPath, patternQuery string
	if patternHasQuery {
		parts := strings.SplitN(pattern, "?", 2)
		patternPath = parts[0]
		patternQuery = parts[1]
	} else {
		patternPath = pattern
	}

	// 首先检查路径部分是否匹配
	// First check if path part matches
	if !r.matchesPathPattern(path, patternPath) {
		return nil
	}

	// 如果规则包含查询字符串，也需要匹配查询字符串
	// If rule contains query string, also need to match query string
	if patternHasQuery {
		// 从 fullPath 中提取查询字符串
		// Extract query string from fullPath
		if !strings.Contains(fullPath, "?") {
			return nil
		}
		parts := strings.SplitN(fullPath, "?", 2)
		actualQuery := parts[1]

		// 匹配查询字符串（支持通配符）
		// Match query string (supports wildcards)
		if !r.matchesQueryPattern(actualQuery, patternQuery) {
			return nil
		}
	}

	// 计算特异性：路径长度 + 查询字符串长度（如果存在）
	// Calculate specificity: path length + query string length (if exists)
	specificity := len(patternPath)
	if patternHasQuery {
		specificity += len(patternQuery) * 10 // 查询字符串的权重更高 / Query string has higher weight
	}

	return &ruleMatch{
		path:        pattern,
		allowed:     false, // 这个值会在调用处设置 / This value will be set at the call site
		specificity: specificity,
	}
}

// matchesQueryPattern 检查查询字符串是否匹配模式（支持通配符和参数匹配）
// matchesQueryPattern checks if query string matches pattern (supports wildcards and parameter matching)
func (r *RobotsTxt) matchesQueryPattern(query, pattern string) bool {
	// 如果模式是 "*"，匹配所有查询字符串
	// If pattern is "*", match all query strings
	if pattern == "*" {
		return true
	}

	// 如果模式包含 "="，说明是参数匹配（如 "ok_url=*" 或 "ok_url=123"）
	// If pattern contains "=", it's parameter matching (e.g., "ok_url=*" or "ok_url=123")
	if strings.Contains(pattern, "=") {
		return r.matchesQueryParameter(query, pattern)
	}

	// 支持通配符匹配
	// Support wildcard matching
	if strings.Contains(pattern, "*") {
		return r.matchesWildcard(query, pattern)
	}

	// 检查查询字符串中是否包含该模式（不限于前缀）
	// Check if query string contains the pattern (not limited to prefix)
	return strings.Contains(query, pattern)
}

// getMatchingGroup 获取匹配的用户代理组
// getMatchingGroup gets the matching user agent group
func (r *RobotsTxt) getMatchingGroup(userAgent string) *UserAgentGroup {
	// 首先尝试精确匹配
	// First try exact match
	if group, ok := r.UserAgents[userAgent]; ok {
		return group
	}

	// 尝试匹配 "*"（所有用户代理）
	// Try to match "*" (all user agents)
	if group, ok := r.UserAgents["*"]; ok {
		return group
	}

	// 尝试部分匹配
	// Try partial match
	for ua, group := range r.UserAgents {
		if strings.Contains(userAgent, ua) || strings.Contains(ua, userAgent) {
			return group
		}
	}

	return nil
}

// isSameDomain 检查 URL 是否与 robots.txt 的源域名匹配
// isSameDomain checks if URL matches the source domain of robots.txt
func (r *RobotsTxt) isSameDomain(url *nurl.URL) bool {
	// 如果没有源 URL，无法验证，默认允许
	// If no source URL, cannot verify, default to allowed
	if r.SourceURL == "" {
		return true
	}

	// 解析源 URL
	// Parse source URL
	sourceURL, err := nurl.Parse(r.SourceURL)
	if err != nil {
		return true // 如果解析失败，默认允许 / If parsing fails, default to allowed
	}

	// 比较协议、主机和端口
	// Compare protocol, host, and port
	// 注意：robots.txt 只适用于相同的协议、主机和端口
	// Note: robots.txt only applies to the same protocol, host, and port
	if sourceURL.Scheme != url.Scheme {
		return false
	}

	// 规范化主机名（去除 www 前缀进行比较，但严格来说应该完全匹配）
	// Normalize hostname (remove www prefix for comparison, but strictly should be exact match)
	sourceHost := r.normalizeHost(sourceURL.Host)
	targetHost := r.normalizeHost(url.Host)

	// 严格匹配：协议、主机和端口必须完全相同
	// Strict match: protocol, host, and port must be exactly the same
	// 但为了更灵活，我们也可以比较规范化后的主机名
	// But for flexibility, we can also compare normalized hostnames
	return sourceHost == targetHost
}

// normalizeHost 规范化主机名（提取主机和端口）
// normalizeHost normalizes hostname (extracts host and port)
func (r *RobotsTxt) normalizeHost(host string) string {
	// 分离主机和端口
	// Separate host and port
	hostname := host
	port := ""
	if idx := strings.Index(host, ":"); idx != -1 {
		hostname = host[:idx]
		port = host[idx+1:]
	}

	// 转换为小写（主机名不区分大小写）
	// Convert to lowercase (hostnames are case-insensitive)
	hostname = strings.ToLower(hostname)

	// 如果端口为空，使用默认端口（但这里我们保留原始端口信息）
	// If port is empty, use default port (but here we keep original port info)
	if port != "" {
		return hostname + ":" + port
	}
	return hostname
}

// matchesPathPattern 检查路径是否匹配模式（支持通配符，符合 Google robots.txt 规范）
// matchesPathPattern checks if path matches pattern (supports wildcards, conforms to Google robots.txt spec)
func (r *RobotsTxt) matchesPathPattern(path, pattern string) bool {
	if pattern == "" {
		return false
	}

	// 如果模式是 "/"，匹配所有路径
	// If pattern is "/", match all paths
	if pattern == "/" {
		return true
	}

	// 支持通配符匹配
	// Support wildcard matching
	if strings.Contains(pattern, "*") {
		return r.matchesWildcard(path, pattern)
	}

	// 前缀匹配：路径必须以模式开头
	// Prefix match: path must start with pattern
	// 根据 Google 规范，这是前缀匹配，不是精确匹配
	// According to Google spec, this is prefix matching, not exact matching
	return strings.HasPrefix(path, pattern)
}

// matchesWildcard 使用通配符匹配路径（支持 * 在任意位置，符合 robots.txt 规范）
// matchesWildcard matches path using wildcards (supports * at any position, conforms to robots.txt spec)
func (r *RobotsTxt) matchesWildcard(text, pattern string) bool {
	// 如果模式是 "*"，匹配所有
	// If pattern is "*", match all
	if pattern == "*" {
		return true
	}

	// 使用简单的通配符匹配实现
	// Use simple wildcard matching implementation
	return r.simpleWildcardMatch(text, pattern)
}

// matchesQueryParameter 检查查询字符串中是否包含匹配的参数
// matchesQueryParameter checks if query string contains matching parameter
func (r *RobotsTxt) matchesQueryParameter(query, pattern string) bool {
	// 解析模式：param=value 或 param=*
	// Parse pattern: param=value or param=*
	parts := strings.SplitN(pattern, "=", 2)
	if len(parts) != 2 {
		return false
	}
	paramName := parts[0]
	paramValuePattern := parts[1]

	// 首先尝试标准解析：检查参数是否作为独立参数存在
	// First try standard parsing: check if parameter exists as independent parameter
	queryParams := r.parseQueryString(query)
	if actualValue, exists := queryParams[paramName]; exists {
		// 如果模式值是 "*"，只要参数存在就匹配
		// If pattern value is "*", match if parameter exists
		if paramValuePattern == "*" {
			return true
		}

		// 如果模式值包含通配符，使用通配符匹配
		// If pattern value contains wildcard, use wildcard matching
		if strings.Contains(paramValuePattern, "*") {
			return r.matchesWildcard(actualValue, paramValuePattern)
		}

		// 精确匹配参数值
		// Exact match parameter value
		return actualValue == paramValuePattern
	}

	// 如果标准解析没找到，检查查询字符串中是否包含 "paramName=" 模式
	// If standard parsing didn't find it, check if query string contains "paramName=" pattern
	// 这处理了参数值中包含 "paramName=value" 的情况（如 q=qwe1232131ok_url=123）
	// This handles cases where parameter value contains "paramName=value" (e.g., q=qwe1232131ok_url=123)
	paramPattern := paramName + "="
	if !strings.Contains(query, paramPattern) {
		return false
	}

	// 如果模式值是 "*"，只要找到参数名就匹配
	// If pattern value is "*", match if parameter name is found
	if paramValuePattern == "*" {
		return true
	}

	// 尝试从查询字符串中提取参数值
	// Try to extract parameter value from query string
	// 查找所有 "paramName=" 出现的位置
	// Find all occurrences of "paramName="
	idx := strings.Index(query, paramPattern)
	if idx == -1 {
		return false
	}

	// 从 "paramName=" 之后提取值
	// Extract value after "paramName="
	valueStart := idx + len(paramPattern)
	valueEnd := valueStart

	// 查找值的结束位置（& 或字符串结尾）
	// Find end of value (& or end of string)
	for valueEnd < len(query) && query[valueEnd] != '&' {
		valueEnd++
	}

	actualValue := query[valueStart:valueEnd]

	// 如果模式值包含通配符，使用通配符匹配
	// If pattern value contains wildcard, use wildcard matching
	if strings.Contains(paramValuePattern, "*") {
		return r.matchesWildcard(actualValue, paramValuePattern)
	}

	// 精确匹配参数值
	// Exact match parameter value
	return actualValue == paramValuePattern
}

// parseQueryString 解析查询字符串为参数映射
// parseQueryString parses query string into parameter map
func (r *RobotsTxt) parseQueryString(query string) map[string]string {
	params := make(map[string]string)

	// 分割查询字符串
	// Split query string
	pairs := strings.Split(query, "&")

	for _, pair := range pairs {
		// 分割键值对
		// Split key-value pair
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]
			// 如果同一个参数出现多次，保留第一个值（或可以合并，这里保留第一个）
			// If same parameter appears multiple times, keep first value
			if _, exists := params[key]; !exists {
				params[key] = value
			}
		} else if len(kv) == 1 && kv[0] != "" {
			// 只有键没有值的情况
			// Case where only key exists without value
			params[kv[0]] = ""
		}
	}

	return params
}

// simpleWildcardMatch 简单的通配符匹配实现（符合 robots.txt 规范）
// simpleWildcardMatch simple wildcard matching implementation (conforms to robots.txt spec)
func (r *RobotsTxt) simpleWildcardMatch(text, pattern string) bool {
	// 处理特殊情况
	// Handle special cases
	if pattern == "*" {
		return true
	}
	if pattern == "" {
		return text == ""
	}

	// 将模式分割成多个部分（按 * 分割）
	// Split pattern into multiple parts (split by *)
	parts := strings.Split(pattern, "*")

	// 如果没有 *，直接进行前缀匹配
	// If no *, directly do prefix matching
	if len(parts) == 1 {
		return strings.HasPrefix(text, pattern)
	}

	// 检查第一个部分（必须在开头）
	// Check first part (must be at the beginning)
	if parts[0] != "" {
		if !strings.HasPrefix(text, parts[0]) {
			return false
		}
		text = text[len(parts[0]):]
	}

	// 检查最后一个部分（必须在结尾，如果非空）
	// Check last part (must be at the end, if not empty)
	if len(parts) > 1 && parts[len(parts)-1] != "" {
		lastPart := parts[len(parts)-1]
		if !strings.HasSuffix(text, lastPart) {
			return false
		}
		text = text[:len(text)-len(lastPart)]
		parts = parts[:len(parts)-1]
	}

	// 检查中间部分（必须按顺序出现）
	// Check middle parts (must appear in order)
	textIdx := 0
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		idx := strings.Index(text[textIdx:], parts[i])
		if idx == -1 {
			return false
		}
		textIdx += idx + len(parts[i])
	}

	return true
}

// GetCrawlDelay 获取指定 User-Agent 的爬取延迟
// GetCrawlDelay gets the crawl delay for the specified User-Agent
func (r *RobotsTxt) GetCrawlDelay(userAgent string) int {
	group := r.getMatchingGroup(userAgent)
	if group == nil {
		return 0
	}
	return group.CrawlDelay
}

// GetSitemaps 获取所有 sitemap URL
// GetSitemaps gets all sitemap URLs
func (r *RobotsTxt) GetSitemaps() []string {
	return r.Sitemaps
}

// GetDisallowedPaths 获取指定 User-Agent 禁止的路径列表
// GetDisallowedPaths gets the list of disallowed paths for the specified User-Agent
func (r *RobotsTxt) GetDisallowedPaths(userAgent string) []string {
	group := r.getMatchingGroup(userAgent)
	if group == nil {
		return []string{}
	}
	return group.Disallow
}

// GetAllowedPaths 获取指定 User-Agent 允许的路径列表
// GetAllowedPaths gets the list of allowed paths for the specified User-Agent
func (r *RobotsTxt) GetAllowedPaths(userAgent string) []string {
	group := r.getMatchingGroup(userAgent)
	if group == nil {
		return []string{}
	}
	return group.Allow
}
