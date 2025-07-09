package rss

import (
	"encoding/xml"

	"github.com/karosown/katool-go/container/stream"
)

// RSS RSS订阅源根结构体
// RSS represents the root structure of an RSS feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	DC      string   `xml:"xmlns:dc,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
	error
}

// Channel RSS频道信息结构体
// Channel represents RSS channel information structure
type Channel struct {
	Title         string   `xml:"title"`         // 频道标题 / Channel title
	Description   string   `xml:"description"`   // 频道描述 / Channel description
	Link          string   `xml:"link"`          // 频道链接 / Channel link
	Generator     string   `xml:"generator"`     // 生成器 / Generator
	LastBuildDate string   `xml:"lastBuildDate"` // 最后构建日期 / Last build date
	AtomLink      AtomLink `xml:"atom:link"`     // Atom链接 / Atom link
	Items         []Item   `xml:"item"`          // 文章项目列表 / Article items list
}

// Items RSS文章项目集合
// Items represents a collection of RSS article items
type Items []Item

// Sites 获取频道中的所有站点项目
// Sites gets all site items in the channel
func (c Channel) Sites() Items {
	return c.Items
}

// Stream 将项目集合转换为流
// Stream converts the items collection to a stream
func (c Items) Stream() *stream.Stream[Item, Items] {
	return stream.ToStream(&c)
}

// AtomLink Atom链接结构体
// AtomLink represents an Atom link structure
type AtomLink struct {
	Href string `xml:"href,attr"` // 链接地址 / Link URL
	Rel  string `xml:"rel,attr"`  // 关系 / Relationship
	Type string `xml:"type,attr"` // 类型 / Type
}

// Item RSS文章项目结构体
// Item represents an RSS article item structure
type Item struct {
	Title       string `xml:"title" json:"title,omitempty"`                 // 标题 / Title
	Description string `xml:"description" json:"description,omitempty"`     // 描述 / Description
	Link        string `xml:"link" json:"link,omitempty"`                   // 链接 / Link
	GUID        GUID   `xml:"guid" json:"guid,omitempty"`                   // 全局唯一标识符 / Global Unique Identifier
	Category    string `xml:"category,omitempty" json:"category,omitempty"` // 分类 / Category
	PubDate     string `xml:"pubDate" json:"pubDate,omitempty"`             // 发布日期 / Publication date
	Extra       string `xml:"extra" json:"extra,omitempty"`                 // 额外信息 / Extra information
}

// GUID 全局唯一标识符结构体
// GUID represents a Global Unique Identifier structure
type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr" json:"isPermaLink,omitempty"` // 是否为永久链接 / Whether it's a permanent link
	Value       string `xml:",chardata" json:"value,omitempty"`              // GUID值 / GUID value
}

// IsErr 检查是否有错误
// IsErr checks if there is an error
func (r *RSS) IsErr() bool {
	return r.error != nil
}

// SetErr 设置错误信息
// SetErr sets error information
func (r *RSS) SetErr(err error) *RSS {
	r.error = err
	return r
}
