# Robots.txt 扫描器

Robots.txt 扫描器插件，用于解析和检查网站的 robots.txt 文件，判断爬虫是否被允许访问特定 URL。

## 功能特性

- ✅ 解析 robots.txt 文件
- ✅ 检查 URL 是否被允许访问
- ✅ 支持多个 User-Agent 规则
- ✅ 获取爬取延迟（Crawl-delay）
- ✅ 获取 Sitemap 信息
- ✅ 支持通配符路径匹配
- ✅ 支持使用自定义 Client 进行请求
- ✅ **支持自定义 robots.txt 路径**（新增）
- ✅ **自动尝试多个常见路径**（新增）
- ✅ **支持直接使用完整 URL**（新增）

## 使用方法

### 基本用法

```go
package main

import (
    "fmt"
    "time"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    // 从 URL 获取 robots.txt
    robotsTxt, err := robots.FetchFromURL("https://example.com", 10*time.Second)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 检查 URL 是否被允许
    allowed := robotsTxt.IsAllowed("https://example.com/admin/", "MyBot/1.0")
    fmt.Printf("Allowed: %v\n", allowed)

    // 获取 Sitemap
    sitemaps := robotsTxt.GetSitemaps()
    fmt.Printf("Sitemaps: %v\n", sitemaps)

    // 获取爬取延迟
    delay := robotsTxt.GetCrawlDelay("MyBot/1.0")
    fmt.Printf("Crawl delay: %d seconds\n", delay)
}
```

### 使用自定义 Client

```go
package main

import (
    "fmt"
    "time"
    "github.com/karosown/katool-go/web_crawler"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    // 创建自定义客户端
    client := web_crawler.NewClient(nil)

    // 使用客户端获取 robots.txt
    robotsTxt, err := robots.FetchFromClient(client, "https://example.com", 10*time.Second)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 检查 URL 是否被允许
    allowed := robotsTxt.IsAllowed("https://example.com/test/", "MyBot/1.0")
    fmt.Printf("Allowed: %v\n", allowed)
}
```

### 使用自定义路径

```go
package main

import (
    "fmt"
    "time"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    // 方式1: 指定自定义路径（相对于网站根目录）
    // Method 1: Specify custom path (relative to website root)
    robotsTxt, err := robots.FetchFromURL("https://example.com", 10*time.Second, "/custom/robots.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 方式2: 直接使用完整 URL
    // Method 2: Use full URL directly
    robotsTxt, err = robots.FetchFromDirectURL("https://example.com/custom/path/robots.txt", 10*time.Second)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    allowed := robotsTxt.IsAllowed("https://example.com/test/", "MyBot/1.0")
    fmt.Printf("Allowed: %v\n", allowed)
}
```

### 自动尝试多个路径

```go
package main

import (
    "fmt"
    "time"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    // 自动尝试多个常见路径：/robots.txt, /robots, /ROBOTS.TXT, /Robots.txt
    // 还可以添加自定义路径
    // Automatically try multiple common paths: /robots.txt, /robots, /ROBOTS.TXT, /Robots.txt
    // Can also add custom paths
    robotsTxt, err := robots.FetchFromURLWithAutoDetect(
        "https://example.com", 
        10*time.Second,
        "/custom/robots",  // 自定义路径 / Custom path
        "/another/path/robots.txt",
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    allowed := robotsTxt.IsAllowed("https://example.com/test/", "MyBot/1.0")
    fmt.Printf("Allowed: %v\n", allowed)
}
```

### 使用自定义 Client 和自定义路径

```go
package main

import (
    "fmt"
    "time"
    "github.com/karosown/katool-go/web_crawler"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    client := web_crawler.NewClient(nil)

    // 使用自定义路径
    // Use custom path
    robotsTxt, err := robots.FetchFromClient(client, "https://example.com", 10*time.Second, "/custom/robots.txt")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 或者自动尝试多个路径
    // Or automatically try multiple paths
    robotsTxt, err = robots.FetchFromClientWithAutoDetect(client, "https://example.com", 10*time.Second)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    allowed := robotsTxt.IsAllowed("https://example.com/test/", "MyBot/1.0")
    fmt.Printf("Allowed: %v\n", allowed)
}
```

