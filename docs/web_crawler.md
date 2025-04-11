# Web 爬虫工具

Web 爬虫模块提供了从网页抓取内容、解析文章和处理 Web 内容的工具。该模块结合使用 `go-readability` 和 `go-rod` 等库，支持常规 HTTP 请求和浏览器模拟。

## 主要功能

- **网页内容提取**：从 URL 获取 HTML 内容
- **可读性解析**：提取文章主体内容，过滤广告和导航等无关内容
- **JS 渲染支持**：支持执行 JavaScript 渲染页面后再提取内容
- **错误处理**：内置错误处理机制，简化错误处理流程
- **流式处理**：支持将处理结果转换为流进行后续操作

## 基本用法

### 从 URL 获取文章内容

```go
// 从 URL 获取文章内容
article, err := web_crawler.FromURL("https://example.com/article", 10*time.Second)
if err != nil {
    fmt.Println("获取文章失败:", err)
    return
}

// 访问文章内容
fmt.Println("标题:", article.Title)
fmt.Println("内容:", article.TextContent)
fmt.Println("HTML:", article.Content)
```

### 添加请求自定义选项

```go
// 添加自定义请求头或其他选项
article, err := web_crawler.FromURL("https://example.com/article", 10*time.Second, 
    func(r *http.Request) {
        r.Header.Set("User-Agent", "Mozilla/5.0 ...")
        r.Header.Set("Referer", "https://example.com")
    })
```

### 使用浏览器引擎执行 JavaScript

```go
// 初始化 Chrome 实例（通常在程序启动时）
web_crawler.WebChrome = core.NewContain()

// 执行 JavaScript 并获取结果
result, err := web_crawler.ExecFun(
    "https://example.com/spa", 
    "return document.querySelector('.article-content').textContent",
    func(page *rod.Page) {
        // 可选：等待特定元素出现
        page.MustElement(".article-content").MustWaitVisible()
        // 可选：点击、滚动等操作
    })
```

### 解析相对路径

```go
// 将相对路径转换为绝对 URL
absoluteURL := web_crawler.ParsePath(
    "https://example.com/article/", 
    "../images/photo.jpg")
// 返回: https://example.com/images/photo.jpg
```

### 错误处理

使用内置的错误处理类型：

```go
var webErr web_crawler.WebReaderError
// ... 执行操作 ...

// 检查是否有错误
if webErr.IsErr() {
    fmt.Println("发生错误:", webErr.error)
}

// 将错误添加到现有错误列表
errors := someOtherErrors // 已有的错误
errors = webErr.SolveErrors(errors) // 合并错误
```

### 使用 WebReaderValue 和流处理

```go
// 创建 WebReaderValue
values := web_crawler.NewWebReaderValue("value1", "value2", "value3")

// 转换为流进行处理
results := values.Stream().
    Filter(func(v web_crawler.WebReaderString) bool {
        return strings.HasPrefix(v.String(), "value")
    }).
    Map(func(v web_crawler.WebReaderString) any {
        return "Processed: " + v.String()
    }).
    ToList()
```

## 高级用法

### 批量抓取和处理

```go
// 准备要抓取的 URL 列表
urls := []string{
    "https://example.com/article1",
    "https://example.com/article2",
    "https://example.com/article3",
}

// 转换为流进行批量处理
results := stream.ToStream(&urls).
    Map(func(url string) any {
        article, err := web_crawler.FromURL(url, 10*time.Second)
        if err != nil {
            return nil // 或创建一个错误对象
        }
        // 提取所需内容
        return map[string]string{
            "url":   url,
            "title": article.Title,
            "text":  article.TextContent,
        }
    }).
    Filter(func(result any) bool {
        return result != nil // 过滤失败的请求
    }).
    ToList()
```

### 从页面提取链接并递归抓取

```go
func crawlPage(url string, depth int) []map[string]string {
    if depth <= 0 {
        return nil
    }
    
    results := make([]map[string]string, 0)
    
    // 抓取当前页面
    article, err := web_crawler.FromURL(url, 10*time.Second)
    if err != nil {
        return results
    }
    
    // 保存当前页面结果
    results = append(results, map[string]string{
        "url":   url,
        "title": article.Title,
        "text":  article.TextContent,
    })
    
    // 提取页面中的链接
    links := extractLinks(article.Content) // 使用正则或 HTML 解析提取链接
    
    // 递归抓取链接（处理相对路径）
    for _, link := range links {
        absoluteLink := web_crawler.ParsePath(url, link)
        subResults := crawlPage(absoluteLink, depth-1)
        results = append(results, subResults...)
    }
    
    return results
}
```

## 注意事项

1. 使用爬虫时请遵守网站的 robots.txt 规则和使用条款
2. 添加适当的延迟，避免对目标服务器造成过大负载
3. 对于动态加载内容，使用 `ExecFun` 方法通过浏览器引擎执行
4. 处理大量页面时，考虑使用流和并行处理来提高效率
5. 注意处理 URL 编码和解码，特别是处理含有非 ASCII 字符的 URL 时 


