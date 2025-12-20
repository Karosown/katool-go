package kfmq

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/karosown/katool-go/mq"
	"github.com/segmentio/kafka-go"
)

type kafkaMessage struct {
	payload []byte
	meta    mq.Metadata

	reader   *kafka.Reader
	origin   kafka.Message
	useCommit bool

	doneOnce sync.Once
	doneErr  error
}

func (m *kafkaMessage) Payload() []byte          { return m.payload }
func (m *kafkaMessage) GetMetadata() mq.Metadata { return m.meta }

func (m *kafkaMessage) Ack() error {
	return m.finalize(func() error {
		if !m.useCommit {
			return nil
		}
		return m.reader.CommitMessages(context.Background(), m.origin)
	})
}

func (m *kafkaMessage) Nack(requeue bool) error {
	return m.finalize(func() error {
		if m.useCommit {
			if requeue {
				return nil
			}
			return m.reader.CommitMessages(context.Background(), m.origin)
		}
		if requeue {
			return m.reader.SetOffset(m.origin.Offset)
		}
		return nil
	})
}

func (m *kafkaMessage) finalize(fn func() error) error {
	m.doneOnce.Do(func() {
		m.doneErr = fn()
	})
	return m.doneErr
}

func (c *KafkaClient) Subscribe(ctx context.Context, topic string, handler mq.Handler, opts ...mq.SubscribeOption) error {
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
	assignPartitions := len(options.SpecificPartitions) > 0
	if !assignPartitions && options.Group == "" {
		options.Group = "default_" + uuid.New().String()
	}

	if assignPartitions {
		for _, p := range options.SpecificPartitions {
			if p < 0 {
				return mq.ErrInvalidPartition
			}
			reader := c.newReader(topic, options.Group, options.ConsumerID, p, true)
			c.registerReader(reader)
			go c.consumeReader(ctx, reader, handler, options.Filter)
		}
		return nil
	}

	reader := c.newReader(topic, options.Group, options.ConsumerID, 0, false)
	c.registerReader(reader)
	go c.consumeReader(ctx, reader, handler, options.Filter)
	return nil
}

func (c *KafkaClient) consumeReader(ctx context.Context, reader *kafka.Reader, handler mq.Handler, filter mq.FilterFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			continue
		}

		meta := buildKafkaMetadata(msg)
		sm := &kafkaMessage{
			payload:  msg.Value,
			meta:     meta,
			reader:   reader,
			origin:   msg,
			useCommit: reader.Config().GroupID != "",
		}

		if filter != nil && !filter(meta) {
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

func (c *KafkaClient) newReader(topic, group, consumerID string, partition int, assign bool) *kafka.Reader {
	cfg := c.readerConfig
	if !c.hasReaderConfig {
		cfg = kafka.ReaderConfig{}
	}
	if len(cfg.Brokers) == 0 {
		cfg.Brokers = c.brokers
	}
	cfg.Topic = topic
	if assign {
		cfg.GroupID = ""
		cfg.GroupInstanceID = ""
		cfg.Partition = partition
	} else {
		cfg.GroupID = group
		if consumerID != "" {
			cfg.GroupInstanceID = consumerID
		}
		cfg.Partition = 0
	}
	if cfg.MinBytes == 0 {
		cfg.MinBytes = 1e3
	}
	if cfg.MaxBytes == 0 {
		cfg.MaxBytes = 10e6
	}
	if cfg.MaxWait == 0 {
		cfg.MaxWait = 5 * time.Second
	}
	return kafka.NewReader(cfg)
}

func buildKafkaMetadata(msg kafka.Message) mq.Metadata {
	extra := parseExtraHeader(msg.Headers)
	ts := msg.Time
	if ts.IsZero() {
		ts = time.Now()
	}
	return mq.Metadata{
		Topic:     msg.Topic,
		Key:       string(msg.Key),
		Partition: msg.Partition,
		MessageID: strconv.FormatInt(msg.Offset, 10),
		Timestamp: ts,
		Extra:     extra,
	}
}

func parseExtraHeader(headers []kafka.Header) map[string]any {
	for _, h := range headers {
		if h.Key != "extra" {
			continue
		}
		var extra map[string]any
		if err := json.Unmarshal(h.Value, &extra); err == nil {
			return extra
		}
		break
	}
	return nil
}
