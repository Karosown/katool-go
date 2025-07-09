package feed

import (
	"encoding/xml"
	"time"

	"github.com/karosown/katool-go/container/stream"
	"github.com/karosown/katool-go/web_crawler/rss"
)

// AtomFeed Atom格式的订阅源结构体
// AtomFeed represents an Atom format feed structure
type AtomFeed struct {
	XMLName xml.Name       `xml:"feed"`
	XMLNS   string         `xml:"xmlns,attr"`
	IDX     string         `xml:"xmlns:idx,attr,omitempty"`
	ID      string         `xml:"id"`
	Title   string         `xml:"title"`
	Link    []rss.AtomLink `xml:"link"` // 使用切片，因为可能有多个link / Use slice as there may be multiple links
	Updated string         `xml:"updated"`
	Entries []AtomEntry    `xml:"entry"`
	error
}

// AtomEntry Atom订阅源中的条目结构体
// AtomEntry represents an entry structure in Atom feed
type AtomEntry struct {
	ID        string         `xml:"id" json:"id,omitempty"`                         // 条目ID / Entry ID
	Title     AtomText       `xml:"title" json:"title,omitempty"`                   // 标题 / Title
	Link      rss.AtomLink   `xml:"link" json:"link,omitempty"`                     // 链接 / Link
	Published string         `xml:"published" json:"published,omitempty"`           // 发布时间 / Published time
	Updated   string         `xml:"updated" json:"updated,omitempty"`               // 更新时间 / Updated time
	Content   AtomText       `xml:"content" json:"content,omitempty"`               // 内容 / Content
	Author    AtomAuthor     `xml:"author" json:"author,omitempty"`                 // 作者 / Author
	Category  []AtomCategory `xml:"category,omitempty" json:"categories,omitempty"` // 分类 / Categories
}

// AtomEntries Atom条目集合
// AtomEntries represents a collection of Atom entries
type AtomEntries []AtomEntry

// AtomText 带有类型的文本内容结构体
// AtomText represents text content structure with type information
type AtomText struct {
	Type  string `xml:"type,attr,omitempty" json:"type,omitempty"` // 内容类型 / Content type
	Value string `xml:",chardata" json:"value,omitempty"`          // 文本值 / Text value
}

// AtomAuthor 条目作者信息结构体
// AtomAuthor represents entry author information structure
type AtomAuthor struct {
	Name  string `xml:"name" json:"name,omitempty"`             // 姓名 / Name
	Email string `xml:"email,omitempty" json:"email,omitempty"` // 邮箱 / Email
	URI   string `xml:"uri,omitempty" json:"uri,omitempty"`     // URI地址 / URI address
}

// AtomCategory 条目分类结构体
// AtomCategory represents entry category structure
type AtomCategory struct {
	Term   string `xml:"term,attr" json:"term,omitempty"`               // 分类术语 / Category term
	Scheme string `xml:"scheme,attr,omitempty" json:"scheme,omitempty"` // 分类方案 / Category scheme
	Label  string `xml:"label,attr,omitempty" json:"label,omitempty"`   // 分类标签 / Category label
}

// GetEntries 获取订阅源的所有条目
// GetEntries gets all entries in the feed
func (f AtomFeed) GetEntries() AtomEntries {
	return f.Entries
}

// Stream 将条目集合转换为流
// Stream converts the entries collection to a stream
func (e AtomEntries) Stream() *stream.Stream[AtomEntry, AtomEntries] {
	return stream.ToStream(&e)
}

// IsErr 检查是否有错误
// IsErr checks if there is an error
func (f *AtomFeed) IsErr() bool {
	return f.error != nil
}

// SetErr 设置错误信息
// SetErr sets error information
func (f *AtomFeed) SetErr(err error) *AtomFeed {
	f.error = err
	return f
}

// GetSelfLink 获取rel="self"的链接
// GetSelfLink gets the link with rel="self"
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
// GetMainLink gets the main link (no rel attribute or rel="alternate")
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

// GetTitleText 获取不带类型属性的标题文本
// GetTitleText gets title text without type attribute
func (e AtomEntry) GetTitleText() string {
	return e.Title.Value
}

// GetLinkHref 获取条目链接地址
// GetLinkHref gets the entry link URL
func (e AtomEntry) GetLinkHref() string {
	return e.Link.Href[len("https://www.google.com/url?rct=j&sa=t&url="):]
}

// GetContentText 获取内容文本
// GetContentText gets the content text
func (e AtomEntry) GetContentText() string {
	return e.Content.Value
}

// GetUpdatedTime 将订阅源updated字段解析为time.Time
// GetUpdatedTime parses the feed's updated field to time.Time
func (f AtomFeed) GetUpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, f.Updated)
}

// GetPublishedTime 将条目published字段解析为time.Time
// GetPublishedTime parses the entry's published field to time.Time
func (e AtomEntry) GetPublishedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Published)
}

// GetUpdatedTime 将条目updated字段解析为time.Time
// GetUpdatedTime parses the entry's updated field to time.Time
func (e AtomEntry) GetUpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Updated)
}

// ToRSSItem 将AtomEntry转换为RSS的Item
// ToRSSItem converts AtomEntry to RSS Item
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
// ToRSS converts AtomFeed to RSS
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

// ParseAtomFeed 解析Atom订阅源数据
// ParseAtomFeed parses Atom feed data
func ParseAtomFeed(data []byte) (*AtomFeed, error) {
	var feed AtomFeed
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, feed.SetErr(err)
	}
	return &feed, nil
}

// FeedParser 统一订阅源接口，便于同时处理RSS和Atom
// FeedParser is a unified feed interface for handling both RSS and Atom
type FeedParser interface {
	IsErr() bool             // 检查是否有错误 / Check if there is an error
	SetErr(error) FeedParser // 设置错误信息 / Set error information
}

// ParseFeed 根据内容自动判断并解析订阅源类型
// ParseFeed automatically detects and parses feed type based on content
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
