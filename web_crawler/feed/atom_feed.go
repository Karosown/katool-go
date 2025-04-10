package feed

import (
	"encoding/xml"
	"time"

	"github.com/karosown/katool/container/stream"
	"github.com/karosown/katool/web_crawler/rss"
)

// AtomFeed 表示Atom格式的feed
type AtomFeed struct {
	XMLName xml.Name       `xml:"feed"`
	XMLNS   string         `xml:"xmlns,attr"`
	IDX     string         `xml:"xmlns:idx,attr,omitempty"`
	ID      string         `xml:"id"`
	Title   string         `xml:"title"`
	Link    []rss.AtomLink `xml:"link"` // 使用切片，因为可能有多个link
	Updated string         `xml:"updated"`
	Entries []AtomEntry    `xml:"entry"`
	error
}

// AtomEntry 表示Atom feed中的条目
type AtomEntry struct {
	ID        string         `xml:"id" json:"id,omitempty"`
	Title     AtomText       `xml:"title" json:"title,omitempty"`
	Link      rss.AtomLink   `xml:"link" json:"link,omitempty"`
	Published string         `xml:"published" json:"published,omitempty"`
	Updated   string         `xml:"updated" json:"updated,omitempty"`
	Content   AtomText       `xml:"content" json:"content,omitempty"`
	Author    AtomAuthor     `xml:"author" json:"author,omitempty"`
	Category  []AtomCategory `xml:"category,omitempty" json:"categories,omitempty"`
}

// AtomEntries 表示多个AtomEntry的集合
type AtomEntries []AtomEntry

// AtomText 表示可能带有类型的文本内容
type AtomText struct {
	Type  string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Value string `xml:",chardata" json:"value,omitempty"`
}

// AtomAuthor 表示条目作者信息
type AtomAuthor struct {
	Name  string `xml:"name" json:"name,omitempty"`
	Email string `xml:"email,omitempty" json:"email,omitempty"`
	URI   string `xml:"uri,omitempty" json:"uri,omitempty"`
}

// AtomCategory 表示条目分类
type AtomCategory struct {
	Term   string `xml:"term,attr" json:"term,omitempty"`
	Scheme string `xml:"scheme,attr,omitempty" json:"scheme,omitempty"`
	Label  string `xml:"label,attr,omitempty" json:"label,omitempty"`
}

// 实现类似RSS的辅助方法

func (f AtomFeed) GetEntries() AtomEntries {
	return f.Entries
}

func (e AtomEntries) Stream() *stream.Stream[AtomEntry, AtomEntries] {
	return stream.ToStream(&e)
}

func (f *AtomFeed) IsErr() bool {
	return f.error != nil
}

func (f *AtomFeed) SetErr(err error) *AtomFeed {
	f.error = err
	return f
}

// GetSelfLink 获取rel="self"的链接
func (f AtomFeed) GetSelfLink() string {
	for _, link := range f.Link {
		if link.Rel == "self" {
			return link.Href
		}
	}
	// 如果没有找到rel="self"的链接，返回第一个链接
	if len(f.Link) > 0 {
		return f.Link[0].Href
	}
	return ""
}

// GetMainLink 获取主链接（无rel属性或rel="alternate"）
func (f AtomFeed) GetMainLink() string {
	for _, link := range f.Link {
		if link.Rel == "" || link.Rel == "alternate" {
			return link.Href
		}
	}
	// 如果没有找到主链接，返回第一个链接
	if len(f.Link) > 0 {
		return f.Link[0].Href
	}
	return ""
}

// 辅助方法，获取不带类型属性的标题文本
func (e AtomEntry) GetTitleText() string {
	return e.Title.Value
}

// 辅助方法，获取条目链接地址
func (e AtomEntry) GetLinkHref() string {
	return e.Link.Href[len("https://www.google.com/url?rct=j&sa=t&url="):]
}

// 辅助方法，获取内容文本
func (e AtomEntry) GetContentText() string {
	return e.Content.Value
}

// 时间解析相关方法

// GetUpdatedTime 将updated字段解析为time.Time
func (f AtomFeed) GetUpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, f.Updated)
}

// GetPublishedTime 将published字段解析为time.Time
func (e AtomEntry) GetPublishedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Published)
}

// GetUpdatedTime 将updated字段解析为time.Time
func (e AtomEntry) GetUpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Updated)
}

// 转换为RSS的方法

// ToRSSItem 将AtomEntry转换为RSS的Item
func (e AtomEntry) ToRSSItem() rss.Item {
	return rss.Item{
		Title:       e.GetTitleText(),
		Description: e.GetContentText(),
		Link:        e.GetLinkHref(),
		GUID: rss.GUID{
			IsPermaLink: true,
			Value:       e.ID,
		},
		PubDate: e.Published,
	}
}

// ToRSS 将AtomFeed转换为RSS
func (f AtomFeed) ToRSS() rss.RSS {
	items := make([]rss.Item, len(f.Entries))
	for i, entry := range f.Entries {
		items[i] = entry.ToRSSItem()
	}

	// 获取link
	link := f.GetMainLink()

	return rss.RSS{
		Version: "2.0",
		Channel: rss.Channel{
			Title:         f.Title,
			Description:   f.Title, // Atom没有直接对应的description字段，使用title代替
			Link:          link,
			LastBuildDate: f.Updated,
			Items:         items,
		},
	}
}

// 解析Atom Feed的函数
func ParseAtomFeed(data []byte) (*AtomFeed, error) {
	var feed AtomFeed
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, feed.SetErr(err)
	}
	return &feed, nil
}

// 统一Feed接口，便于同时处理RSS和Atom
type FeedParser interface {
	IsErr() bool
	SetErr(error) FeedParser
}

// ParseFeed 根据内容自动判断并解析Feed类型
func ParseFeed(data []byte) (interface{}, error) {
	// 尝试解析为RSS
	var rss rss.RSS
	rssErr := xml.Unmarshal(data, &rss)
	if rssErr == nil && rss.Channel.Title != "" {
		return &rss, nil
	}

	// 尝试解析为Atom
	var atom AtomFeed
	atomErr := xml.Unmarshal(data, &atom)
	if atomErr == nil && atom.Title != "" {
		return &atom, nil
	}

	// 两种格式都解析失败
	if rssErr != nil {
		return nil, rssErr
	}
	return nil, atomErr
}
