package template

import (
	"fmt"
	"github.com/karosown/katool-go/container/xmap"
	"strings"
)

type Engine[T SenderAdapter] struct {
	Text   string
	Mapper xmap.Map[string, string]
	Prefix string
	Suffix string
}

// NewEngine 创建新的模板引擎
func NewEngine[T SenderAdapter](text string) *Engine[T] {
	return &Engine[T]{
		Text:   text,
		Mapper: xmap.NewMap[string, string](),
		Prefix: "#{",
		Suffix: "}",
	}
}

// SetDelimiters 设置分隔符
func (e *Engine[T]) SetDelimiters(prefix, suffix string) *Engine[T] {
	e.Prefix = prefix
	e.Suffix = suffix
	return e
}

// AddMapping 添加映射
func (e *Engine[T]) AddMapping(key, value string) *Engine[T] {
	e.Mapper.Set(key, value)
	return e
}

// AddMappings 批量添加映射
func (e *Engine[T]) AddMappings(mappings map[string]string) *Engine[T] {
	for k, v := range mappings {
		e.Mapper.Set(k, v)
	}
	return e
}

// Load 加载并替换模板
func (e *Engine[T]) Load() string {
	message := e.Text
	e.Mapper.ForEach(func(k, v string) {
		message = strings.ReplaceAll(message,
			fmt.Sprintf("%s%s%s", e.Prefix, k, e.Suffix),
			v)
	})
	return message
}

// Validate 验证模板是否有效
func (e *Engine[T]) Validate() (*string, error) {
	// 检查是否有未替换的模板变量
	message := e.Load()
	if strings.Contains(message, e.Prefix) && strings.Contains(message, e.Suffix) {
		return nil, fmt.Errorf("template contains unresolved variables")
	}
	return &message, nil
}

// Send 发送消息
func (e *Engine[T]) Send(sender Sender[T]) error {
	validate, err := e.Validate()
	if err != nil {
		return err
	}
	return sender.Send(*validate)
}
