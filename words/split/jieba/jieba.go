package jieba

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/words/split"
	"github.com/wangbin/jiebago"
)

// Client 封装 gojieba 的结构体
type Client struct {
	jieba *jiebago.Segmenter
}

// New 创建一个新的 Client 实例
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

// Free 释放资源
func (j *Client) Free() {
	//j.jieba.
}

// Cut 精确模式分词
func (j *Client) Cut(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.Cut(text, true))
}

// CutAll 全模式分词
func (j *Client) CutAll(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.CutAll(text))
}

// CutForSearch 搜索引擎模式分词
func (j *Client) CutForSearch(text string) split.SplitStrings {
	return convert.AwaitChanToArray(j.jieba.CutForSearch(text, true))
}
