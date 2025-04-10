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

type SourceCode struct {
	WebReaderString
	WebReaderError
}

var ChromeRemoteURL string

func ReadRSS(xmlURL string) rss.RSS {
	code := ReadSourceCode(xmlURL,
		rss.SourceCodeGetFunc, func(page *rod.Page) {
			page.MustWaitLoad()
		})
	r := &rss.RSS{}
	// 如果是html，那么取出pre
	err := xml.Unmarshal([]byte(code.String()), r)
	return *r.SetErr(err)
}
func ReadFeed(xmlURL string) rss.RSS {
	code := ReadSourceCode(xmlURL, rss.SourceCodeGetFunc, func(page *rod.Page) {
		page.MustWaitLoad()
	})
	r, err := feed.ParseAtomFeed([]byte(code.String()))
	if err != nil {
		return *((&rss.RSS{}).SetErr(err))
	}
	return r.ToRSS()
}
func ReadSourceCode(url, execJsFunc string, rendorFunc func(*rod.Page)) SourceCode {
	var gen func() SourceCode
	tryNum := 7
	gen = func() SourceCode {
		code := readSourceCode(url, execJsFunc, rendorFunc)
		if code.IsErr() {
			if tryNum != 0 {
				tryNum--
				if tryNum == 0 {
					WebChrome.ReStart()
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

func readSourceCode(url, execJsFunc string, rendorFunc func(*rod.Page)) SourceCode {

	sourceCode, err := execFun(url, optional.IsTrue(execJsFunc != "", execJsFunc, "() => {"+
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

// DiyConvertHtmlToArticle processes an HTML string to extract and clean the article content.
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

// 方法一：使用正则表达式
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
