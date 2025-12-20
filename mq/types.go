package mq

import (
	"time"
)

// Metadata 包含消息的标准元信息
// 使用强类型字段，extra 用于扩展
type Metadata struct {
	Topic     string    // 主题
	Key       string    // 路由键 (决定分区)
	Partition int       // 实际落入的分区 ID
	MessageID string    // 唯一消息 ID
	Timestamp time.Time // 消息生产时间

	// Extra 存储开发者自定义的元数据
	Extra map[string]any
}

// Message 消息抽象接口
type Message interface {
	// Payload 返回消息体字节
	Payload() []byte

	// GetMetadata 返回结构化的元数据对象
	GetMetadata() Metadata

	// Ack 确认消息 (用于可靠投递)
	Ack() error

	// Nack 拒绝消息 (requeue=true 表示重新入队)
	Nack(requeue bool) error
}
