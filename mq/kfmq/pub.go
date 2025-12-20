package kfmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/karosown/katool-go/mq"
	"github.com/segmentio/kafka-go"
)

func (c *KafkaClient) Publish(ctx context.Context, topic string, msg []byte, opts ...mq.PublishOption) error {
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

	var headers []kafka.Header
	if len(options.Extra) > 0 {
		if b, err := json.Marshal(options.Extra); err == nil {
			headers = append(headers, kafka.Header{Key: "extra", Value: b})
		}
	}

	kmsg := kafka.Message{
		Topic:   topic,
		Key:     []byte(options.Key),
		Value:   msg,
		Time:    time.Now(),
		Headers: headers,
	}
	return c.writer.WriteMessages(ctx, kmsg)
}