### 手动解析 robots.txt 内容

```go
package main

import (
    "fmt"
    "strings"
    "github.com/karosown/katool-go/web_crawler/robots"
)

func main() {
    robotsTxtContent := `
User-agent: *
Disallow: /admin/
Allow: /public/

Sitemap: https://example.com/sitemap.xml
`

    robots := robots.NewRobotsTxt()
    reader := strings.NewReader(robotsTxtContent)
    err := robots.Parse(reader)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // 使用解析后的 robots.txt
    allowed := robots.IsAllowed("https://example.com/admin/", "*")
    fmt.Printf("Allowed: %v\n", allowed)
}
```

## API 文档

### 类型

#### RobotsTxt

Robots.txt 解析器结构体。

```go
type RobotsTxt struct {
    UserAgents map[string]*UserAgentGroup
    Sitemaps   []string
    error      error
}
```

#### UserAgentGroup

用户代理组，包含允许和禁止的路径。

```go
type UserAgentGroup struct {
    UserAgent  string
    Allow      []string
    Disallow   []string
    CrawlDelay int
}
```

### 函数

#### NewRobotsTxt

创建新的 robots.txt 解析器。

```go
func NewRobotsTxt() *RobotsTxt
```

#### FetchFromURL

从指定 URL 获取并解析 robots.txt。支持自定义路径（可选参数）。

```go
func FetchFromURL(baseURL string, timeout time.Duration, robotsPath ...string) (*RobotsTxt, error)
```

**参数说明：**
- `baseURL`: 网站基础 URL
- `timeout`: 请求超时时间
- `robotsPath`: 可选的 robots.txt 路径，如果为空则默认使用 `/robots.txt`

**示例：**
```go
// 使用默认路径 /robots.txt
robots, err := robots.FetchFromURL("https://example.com", 10*time.Second)

// 使用自定义路径
robots, err := robots.FetchFromURL("https://example.com", 10*time.Second, "/custom/robots.txt")

// 使用完整 URL
robots, err := robots.FetchFromURL("https://example.com", 10*time.Second, "https://example.com/custom/robots.txt")
```

#### FetchFromURLWithAutoDetect

自动尝试多个常见路径获取 robots.txt。

```go
func FetchFromURLWithAutoDetect(baseURL string, timeout time.Duration, customPaths ...string) (*RobotsTxt, error)
```

**默认尝试的路径：**
- `/robots.txt`
- `/robots`
- `/ROBOTS.TXT`
- `/Robots.txt`

**示例：**
```go
// 只使用默认路径
robots, err := robots.FetchFromURLWithAutoDetect("https://example.com", 10*time.Second)

// 添加自定义路径
robots, err := robots.FetchFromURLWithAutoDetect("https://example.com", 10*time.Second, "/custom/robots")
```

#### FetchFromDirectURL

直接从完整 URL 获取并解析 robots.txt。

```go
func FetchFromDirectURL(robotsURL string, timeout time.Duration) (*RobotsTxt, error)
```

**示例：**
```go
robots, err := robots.FetchFromDirectURL("https://example.com/custom/path/robots.txt", 10*time.Second)
```

#### FetchFromClient

使用指定的 Client 从 URL 获取并解析 robots.txt。支持自定义路径（可选参数）。

```go
func FetchFromClient(client *web_crawler.Client, baseURL string, timeout time.Duration, robotsPath ...string) (*RobotsTxt, error)
```

**参数说明：**
- `client`: web_crawler.Client 实例
- `baseURL`: 网站基础 URL
- `timeout`: 请求超时时间
- `robotsPath`: 可选的 robots.txt 路径，如果为空则默认使用 `/robots.txt`

