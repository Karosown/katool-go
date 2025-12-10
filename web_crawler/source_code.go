package web_crawler

import (
	"bytes"
	"encoding/xml"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/web_crawler/feed"
	"github.com/karosown/katool-go/web_crawler/rss"
	"golang.org/x/net/html"
)

// SourceCode 源代码结构体，包含网页源码和错误信息
// SourceCode represents source code structure containing web page source and error information
type SourceCode struct {
	WebReaderString
	WebReaderError
}

// ChromeRemoteURL Chrome远程调试URL
// ChromeRemoteURL is the Chrome remote debugging URL
var ChromeRemoteURL string

// ReadRSS 读取RSS订阅源
// ReadRSS reads RSS feed source
func ReadRSS(xmlURL string) rss.RSS {
	return DefaultClient.ReadRSS(xmlURL)
}

// ReadRSS 读取RSS订阅源
// ReadRSS reads RSS feed source
func (c *Client) ReadRSS(xmlURL string) rss.RSS {
	code := c.ReadSourceCode(xmlURL,
		rss.SourceCodeGetFunc, func(page *rod.Page) {
			page.MustWaitLoad()
		})
	r := &rss.RSS{}
	// 如果是html，那么取出pre
	err := xml.Unmarshal([]byte(code.String()), r)
	return *r.SetErr(err)
}

// ReadFeed 读取Atom订阅源并转换为RSS格式
// ReadFeed reads Atom feed source and converts to RSS format
func ReadFeed(xmlURL string) rss.RSS {
	return DefaultClient.ReadFeed(xmlURL)
}

// ReadFeed 读取Atom订阅源并转换为RSS格式
// ReadFeed reads Atom feed source and converts to RSS format
func (c *Client) ReadFeed(xmlURL string) rss.RSS {
	code := c.ReadSourceCode(xmlURL, rss.SourceCodeGetFunc, func(page *rod.Page) {
		page.MustWaitLoad()
	})
	r, err := feed.ParseAtomFeed([]byte(code.String()))
	if err != nil {
		return *((&rss.RSS{}).SetErr(err))
	}
	return r.ToRSS()
}

// ReadSourceCode 读取网页源代码（带重试机制）
// ReadSourceCode reads web page source code (with retry mechanism)
func ReadSourceCode(url, execJsFunc string, rendorFunc func(*rod.Page)) SourceCode {
	return DefaultClient.ReadSourceCode(url, execJsFunc, rendorFunc)
}

// ReadSourceCode 读取网页源代码（带重试机制）
// ReadSourceCode reads web page source code (with retry mechanism)
func (c *Client) ReadSourceCode(url, execJsFunc string, rendorFunc func(*rod.Page)) SourceCode {
	var gen func() SourceCode
	tryNum := 7
	gen = func() SourceCode {
		code := c.readSourceCode(url, execJsFunc, rendorFunc)
		if code.IsErr() {
			if tryNum != 0 {
				tryNum--
				if tryNum == 0 {
					if c.getChrome() != nil {
						c.getChrome().ReStart()
					}
				} else {
					time.Sleep(time.Duration(7-tryNum+1) * time.Second)
				}
				return gen()
			}
			return SourceCode{
				WebReaderString: "",
				WebReaderError: WebReaderError{
					errors.New("the Code Get Error"),
				},
			}
		}
		return code
	}
	return gen()
}

// readSourceCode 内部读取源代码方法
// readSourceCode is the internal method for reading source code
func (c *Client) readSourceCode(url, execJsFunc string, rendorFunc func(*rod.Page)) SourceCode {

	sourceCode, err := c.execFun(url, optional.IsTrue(execJsFunc != "", execJsFunc, "() => {"+
		"	return document.documentElement.outerHTML"+
		"}"), rendorFunc)
	if err != nil {
		return SourceCode{
			"", WebReaderError{err},
		}
	}
	return SourceCode{
		WebReaderString(optional.IsTrueByFunc(sourceCode != nil, func() string {
			return sourceCode.Value.String()
		}, optional.EmptyStringFunc)), WebReaderError{err},
	}
}

// DiyConvertHtmlToArticle 处理HTML字符串以提取和清理文章内容
// DiyConvertHtmlToArticle processes an HTML string to extract and clean the article content
func DiyConvertHtmlToArticle(code string) string {
	// Parse the HTML string into a node tree
	doc, err := html.Parse(strings.NewReader(code))
	if err != nil {
		return code // Return original code if parsing fails
	}

	// Locate the <body> tag
	var body *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)
	if body == nil {
		return code // Return original code if no <body> tag is found
	}

	// Remove unwanted elements: <script>, <devsite-header>, and tags with class "nocontent"
	var removeUnwanted func(*html.Node)
	removeUnwanted = func(n *html.Node) {
		var next *html.Node
		for c := n.FirstChild; c != nil; c = next {
			next = c.NextSibling // Store next sibling since we might remove c
			if c.Type == html.ElementNode {
				// Remove <script> and <devsite-header> tags
				if c.Data == "script" || c.Data == "devsite-header" {
					n.RemoveChild(c)
					continue
				}
				// Remove tags with class "nocontent"
				for _, attr := range c.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "nocontent") {
						n.RemoveChild(c)
						break
					}
				}
			}
			// Recursively process child nodes
			removeUnwanted(c)
		}
	}
	removeUnwanted(body)

	// Render the cleaned <body> content to a string
	var buf bytes.Buffer
	if err := html.Render(&buf, body); err != nil {
		return code // Return original code if rendering fails
	}

	// Extract the inner content of the <body> tag
	rendered := buf.String()
	start := strings.Index(rendered, ">") + 1
	end := strings.LastIndex(rendered, "<")
	if start >= end || start < 1 {
		return code // Return original code if extraction fails
	}
	cleanedContent := rendered[start:end]

	// Remove excess whitespace
	cleanedContent = strings.TrimSpace(cleanedContent)

	return cleanedContent
}

// StripHTMLTags 使用正则表达式移除HTML标签
// StripHTMLTags removes HTML tags using regular expressions
func StripHTMLTags(html string) string {
	// 移除所有HTML标签
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	text := tagRegex.ReplaceAllString(html, "")

	//// 处理HTML实体
	//text = html.UnescapeString(text)

	// 清理多余的空白字符
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
