package test

import (
	"fmt"
	"testing"

	"github.com/karosown/katool-go/words/jieba"
)

func Test_jieba(t *testing.T) {
	client := jieba.New()
	defer client.Free()
	fmt.Println(client.Cut("我测试一下测测测测 我Hello World"))
}

func Test_jieba_Frecuence(t *testing.T) {
	client := jieba.New()
	defer client.Free()
	fmt.Println(client.CutAll("下面是一个简洁的Go语言SDK，封装了 gojieba 库以简化中文分词的调用。这个SDK提供了一个清晰的 API，以便于开发者更容易地执行分词操作。").Frequency())
}
