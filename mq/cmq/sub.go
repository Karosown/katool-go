package cmq

import (
	"context"
	"github.com/google/uuid"
	"github.com/karosown/katool-go/mq"
)

func (b *ChanBroker) Subscribe(ctx context.Context, topic string, handler mq.Handler, opts ...mq.SubscribeOption) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-b.ctx.Done():
		return mq.ErrBrokerClosed
	default:
	}

	options := mq.NewSubscribeOptions(opts...)
	if options.Group == "" {
		options.Group = "default_" + uuid.New().String()
	}

	b.mu.Lock()
	numParts := b.ensureTopic(topic)

	if len(options.SpecificPartitions) > 0 {
		for _, p := range options.SpecificPartitions {
			if p < 0 || p >= numParts {
				b.mu.Unlock()
				return mq.ErrInvalidPartition
			}
		}
	}

	// 初始化消费组物理管道
	if _, ok := b.queues[topic][options.Group]; !ok {
		chans := make([]chan *chanMessage, numParts)
		for i := 0; i < numParts; i++ {
			chans[i] = make(chan *chanMessage, 100)
		}
		b.queues[topic][options.Group] = &consumerGroup{
			name:       options.Group,
			partitions: chans,
		}
	}
	group := b.queues[topic][options.Group]
	b.mu.Unlock()

	// 消费逻辑
	consumeFunc := func(pIndex int) {
		ch := group.partitions[pIndex]
		for {
			select {
			case <-ctx.Done():
				return
			case <-b.ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				// 客户端过滤
				if options.Filter != nil {
					if !options.Filter(msg.GetMetadata()) {
						continue
					}
				}
				_ = handler(ctx, msg)
			}
		}
	}

	// 模式判断
	if len(options.SpecificPartitions) > 0 {
		// Assign 模式
		for _, p := range options.SpecificPartitions {
			if p >= 0 && p < numParts {
				go consumeFunc(p)
			}
		}
	} else {
		// 监听所有 (简化版 Rebalance)
		for i := 0; i < numParts; i++ {
			go consumeFunc(i)
		}
	}

	return nil
}
