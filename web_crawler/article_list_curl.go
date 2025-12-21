package web_crawler

import "github.com/go-rod/rod"

// ReadArticleLinksWithCURL 使用 curl 命令读取文章链接和标题。
// ReadArticleLinksWithCURL reads article links and titles using a curl command.
func ReadArticleLinksWithCURL(curl, js string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksWithCURL(curl, js, renderFunc, requestModifiers...)
}

// ReadArticleLinksWithCURL 使用 curl 命令读取文章链接和标题。
// ReadArticleLinksWithCURL reads article links and titles using a curl command.
func (c *Client) ReadArticleLinksWithCURL(curl, js string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	setup, url, err := parseCurlSetup(curl)
	if err != nil {
		return ArticleLinkValue{ArticleLinks: nil, WebReaderError: WebReaderError{err}}
	}
	return c.ReadArticleLinks(url, js, combineRenderFunc(setup, renderFunc), combineRequestModifiers(setup, requestModifiers)...)
}

// ReadArticleLinksAutoWithCURL 使用 curl 命令读取文章链接和标题（自动解析，无 JS 选择器）。
// ReadArticleLinksAutoWithCURL reads article links and titles using a curl command without custom JS selectors.
func ReadArticleLinksAutoWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAutoWithCURL(curl, renderFunc, requestModifiers...)
}

// ReadArticleLinksAutoWithCURL 使用 curl 命令读取文章链接和标题（自动解析，无 JS 选择器）。
// ReadArticleLinksAutoWithCURL reads article links and titles using a curl command without custom JS selectors.
func (c *Client) ReadArticleLinksAutoWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return c.ReadArticleLinksAutoWithCURLOptions(curl, renderFunc, DefaultLinkExtractOptions(), requestModifiers...)
}

// ReadArticleLinksAutoWithCURLOptions 使用 curl 命令读取文章链接和标题（带提取选项）。
// ReadArticleLinksAutoWithCURLOptions reads article links and titles using a curl command with extraction options.
func ReadArticleLinksAutoWithCURLOptions(curl string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.ReadArticleLinksAutoWithCURLOptions(curl, renderFunc, options, requestModifiers...)
}

// ReadArticleLinksAutoWithCURLOptions 使用 curl 命令读取文章链接和标题（带提取选项）。
// ReadArticleLinksAutoWithCURLOptions reads article links and titles using a curl command with extraction options.
func (c *Client) ReadArticleLinksAutoWithCURLOptions(curl string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) ArticleLinkValue {
	setup, url, err := parseCurlSetup(curl)
	if err != nil {
		return ArticleLinkValue{ArticleLinks: nil, WebReaderError: WebReaderError{err}}
	}
	return c.ReadArticleLinksAutoWithOptions(url, combineRenderFunc(setup, renderFunc), options, combineRequestModifiers(setup, requestModifiers)...)
}

// ReadAnchorArticlesWithCURL 使用 curl 命令读取锚点链接并获取标题/内容。
// ReadAnchorArticlesWithCURL reads anchor links and fetches title/content using a curl command.
func ReadAnchorArticlesWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) AnchorArticleValue {
	return DefaultClient.ReadAnchorArticlesWithCURL(curl, renderFunc, requestModifiers...)
}

// ReadAnchorArticlesWithCURL 使用 curl 命令读取锚点链接并获取标题/内容。
// ReadAnchorArticlesWithCURL reads anchor links and fetches title/content using a curl command.
func (c *Client) ReadAnchorArticlesWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) AnchorArticleValue {
	return c.ReadAnchorArticlesWithCURLOptions(curl, renderFunc, DefaultLinkExtractOptions(), requestModifiers...)
}

// ReadAnchorArticlesWithCURLOptions 使用 curl 命令读取锚点链接并获取标题/内容（带提取选项）。
// ReadAnchorArticlesWithCURLOptions reads anchor links and fetches title/content using a curl command with options.
func ReadAnchorArticlesWithCURLOptions(curl string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) AnchorArticleValue {
	return DefaultClient.ReadAnchorArticlesWithCURLOptions(curl, renderFunc, options, requestModifiers...)
}

// ReadAnchorArticlesWithCURLOptions 使用 curl 命令读取锚点链接并获取标题/内容（带提取选项）。
// ReadAnchorArticlesWithCURLOptions reads anchor links and fetches title/content using a curl command with options.
func (c *Client) ReadAnchorArticlesWithCURLOptions(curl string, renderFunc func(*rod.Page), options LinkExtractOptions, requestModifiers ...RequestWith) AnchorArticleValue {
	setup, url, err := parseCurlSetup(curl)
	if err != nil {
		return AnchorArticleValue{AnchorArticles: nil, WebReaderError: WebReaderError{err}}
	}
	return c.ReadAnchorArticlesWithOptions(url, combineRenderFunc(setup, renderFunc), options, combineRequestModifiers(setup, requestModifiers)...)
}

// AllLinkWithCURL 使用 curl 命令读取所有锚点链接（仅标题与链接）。
// AllLinkWithCURL reads all anchor links (titles + links only) using a curl command.
func AllLinkWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	return DefaultClient.AllLinkWithCURL(curl, renderFunc, requestModifiers...)
}

// AllLinkWithCURL 使用 curl 命令读取所有锚点链接（仅标题与链接）。
// AllLinkWithCURL reads all anchor links (titles + links only) using a curl command.
func (c *Client) AllLinkWithCURL(curl string, renderFunc func(*rod.Page), requestModifiers ...RequestWith) ArticleLinkValue {
	setup, url, err := parseCurlSetup(curl)
	if err != nil {
		return ArticleLinkValue{ArticleLinks: nil, WebReaderError: WebReaderError{err}}
	}
	return c.AllLink(url, combineRenderFunc(setup, renderFunc), combineRequestModifiers(setup, requestModifiers)...)
}