**示例：**
```go
client := web_crawler.NewClient(nil)

// 使用默认路径
robots, err := robots.FetchFromClient(client, "https://example.com", 10*time.Second)

// 使用自定义路径
robots, err := robots.FetchFromClient(client, "https://example.com", 10*time.Second, "/custom/robots.txt")
```

#### FetchFromClientWithAutoDetect

使用指定的 Client 自动尝试多个常见路径获取 robots.txt。

```go
func FetchFromClientWithAutoDetect(client *web_crawler.Client, baseURL string, timeout time.Duration, customPaths ...string) (*RobotsTxt, error)
```

**示例：**
```go
client := web_crawler.NewClient(nil)
robots, err := robots.FetchFromClientWithAutoDetect(client, "https://example.com", 10*time.Second, "/custom/robots")
```

#### FetchFromClientDirectURL

使用指定的 Client 直接从完整 URL 获取并解析 robots.txt。

```go
func FetchFromClientDirectURL(client *web_crawler.Client, robotsURL string, timeout time.Duration) (*RobotsTxt, error)
```

**示例：**
```go
client := web_crawler.NewClient(nil)
robots, err := robots.FetchFromClientDirectURL(client, "https://example.com/custom/path/robots.txt", 10*time.Second)
```

### 方法

#### Parse

解析 robots.txt 内容。

```go
func (r *RobotsTxt) Parse(reader io.Reader) error
```

#### IsAllowed

检查指定 URL 是否被允许访问（针对特定 User-Agent）。

```go
func (r *RobotsTxt) IsAllowed(url, userAgent string) bool
```

#### GetCrawlDelay

获取指定 User-Agent 的爬取延迟。

```go
func (r *RobotsTxt) GetCrawlDelay(userAgent string) int
```

#### GetSitemaps

获取所有 sitemap URL。

```go
func (r *RobotsTxt) GetSitemaps() []string
```

#### GetDisallowedPaths

获取指定 User-Agent 禁止的路径列表。

```go
func (r *RobotsTxt) GetDisallowedPaths(userAgent string) []string
```

#### GetAllowedPaths

获取指定 User-Agent 允许的路径列表。

```go
func (r *RobotsTxt) GetAllowedPaths(userAgent string) []string
```

## 支持的 robots.txt 指令

- `User-agent`: 指定用户代理
- `Disallow`: 禁止访问的路径
- `Allow`: 允许访问的路径
- `Crawl-delay`: 爬取延迟（秒）
- `Sitemap`: Sitemap URL

## 路径匹配规则

1. 精确匹配：路径必须完全匹配模式
2. 前缀匹配：路径以模式开头
3. 通配符匹配：支持 `*` 通配符
4. 最长匹配：如果有多个规则匹配，选择最具体的规则

## 自定义路径支持

本插件支持多种方式指定 robots.txt 的路径：

1. **默认路径**：如果不指定路径，默认使用 `/robots.txt`
2. **相对路径**：可以指定相对于网站根目录的路径，如 `/custom/robots.txt`
3. **完整 URL**：可以直接使用完整的 robots.txt URL
4. **自动检测**：可以自动尝试多个常见路径，直到找到可用的文件

**常见非标准路径：**
- `/robots`（无扩展名）
- `/ROBOTS.TXT`（大写）
- `/Robots.txt`（混合大小写）
- `/custom/robots.txt`（自定义目录）

## 注意事项

- 如果没有匹配的 User-Agent 规则，默认允许访问
- 如果 robots.txt 文件不存在或无法访问，需要自行处理错误
- 建议在爬取网站前先检查 robots.txt，遵守网站的爬取规则
- 对于使用非标准路径的网站，可以使用 `FetchFromURLWithAutoDetect` 或 `FetchFromClientWithAutoDetect` 自动尝试多个路径
- 如果知道确切的路径，建议直接指定路径以提高效率

