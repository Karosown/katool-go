package cmq

import "hash/fnv"

const DefaultPartitions = 4

func (b *ChanBroker) ensureTopic(topic string) int {
	if _, ok := b.topicPartitions[topic]; !ok {
		b.topicPartitions[topic] = DefaultPartitions
	}
	if _, ok := b.queues[topic]; !ok {
		b.queues[topic] = make(map[string]*consumerGroup)
	}
	return b.topicPartitions[topic]
}

func (b *ChanBroker) hashKey(key string, numParts int) int {
	if numParts == 0 || key == "" {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % numParts
}
