package web_crawler

import (
	"errors"
	"fmt"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-shiori/go-readability"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/web_crawler/core"
)

// WebChrome Chrome浏览器实例
// WebChrome is a Chrome browser instance
var WebChrome *core.Contain

// WebReaderError Web读取错误包装器
// WebReaderError is a web reading error wrapper
type WebReaderError struct {
	error
}

// WebReaderString Web读取字符串类型
// WebReaderString is a web reading string type
type WebReaderString string

// WebReaderValue Web读取值集合
// WebReaderValue is a collection of web reading values
type WebReaderValue []WebReaderString

// NewWebReaderValue 创建新的Web读取值
// NewWebReaderValue creates a new web reading value
func NewWebReaderValue(strs ...string) WebReaderValue {
	res := []WebReaderString{}
	for _, str := range strs {
		res = append(res, WebReaderString(str))
	}
	return res
}

// String 转换为字符串
// String converts to string
func (w WebReaderString) String() string {
	return string(w)
}

// IsErr 检查是否有错误
// IsErr checks if there is an error
func (w WebReaderError) IsErr() bool {
	return w.error != nil
}

// SolveErrors 处理错误
// SolveErrors handles errors
func (w WebReaderError) SolveErrors(errs error) error {
	err := w
	if w.error == nil {
		return nil
	}
	if errs == nil {
		errs = err
	}
	errors.Join(errs, err)
	return errs
}

// Stream 转换为流
// Stream converts to stream
func (w WebReaderValue) Stream() *stream.Stream[WebReaderString, WebReaderValue] {
	return stream.ToStream(&w)
}

// Article 文章结构体
// Article represents an article structure
type Article struct {
	readability.Article
	WebReaderError
}

// RequestWith 请求修改函数类型
// RequestWith is a function type for request modification
type RequestWith func(r *http.Request)

// fromURL fetch the web page from specified url then parses the response to find
// the readable content.
// fromURL 从指定URL获取网页并解析响应以找到可读内容
// fromURL fetches the web page from specified URL then parses the response to find the readable content
func fromURL(pageURL string, timeout time.Duration, requestModifiers ...RequestWith) (readability.Article, error) {
	// Make sure URL is valid
	parsedURL, err := nurl.ParseRequestURI(pageURL)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Fetch page from URL
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", pageURL, nil)
	if requestModifiers != nil && len(requestModifiers) > 0 {
		for _, modifer := range requestModifiers {
			if modifer != nil {
				modifer(req)
			}
		}
	}
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to fetch the page: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return readability.Article{}, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	// Make sure content type is HTML
	cp := resp.Header.Get("Content-Type")
	if !strings.Contains(cp, "text/html") {
		return readability.Article{}, fmt.Errorf("URL is not a HTML document")
	}

	// Parse content
	parser := readability.NewParser()
	return parser.Parse(resp.Body, parsedURL)
}
func execFun(url, js string, rendorFunc func(*rod.Page)) (*proto.RuntimeRemoteObject, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()
	}()
	mustPage := WebChrome.MustPage(url)
	defer mustPage.Close()
	// Wait for the page loading...
	if rendorFunc != nil {
		rendorFunc(mustPage)
	}
	eval, err := mustPage.
		Eval(js)
	return eval, err
}

func ParsePath(url, path string) string {
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "https") {
		return path
	}
	if strings.HasPrefix(path, "/") || !strings.HasPrefix(path, "./") {
		// 获取url的host地址
		u, err := nurl.Parse(url)
		if err != nil {
			return ""
		}
		host := u.Host
		return url[:strings.Index(url, "://")] + "://" + optional.IsTrue(strings.HasSuffix(host, "/"), host[:len(host)-1]+path, host+path)
	}
	return url + "/" + path
}
