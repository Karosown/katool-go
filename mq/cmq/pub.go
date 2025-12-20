package cmq

import (
	"context"
	"github.com/google/uuid"
	"github.com/karosown/katool-go/mq"
	"time"
)

func (b *ChanBroker) Publish(ctx context.Context, topic string, msg []byte, opts ...mq.PublishOption) error {
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

	options := mq.NewPublishOptions(opts...)

	// 1. 处理延迟
	if options.Delay > 0 {
		time.AfterFunc(options.Delay, func() {
			_ = b.Publish(context.Background(), topic, msg, mq.WithKey(options.Key), func(o *mq.PublishOptions) {
				o.Extra = options.Extra
			})
		})
		return nil
	}

	b.mu.Lock()
	numParts := b.ensureTopic(topic)
	// 计算分区 (Util中定义)
	partIndex := b.hashKey(options.Key, numParts)
	// 获取该Topic下所有订阅组
	groupsMap := b.queues[topic]
	targetGroups := make([]*consumerGroup, 0, len(groupsMap))
	for _, group := range groupsMap {
		targetGroups = append(targetGroups, group)
	}
	b.mu.Unlock()

	// 2. 封装消息
	message := &chanMessage{
		payload: msg,
		meta: mq.Metadata{
			Topic:     topic,
			Key:       options.Key,
			Partition: partIndex,
			MessageID: uuid.New().String(),
			Timestamp: time.Now(),
			Extra:     options.Extra,
		},
	}

	// 3. 广播给所有 Group 的对应分区
	if len(targetGroups) == 0 {
		return nil
	}

	for _, group := range targetGroups {
		if partIndex < len(group.partitions) {
			select {
			case group.partitions[partIndex] <- message:
			case <-ctx.Done():
				return ctx.Err()
			case <-b.ctx.Done():
				return mq.ErrBrokerClosed
			default:
				// 内存满了可以选择丢弃或阻塞，这里简单丢弃
			}
		}
	}
	return nil
}
