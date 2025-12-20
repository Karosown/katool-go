package cmq

import (
	"context"
	"github.com/karosown/katool-go/mq"
	"sync"
)

// -------------------------------------------------------
// Message 实现
// -------------------------------------------------------
type chanMessage struct {
	payload []byte      `json:"Payload,omitempty"`
	meta    mq.Metadata `json:"Meta"`
}

func (m *chanMessage) Payload() []byte          { return m.payload }
func (m *chanMessage) GetMetadata() mq.Metadata { return m.meta }
func (m *chanMessage) Ack() error               { return nil } // 内存版不需要 Ack
func (m *chanMessage) Nack(requeue bool) error  { return nil }

// -------------------------------------------------------
// Broker 结构
// -------------------------------------------------------
type consumerGroup struct {
	name       string
	partitions []chan *chanMessage
}

type ChanBroker struct {
	mu sync.RWMutex
	// 数据结构：Topic -> GroupName -> GroupStruct
	queues          map[string]map[string]*consumerGroup
	topicPartitions map[string]int
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewChanBroker() *ChanBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChanBroker{
		queues:          make(map[string]map[string]*consumerGroup),
		topicPartitions: make(map[string]int),
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (b *ChanBroker) Close() error {
	b.cancel()
	return nil
}
