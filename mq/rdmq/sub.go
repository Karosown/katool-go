package rdmq

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/karosown/katool-go/mq"
	"github.com/redis/go-redis/v9"
)

type streamMessage struct {
	payload []byte
	meta    mq.Metadata

	rdb      *redis.Client
	stream   string
	group    string
	id       string
	extraRaw string
	key      string

	doneOnce sync.Once
	doneErr  error
}

func (m *streamMessage) Payload() []byte          { return m.payload }
func (m *streamMessage) GetMetadata() mq.Metadata { return m.meta }

func (m *streamMessage) Ack() error {
	return m.finalize(func() error {
		return m.rdb.XAck(context.Background(), m.stream, m.group, m.id).Err()
	})
}

func (m *streamMessage) Nack(requeue bool) error {
	return m.finalize(func() error {
		if err := m.rdb.XAck(context.Background(), m.stream, m.group, m.id).Err(); err != nil {
			return err
		}
		if !requeue {
			return nil
		}
		values := map[string]any{
			"payload": m.payload,
			"key":     m.key,
			"ts":      time.Now().UnixNano(),
		}
		if m.extraRaw != "" {
			values["extra"] = m.extraRaw
		}
		return m.rdb.XAdd(context.Background(), &redis.XAddArgs{
			Stream: m.stream,
			Values: values,
		}).Err()
	})
}

func (m *streamMessage) finalize(fn func() error) error {
	m.doneOnce.Do(func() {
		m.doneErr = fn()
	})
	return m.doneErr
}

func (c *RedisClient) Subscribe(ctx context.Context, topic string, handler mq.Handler, opts ...mq.SubscribeOption) error {
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

	options := mq.NewSubscribeOptions(opts...)
	if options.Group == "" {
		options.Group = "default_" + uuid.New().String()
	}
	if options.ConsumerID == "" {
		options.ConsumerID = uuid.New().String()
	}

	partitions := c.partitions
	targetPartitions := options.SpecificPartitions
	if len(targetPartitions) == 0 {
		targetPartitions = make([]int, partitions)
		for i := 0; i < partitions; i++ {
			targetPartitions[i] = i
		}
	} else {
		for _, p := range targetPartitions {
			if p < 0 || p >= partitions {
				return mq.ErrInvalidPartition
			}
		}
	}

	for _, p := range targetPartitions {
		stream := c.streamKey(topic, p)
		if err := createGroupIfNeeded(ctx, c.rdb, stream, options.Group); err != nil {
			return err
		}
		go c.consumeStream(ctx, stream, options.Group, options.ConsumerID, handler, options.Filter)
	}

	return nil
}

func (c *RedisClient) consumeStream(
	ctx context.Context,
	stream string,
	group string,
	consumer string,
	handler mq.Handler,
	filter mq.FilterFunc,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.ctx.Done():
			return
		default:
		}

		resp, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) {
				continue
			}
			continue
		}

		for _, s := range resp {
			for _, msg := range s.Messages {
				sm := c.buildStreamMessage(stream, group, msg)
				if filter != nil && !filter(sm.meta) {
					_ = sm.Ack()
					continue
				}
				if err := handler(ctx, sm); err != nil {
					_ = sm.Nack(true)
					continue
				}
				_ = sm.Ack()
			}
		}
	}
}

func (c *RedisClient) buildStreamMessage(stream, group string, msg redis.XMessage) *streamMessage {
	payload := toBytes(msg.Values["payload"])
	key, _ := msg.Values["key"].(string)
	if key == "" {
		key = toString(msg.Values["key"])
	}
	ts := toInt64(msg.Values["ts"])
	if ts == 0 {
		ts = time.Now().UnixNano()
	}
	extraRaw := toString(msg.Values["extra"])
	var extra map[string]any
	if extraRaw != "" {
		_ = json.Unmarshal([]byte(extraRaw), &extra)
	}

	meta := mq.Metadata{
		Topic:     trimPrefix(stream, c.prefix),
		Key:       key,
		Partition: parsePartition(stream),
		MessageID: msg.ID,
		Timestamp: time.Unix(0, ts),
		Extra:     extra,
	}

	return &streamMessage{
		payload:  payload,
		meta:     meta,
		rdb:      c.rdb,
		stream:   stream,
		group:    group,
		id:       msg.ID,
		extraRaw: extraRaw,
		key:      key,
	}
}

func createGroupIfNeeded(ctx context.Context, rdb *redis.Client, stream, group string) error {
	err := rdb.XGroupCreateMkStream(ctx, stream, group, "$").Err()
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}

func toBytes(v any) []byte {
	switch t := v.(type) {
	case []byte:
		return t
	case string:
		return []byte(t)
	default:
		return nil
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case string:
		i, _ := strconv.ParseInt(t, 10, 64)
		return i
	case []byte:
		i, _ := strconv.ParseInt(string(t), 10, 64)
		return i
	default:
		return 0
	}
}

func trimPrefix(stream, prefix string) string {
	if strings.HasPrefix(stream, prefix) {
		stream = strings.TrimPrefix(stream, prefix)
	}
	// stream format: topic:partition
	if idx := strings.LastIndex(stream, ":"); idx > 0 {
		return stream[:idx]
	}
	return stream
}

func parsePartition(stream string) int {
	if idx := strings.LastIndex(stream, ":"); idx > 0 && idx < len(stream)-1 {
		p, _ := strconv.Atoi(stream[idx+1:])
		return p
	}
	return 0
}
