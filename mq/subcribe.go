package mq

// FilterFunc 客户端过滤器函数
type FilterFunc func(md Metadata) bool

// SubscribeOptions 订阅配置结构体
type SubscribeOptions struct {
	Group              string     // 消费组
	ConsumerID         string     // 消费者ID
	SpecificPartitions []int      // 指定分区 (Assign模式)
	Filter             FilterFunc // 过滤器
}

// NewSubscribeOptions 辅助函数：初始化并应用选项
// 供 Broker 实现层调用，简化代码
func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	// 初始化一个零值对象
	options := SubscribeOptions{}

	// 依次应用选项
	for _, o := range opts {
		o(&options)
	}

	return options
}

// SubscribeOption 函数式选项类型
type SubscribeOption func(*SubscribeOptions)
