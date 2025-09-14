package jieba

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/words/split"
	"github.com/wangbin/jiebago"
)

// Client 中文分词客户端结构体（基于jiebago库）
// Client is a Chinese word segmentation aiclient structure (based on jiebago library)
type Client struct {
	jieba *jiebago.Segmenter // 分词器实例 / Segmenter instance
}

// New 创建一个新的中文分词客户端实例
// New creates a new Chinese word segmentation aiclient instance
func New(path ...string) *Client {
	jieba := &jiebago.Segmenter{}
	if cutil.IsEmpty(path) || cutil.IsBlank(path[0]) {
		dictPath := fileutil.CurrentPath() + "/static/dict.txt"
		path = []string{dictPath}
	}
	jieba.LoadDictionary(path[0])
	return &Client{
		jieba: jieba,
	}
}

// Free 释放分词器资源
// Free releases segmenter resources
func (j *Client) Free() {
	//j.jieba.
}

// Cut 精确模式分词，追求分词的准确度
// Cut performs accurate word segmentation for precision
func (j *Client) Cut(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.Cut(text, true))
}

// CutAll 全模式分词，输出所有可能的词语
// CutAll performs full mode segmentation, outputs all possible words
func (j *Client) CutAll(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.CutAll(text))
}

// CutForSearch 搜索引擎模式分词，适合用于搜索引擎的索引
// CutForSearch performs search engine mode segmentation, suitable for search engine indexing
func (j *Client) CutForSearch(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.CutForSearch(text, true))
}
