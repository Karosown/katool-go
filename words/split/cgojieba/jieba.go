package cgojieba

import (
	"github.com/karosown/katool-go/words/split"
	"github.com/yanyiwu/gojieba"
)

// Client 封装 gojieba 的结构体
type Client struct {
	jieba *gojieba.Jieba
}

// New 创建一个新的 Client 实例
func New(paths ...string) *Client {
	return &Client{
		jieba: gojieba.NewJieba(),
	}
}

// Free 释放资源
func (j *Client) Free() {
	j.jieba.Free()
}

// Cut 精确模式分词
func (j *Client) Cut(text string) split.SplitStrings {
	return j.jieba.Cut(text, true)
}

// CutAll 全模式分词
func (j *Client) CutAll(text string) split.SplitStrings {
	return j.jieba.CutAll(text)
}

// CutForSearch 搜索引擎模式分词
func (j *Client) CutForSearch(text string) split.SplitStrings {
	return j.jieba.CutForSearch(text, true)
}
