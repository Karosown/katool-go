package mq

import (
	"context"
)

// Handler 消费者回调函数定义
type Handler func(ctx context.Context, msg Message) error

// Client 消息队列核心接口
// 定义了发布和订阅的标准行为
type Client interface {
	// Publish 发送消息
	Publish(ctx context.Context, topic string, msg []byte, opts ...PublishOption) error

	// Subscribe 订阅消息
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) error

	// Close 关闭客户端连接
	Close() error
}
