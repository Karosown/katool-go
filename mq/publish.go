package mq

import (
	"time"
)

// PublishOptions 发布配置结构体
type PublishOptions struct {
	Key   string         // 路由键
	Delay time.Duration  // 延迟时间
	Extra map[string]any // 自定义元数据
}

// NewPublishOptions 辅助函数：初始化并应用选项
func NewPublishOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{
		Extra: make(map[string]any), // 预初始化 map，防止 nil panic
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// PublishOption 函数式选项类型
type PublishOption func(*PublishOptions)
