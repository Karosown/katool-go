package cgojieba

import (
	"github.com/karosown/katool-go/words/split"
	"github.com/yanyiwu/gojieba"
)

// Client CGO版本的中文分词客户端结构体（基于gojieba库）
// Client is a CGO-based Chinese word segmentation aiclient structure (based on gojieba library)
type Client struct {
	jieba *gojieba.Jieba // CGO分词器实例 / CGO segmenter instance
}

// New 创建一个新的CGO版本中文分词客户端实例
// New creates a new CGO-based Chinese word segmentation aiclient instance
func New(paths ...string) *Client {
	return &Client{
		jieba: gojieba.NewJieba(),
	}
}

// Free 释放CGO分词器资源（重要：必须调用以避免内存泄漏）
// Free releases CGO segmenter resources (important: must be called to avoid memory leaks)
func (j *Client) Free() {
	j.jieba.Free()
}

// Cut 精确模式分词，追求分词的准确度
// Cut performs accurate word segmentation for precision
func (j *Client) Cut(text string) split.SplitStrings {
	return j.jieba.Cut(text, true)
}

// CutAll 全模式分词，输出所有可能的词语
// CutAll performs full mode segmentation, outputs all possible words
func (j *Client) CutAll(text string) split.SplitStrings {
	return j.jieba.CutAll(text)
}

// CutForSearch 搜索引擎模式分词，适合用于搜索引擎的索引
// CutForSearch performs search engine mode segmentation, suitable for search engine indexing
func (j *Client) CutForSearch(text string) split.SplitStrings {
	return j.jieba.CutForSearch(text, true)
}
