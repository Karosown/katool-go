package web_crawler

import (
	"errors"
	nurl "net/url"
	"regexp"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/go-rod/rod"
	"github.com/karosown/katool-go/container/stream"
	"golang.org/x/net/html"
)

// ArticleLink holds a single article title and link.
type ArticleLink struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"link,omitempty"`
}

// ArticleLinks is a collection of article links.
type ArticleLinks []ArticleLink

// Stream converts the collection to a stream.
func (a ArticleLinks) Stream() *stream.Stream[ArticleLink, ArticleLinks] {
	return stream.ToStream(&a)
}

// ArticleLinkValue wraps the article links with error information.
type ArticleLinkValue struct {
	ArticleLinks
	WebReaderError
}

// AnchorArticle holds article content fetched from an anchor link.
type AnchorArticle struct {
	Title       string `json:"title,omitempty"`
	Link        string `json:"link,omitempty"`
	Content     string `json:"content,omitempty"`
	TextContent string `json:"text_content,omitempty"`
}

// AnchorArticles is a collection of anchor articles.
type AnchorArticles []AnchorArticle

// Stream converts the collection to a stream.
func (a AnchorArticles) Stream() *stream.Stream[AnchorArticle, AnchorArticles] {
	return stream.ToStream(&a)
}

// AnchorArticleValue wraps the anchor articles with error information.
type AnchorArticleValue struct {
	AnchorArticles
	WebReaderError
}

// LinkExtractOptions controls how links are selected from a list page.
type LinkExtractOptions struct {
	Selector      string
	IncludeHref   []string
	ExcludeHref   []string
	IncludeText   []string
	ExcludeText   []string
	AllowExternal bool
	MaxLinks      int
}

// DefaultLinkExtractOptions returns the default extraction options.
func DefaultLinkExtractOptions() LinkExtractOptions {
	return LinkExtractOptions{
		AllowExternal: false,
		MaxLinks:      0,
	}
}

// AllAnchorOptions returns options that read all anchors, including external links.
func AllAnchorOptions() LinkExtractOptions {
	opts := DefaultLinkExtractOptions()
	opts.AllowExternal = true
	return opts
}

// ReadArticleLinks reads article links and titles from a list page.
func ReadArticleLinks(url, js string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinks(url, js, renderFunc, requestModifiers...)
}

// ReadArticleLinks reads article links and titles from a list page.
func (c *Client) ReadArticleLinks(url, js string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	if strings.TrimSpace(js) == "" {
		return c.ReadArticleLinksAuto(url, renderFunc, requestModifiers...)
	}
	return c.readArticleLinks(url, js, renderFunc, func(link string) Article {
		return c.GetArticleWithURL(link, requestModifiers...)
	})
}

// ReadArticleLinksWithChrome reads article links and titles using Chrome for article pages.
func ReadArticleLinksWithChrome(url, js string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksWithChrome(url, js, listRenderFunc, articleRenderFunc, restartCondition, i...)
}

// ReadArticleLinksWithChrome reads article links and titles using Chrome for article pages.
func (c *Client) ReadArticleLinksWithChrome(url, js string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) ArticleLinkValue {
	if strings.TrimSpace(js) == "" {
		return c.ReadArticleLinksAutoWithChrome(url, listRenderFunc, articleRenderFunc, restartCondition, i...)
	}
	return c.readArticleLinks(url, js, listRenderFunc, func(link string) Article {
		return c.GetArticleWithChrome(link, articleRenderFunc, restartCondition, i...)
	})
}

// ReadArticleLinksAuto reads article links and titles without custom JS selectors.
func ReadArticleLinksAuto(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAuto(url, renderFunc, requestModifiers...)
}

// ReadArticleLinksAuto reads article links and titles without custom JS selectors.
func (c *Client) ReadArticleLinksAuto(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return c.ReadArticleLinksAutoWithOptions(url, renderFunc, DefaultLinkExtractOptions(), requestModifiers...)
}

// ReadArticleLinksAutoWithOptions reads article links and titles with extraction options.
func ReadArticleLinksAutoWithOptions(url string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAutoWithOptions(url, renderFunc, options, requestModifiers...)
}

// ReadArticleLinksAutoWithOptions reads article links and titles with extraction options.
func (c *Client) ReadArticleLinksAutoWithOptions(url string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) ArticleLinkValue {
	return c.readArticleLinksAuto(url, renderFunc, options, func(link string) Article {
		return c.GetArticleWithURL(link, requestModifiers...)
	})
}

// ReadArticleLinksAutoWithChrome reads article links and titles without custom JS selectors.
func ReadArticleLinksAutoWithChrome(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAutoWithChrome(url, listRenderFunc, articleRenderFunc, restartCondition, i...)
}

// ReadArticleLinksAutoWithChrome reads article links and titles without custom JS selectors.
func (c *Client) ReadArticleLinksAutoWithChrome(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) ArticleLinkValue {
	return c.ReadArticleLinksAutoWithChromeOptions(url, listRenderFunc, articleRenderFunc, restartCondition, DefaultLinkExtractOptions(), i...)
}

// ReadArticleLinksAutoWithChromeOptions reads article links and titles with extraction options.
func ReadArticleLinksAutoWithChromeOptions(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, options LinkExtractOptions, i ...int) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAutoWithChromeOptions(url, listRenderFunc, articleRenderFunc, restartCondition, options, i...)
}

// ReadArticleLinksAutoWithChromeOptions reads article links and titles with extraction options.
func (c *Client) ReadArticleLinksAutoWithChromeOptions(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, options LinkExtractOptions, i ...int) ArticleLinkValue {
	return c.readArticleLinksAuto(url, listRenderFunc, options, func(link string) Article {
		return c.GetArticleWithChrome(link, articleRenderFunc, restartCondition, i...)
	})
}

// ReadAllAnchorArticles reads all anchor links and fetches title/content.
func ReadAllAnchorArticles(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) AnchorArticleValue {
	return DefaultClient.ReadAllAnchorArticles(url, renderFunc, requestModifiers...)
}

// ReadAllAnchorArticles reads all anchor links and fetches title/content.
func (c *Client) ReadAllAnchorArticles(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) AnchorArticleValue {
	return c.ReadAnchorArticlesWithOptions(url, renderFunc, AllAnchorOptions(), requestModifiers...)
}

// ReadAllAnchorArticlesWithChrome reads all anchor links and fetches title/content with Chrome.
func ReadAllAnchorArticlesWithChrome(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) AnchorArticleValue {
	return DefaultClient.ReadAllAnchorArticlesWithChrome(url, listRenderFunc, articleRenderFunc, restartCondition, i...)
}

// ReadAllAnchorArticlesWithChrome reads all anchor links and fetches title/content with Chrome.
func (c *Client) ReadAllAnchorArticlesWithChrome(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, i ...int) AnchorArticleValue {
	return c.ReadAnchorArticlesWithChromeOptions(url, listRenderFunc, articleRenderFunc, restartCondition, AllAnchorOptions(), i...)
}

// AllLink reads all anchor links on a page (titles + links only).
func AllLink(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.AllLink(url, renderFunc, requestModifiers...)
}

// AllLink reads all anchor links on a page (titles + links only).
func (c *Client) AllLink(url string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return c.allLinkAuto(url, renderFunc, AllAnchorOptions())
}

// ReadAnchorArticlesWithOptions reads anchor links and fetches title/content with options.
func ReadAnchorArticlesWithOptions(url string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) AnchorArticleValue {
	return DefaultClient.ReadAnchorArticlesWithOptions(url, renderFunc, options, requestModifiers...)
}

// ReadAnchorArticlesWithOptions reads anchor links and fetches title/content with options.
func (c *Client) ReadAnchorArticlesWithOptions(url string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) AnchorArticleValue {
	return c.readAnchorArticlesAuto(url, renderFunc, options, func(link string) Article {
		return c.GetArticleWithURL(link, requestModifiers...)
	})
}

// ReadAnchorArticlesWithChromeOptions reads anchor links and fetches title/content with Chrome and options.
func ReadAnchorArticlesWithChromeOptions(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, options LinkExtractOptions, i ...int) AnchorArticleValue {
	return DefaultClient.ReadAnchorArticlesWithChromeOptions(url, listRenderFunc, articleRenderFunc, restartCondition, options, i...)
}

// ReadAnchorArticlesWithChromeOptions reads anchor links and fetches title/content with Chrome and options.
func (c *Client) ReadAnchorArticlesWithChromeOptions(url string, listRenderFunc func(*rod.Page), articleRenderFunc func(*rod.Page), restartCondition func(Article) bool, options LinkExtractOptions, i ...int) AnchorArticleValue {
	return c.readAnchorArticlesAuto(url, listRenderFunc, options, func(link string) Article {
		return c.GetArticleWithChrome(link, articleRenderFunc, restartCondition, i...)
	})
}

func (c *Client) readArticleLinks(url, js string, renderFunc func(*rod.Page), fetch func(string) Article) ArticleLinkValue {
	list := c.ReadArray(url, js, renderFunc)
	if list.IsErr() {
		return ArticleLinkValue{
			ArticleLinks:   nil,
			WebReaderError: list.WebReaderError,
		}
	}

	res := make(ArticleLinks, 0, len(list.WebReaderValue))
	var errs error
	for _, item := range list.WebReaderValue {
		path := ParsePath(url, item.String())
		if path == "" {
			continue
		}
		article := fetch(path)
		if article.IsErr() {
			errs = errors.Join(errs, article.WebReaderError)
			res = append(res, ArticleLink{Title: "", Link: path})
			continue
		}
		res = append(res, ArticleLink{Title: article.Title, Link: path})
	}

	return ArticleLinkValue{
		ArticleLinks:   res,
		WebReaderError: WebReaderError{errs},
	}
}

func (c *Client) readArticleLinksAuto(url string, renderFunc func(*rod.Page), options LinkExtractOptions, fetch func(string) Article) ArticleLinkValue {
	code := c.ReadSourceCode(url, "", renderFunc)
	if code.IsErr() {
		return ArticleLinkValue{
			ArticleLinks:   nil,
			WebReaderError: code.WebReaderError,
		}
	}

	links, err := extractArticleLinksWithOptions(code.String(), url, options)
	if err != nil {
		return ArticleLinkValue{
			ArticleLinks:   nil,
			WebReaderError: WebReaderError{err},
		}
	}
	var errs error
	for i := range links {
		if links[i].Title != "" {
			continue
		}
		article := fetch(links[i].Link)
		if article.IsErr() {
			errs = errors.Join(errs, article.WebReaderError)
			continue
		}
		links[i].Title = article.Title
	}

	return ArticleLinkValue{
		ArticleLinks:   links,
		WebReaderError: WebReaderError{errs},
	}
}

func (c *Client) readAnchorArticlesAuto(url string, renderFunc func(*rod.Page), options LinkExtractOptions, fetch func(string) Article) AnchorArticleValue {
	code := c.ReadSourceCode(url, "", renderFunc)
	if code.IsErr() {
		return AnchorArticleValue{
			AnchorArticles: nil,
			WebReaderError: code.WebReaderError,
		}
	}

	links, err := extractArticleLinksWithOptions(code.String(), url, options)
	if err != nil {
		return AnchorArticleValue{
			AnchorArticles: nil,
			WebReaderError: WebReaderError{err},
		}
	}

	res := make(AnchorArticles, 0, len(links))
	var errs error
	for _, link := range links {
		article := fetch(link.Link)
		if article.IsErr() {
			errs = errors.Join(errs, article.WebReaderError)
			res = append(res, AnchorArticle{
				Title: link.Title,
				Link:  link.Link,
			})
			continue
		}
		title := article.Title
		if title == "" {
			title = link.Title
		}
		res = append(res, AnchorArticle{
			Title:       title,
			Link:        link.Link,
			Content:     article.Content,
			TextContent: article.TextContent,
		})
	}

	return AnchorArticleValue{
		AnchorArticles: res,
		WebReaderError: WebReaderError{errs},
	}
}

func (c *Client) allLinkAuto(url string, renderFunc func(*rod.Page), options LinkExtractOptions) ArticleLinkValue {
	code := c.ReadSourceCode(url, "", renderFunc)
	if code.IsErr() {
		return ArticleLinkValue{
			ArticleLinks:   nil,
			WebReaderError: code.WebReaderError,
		}
	}

	links, err := extractArticleLinksWithOptions(code.String(), url, options)
	if err != nil {
		return ArticleLinkValue{
			ArticleLinks:   nil,
			WebReaderError: WebReaderError{err},
		}
	}
	return ArticleLinkValue{
		ArticleLinks:   links,
		WebReaderError: WebReaderError{},
	}
}

func extractArticleLinksWithOptions(htmlStr, baseURL string, options LinkExtractOptions) (ArticleLinks, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}
	base, _ := nurl.Parse(baseURL)
	result := make(map[string]ArticleLink)

	includeHref, err := compilePatterns(options.IncludeHref)
	if err != nil {
		return nil, err
	}
	excludeHref, err := compilePatterns(options.ExcludeHref)
	if err != nil {
		return nil, err
	}
	includeText, err := compilePatterns(options.IncludeText)
	if err != nil {
		return nil, err
	}
	excludeText, err := compilePatterns(options.ExcludeText)
	if err != nil {
		return nil, err
	}

	maxLinks := options.MaxLinks
	if maxLinks < 0 {
		maxLinks = 0
	}

	var nodes []*html.Node
	if strings.TrimSpace(options.Selector) != "" {
		sel, selErr := cascadia.Compile(options.Selector)
		if selErr != nil {
			return nil, selErr
		}
		nodes = sel.MatchAll(doc)
	} else {
		nodes = []*html.Node{doc}
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if maxLinks > 0 && len(result) >= maxLinks {
			return
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			href := strings.TrimSpace(getAttr(n, "href"))
			if shouldSkipHref(href) {
				return
			}
			abs := ParsePath(baseURL, href)
			if abs == "" {
				return
			}
			if !options.AllowExternal && !sameHost(base, abs) {
				return
			}

			if !matchInclude(includeHref, abs) || matchExclude(excludeHref, abs) {
				return
			}

			title := strings.TrimSpace(collapseSpaces(nodeText(n)))
			if !matchInclude(includeText, title) || matchExclude(excludeText, title) {
				return
			}
			existing, ok := result[abs]
			if !ok || (existing.Title == "" && title != "") {
				result[abs] = ArticleLink{Title: title, Link: abs}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, n := range nodes {
		walk(n)
		if maxLinks > 0 && len(result) >= maxLinks {
			break
		}
	}

	res := make(ArticleLinks, 0, len(result))
	for _, v := range result {
		res = append(res, v)
	}
	return res, nil
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func nodeText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func shouldSkipHref(href string) bool {
	if href == "" {
		return true
	}
	lower := strings.ToLower(href)
	return strings.HasPrefix(lower, "#") ||
		strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "tel:")
}

func sameHost(base *nurl.URL, link string) bool {
	if base == nil || base.Host == "" {
		return true
	}
	u, err := nurl.Parse(link)
	if err != nil || u.Host == "" {
		return true
	}
	return normalizeHost(base.Host) == normalizeHost(u.Host)
}

func normalizeHost(host string) string {
	host = strings.ToLower(host)
	return strings.TrimPrefix(host, "www.")
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		return nil, nil
	}
	res := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if strings.TrimSpace(p) == "" {
			continue
		}
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func matchInclude(patterns []*regexp.Regexp, value string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if p.MatchString(value) {
			return true
		}
	}
	return false
}

func matchExclude(patterns []*regexp.Regexp, value string) bool {
	if len(patterns) == 0 {
		return false
	}
	for _, p := range patterns {
		if p.MatchString(value) {
			return true
		}
	}
	return false
}
