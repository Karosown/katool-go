package kfmq

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/segmentio/kafka-go"
)

// KafkaClient implements mq.Client on top of Kafka.
type KafkaClient struct {
	brokers []string
	writer  *kafka.Writer

	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool

	mu      sync.Mutex
	readers []*kafka.Reader

	readerConfig    kafka.ReaderConfig
	hasReaderConfig bool
}

// Option customizes KafkaClient.
type Option func(*KafkaClient)

// WithWriter replaces the default writer.
func WithWriter(w *kafka.Writer) Option {
	return func(c *KafkaClient) {
		if w != nil {
			c.writer = w
		}
	}
}

// WithWriterConfig builds a writer from the provided config.
func WithWriterConfig(cfg kafka.WriterConfig) Option {
	return func(c *KafkaClient) {
		c.writer = kafka.NewWriter(cfg)
	}
}

// WithReaderConfig sets the default reader config used for subscriptions.
// Topic/Group/Partition/Brokers fields are overridden per subscription.
func WithReaderConfig(cfg kafka.ReaderConfig) Option {
	return func(c *KafkaClient) {
		c.readerConfig = cfg
		c.hasReaderConfig = true
	}
}

// NewKafkaClient creates a Kafka-backed mq client.
func NewKafkaClient(brokers []string, opts ...Option) *KafkaClient {
	ctx, cancel := context.WithCancel(context.Background())
	c := &KafkaClient{
		brokers: brokers,
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.Hash{},
		},
		ctx:    ctx,
		cancel: cancel,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Close stops background work and closes the writer/readers.
func (c *KafkaClient) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil
	}
	c.cancel()
	c.mu.Lock()
	readers := c.readers
	c.readers = nil
	c.mu.Unlock()

	for _, r := range readers {
		_ = r.Close()
	}
	return c.writer.Close()
}

func (c *KafkaClient) registerReader(r *kafka.Reader) {
	c.mu.Lock()
	c.readers = append(c.readers, r)
	c.mu.Unlock()
}
