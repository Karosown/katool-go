package rss

import (
	"encoding/xml"

	"github.com/karosown/katool-go/container/stream"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	DC      string   `xml:"xmlns:dc,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
	error
}

type Channel struct {
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	Generator     string   `xml:"generator"`
	LastBuildDate string   `xml:"lastBuildDate"`
	AtomLink      AtomLink `xml:"atom:link"`
	Items         []Item   `xml:"item"`
}
type Items []Item

func (c Channel) Sites() Items {
	return c.Items
}
func (c Items) Stream() *stream.Stream[Item, Items] {
	return stream.ToStream(&c)
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       string `xml:"title" json:"title,omitempty"`
	Description string `xml:"description" json:"description,omitempty"`
	Link        string `xml:"link" json:"link,omitempty"`
	GUID        GUID   `xml:"guid" json:"guid,omitempty"`
	Category    string `xml:"category,omitempty" json:"category,omitempty"`
	PubDate     string `xml:"pubDate" json:"pubDate,omitempty"`
	Extra       string `xml:"extra" json:"extra,omitempty"`
}

type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr" json:"isPermaLink,omitempty"`
	Value       string `xml:",chardata" json:"value,omitempty"`
}

func (r *RSS) IsErr() bool {
	return r.error != nil
}

func (r *RSS) SetErr(err error) *RSS {
	r.error = err
	return r
}
