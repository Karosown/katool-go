package mq

import (
	"errors"
	"time"
)

// 定义一些通用的错误类型，供实现层使用

var (
	// ErrBrokerClosed 当尝试向已关闭的 Broker 操作时返回
	ErrBrokerClosed = errors.New("mq: broker is closed")

	// ErrInvalidPartition 当请求的分区索引越界时返回
	ErrInvalidPartition = errors.New("mq: invalid partition index")
)

func WithKey(key string) PublishOption {
	return func(o *PublishOptions) { o.Key = key }
}

func WithDelay(d time.Duration) PublishOption {
	return func(o *PublishOptions) { o.Delay = d }
}

func WithExtra(key string, value any) PublishOption {
	return func(o *PublishOptions) {
		if o.Extra == nil {
			o.Extra = make(map[string]any)
		}
		o.Extra[key] = value
	}
}
func WithGroup(group string) SubscribeOption {
	return func(o *SubscribeOptions) { o.Group = group }
}

func WithConsumerID(id string) SubscribeOption {
	return func(o *SubscribeOptions) { o.ConsumerID = id }
}

func WithPartitions(partitions ...int) SubscribeOption {
	return func(o *SubscribeOptions) {
		if len(partitions) == 0 {
			o.SpecificPartitions = nil
			return
		}
		cp := make([]int, len(partitions))
		copy(cp, partitions)
		o.SpecificPartitions = cp
	}
}

func WithFilter(f FilterFunc) SubscribeOption {
	return func(o *SubscribeOptions) {
		if f == nil {
			return
		}
		prev := o.Filter
		if prev == nil {
			o.Filter = f
		} else {
			o.Filter = func(md Metadata) bool {
				return prev(md) && f(md)
			}
		}
	}
}