```go
package test

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/karosown/katool-go/sys"
	"github.com/karosown/katool-go/web_crawler"
	"github.com/karosown/katool-go/web_crawler/core"
	"github.com/karosown/katool-go/web_crawler/render"
	"github.com/karosown/katool-go/web_crawler/rss"
)

func TestChorme(t *testing.T) {
	container := launcher.MustResolveURL(web_crawler.ChromeRemoteURL)
	browser := rod.New().ControlURL(container).MustConnect().MustPage("https://openai.com/news/rss.xml")
	fmt.Println(
		browser.MustEval("() => {" +
			"	return document.documentElement.outerHTML" +
			"}"),
	)
}
func TestWebReader(t *testing.T) {
	u := "https://juejin.cn/post/7357922389417017353"
	article := web_crawler.GetArticleWithURL(u, func(r *http.Request) {
		r.Header = http.Header{
			"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3\""},
		}
	})
	if article.IsErr() {
		fmt.Println(article.Error())
	}
	fmt.Printf("URL     : %s\n", u)
	fmt.Printf("Title   : %s\n", article.Title)
	fmt.Printf("Author  : %s\n", article.Byline)
	fmt.Printf("Length  : %d\n", article.Length)
	fmt.Printf("Excerpt : %s\n", article.Excerpt)
	fmt.Printf("S  : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Printf("Content : %s\n", article.Content)
	fmt.Printf("TextContent : %s\n", article.SiteName)
	fmt.Printf("ImageiteName: %s\n", article.TextContent)
	fmt.Println()
}

func TestReadSourceCode(t *testing.T) {
	code := web_crawler.ReadSourceCode("https://openai.com/news/rss.xml", "", func(page *rod.Page) {
		page.MustStopLoading()
	})
	fmt.Println(code.String())
	r := &rss.RSS{}
	err := xml.Unmarshal([]byte(code.String()), r)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(r)
}
func TestReadChrome(t *testing.T) {
	for {
		u := "https://openai.com/index/openai-and-the-csu-system/"
		article := web_crawler.GetArticleWithChrome(u, func(page *rod.Page) {
			page.MustWaitLoad()
		}, nil)
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", u)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Content : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
	}
}
func TestReadClaude(t *testing.T) {
	//for {
	u := "https://www.anthropic.com/news"
	web_crawler.ReadArray(u, `()=>{
			// 使用 querySelectorAll 选择所有匹配的 <a> 标签  
			let linkElements = document.querySelectorAll('a.PostCard_post-card__z_Sqq.PostList_post-card__1g0fm');  
			
			// 将 NodeList 转换为数组并获取所有 href  
			let hrefValues = Array.from(linkElements).map(link => link.getArray('href'));  
			return hrefValues
		}`, func(page *rod.Page) {
		page.MustWaitLoad()
	}).Stream().Map(func(i web_crawler.WebReaderString) any {
		path := web_crawler.ParsePath(u, i.String())
		article := web_crawler.GetArticleWithURL(path, nil)
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", path)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Content : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
		return article
	}).ToList()

	//}
}
func TestReadOpenAI(t *testing.T) {
	//for {
	u := "https://www.landiannews.com/feed"
	web_crawler.ReadRSS(u).Channel.Sites().Stream().Map(func(i rss.Item) any {
		path := web_crawler.ParsePath(u, i.Link)
		article := web_crawler.GetArticleWithChrome(path, nil, func(article web_crawler.Article) bool {
			return article.Title == ""
		})
		if article.IsErr() {
			sys.Panic(article.Error())
		}
		fmt.Printf("URL     : %s\n", path)
		fmt.Printf("Title   : %s\n", article.Title)
		fmt.Printf("Author  : %s\n", article.Byline)
		fmt.Printf("Length  : %d\n", article.Length)
		fmt.Printf("Excerpt : %s\n", article.Excerpt)
		fmt.Printf("Image  : %s\n", article.Image)
		fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Content : %s\n", article.Content)
		fmt.Printf("TextContent : %s\n", article.TextContent)
		fmt.Printf("SiteName: %s\n", article.SiteName)
		fmt.Printf("PubTime:%s\n", article.PublishedTime)
		return article
	}).ToList()

	//}
}

func TestSourceRead(t *testing.T) {
	code := web_crawler.ReadSourceCode("https://www.cnbeta.com.tw/articles/tech/1483576.htm", "", render.Render)
	println(web_crawler.StripHTMLTags(web_crawler.DiyConvertHtmlToArticle(code.String())))
}
func init() {
	web_crawler.ChromeRemoteURL = "127.0.0.1:9222"
	web_crawler.WebChrome = core.NewCotain("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", true)
}
```