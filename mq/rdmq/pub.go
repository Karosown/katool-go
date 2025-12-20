package rdmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/karosown/katool-go/mq"
	"github.com/redis/go-redis/v9"
)

func (c *RedisClient) Publish(ctx context.Context, topic string, msg []byte, opts ...mq.PublishOption) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-c.ctx.Done():
		return mq.ErrBrokerClosed
	default:
	}

	options := mq.NewPublishOptions(opts...)

	if options.Delay > 0 {
		time.AfterFunc(options.Delay, func() {
			_ = c.Publish(context.Background(), topic, msg, mq.WithKey(options.Key), func(o *mq.PublishOptions) {
				o.Extra = options.Extra
			})
		})
		return nil
	}

	partition := hashKey(options.Key, c.partitions)
	fields := map[string]any{
		"payload": msg,
		"key":     options.Key,
		"ts":      time.Now().UnixNano(),
	}
	if len(options.Extra) > 0 {
		if b, err := json.Marshal(options.Extra); err == nil {
			fields["extra"] = string(b)
		}
	}

	_, err := c.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: c.streamKey(topic, partition),
		Values: fields,
	}).Result()
	return err
}
