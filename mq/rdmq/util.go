package rdmq

import (
	"hash/fnv"
	"strconv"
)

// Option customizes RedisClient.
type Option func(*RedisClient)

// WithPartitions sets a fixed partition count for all topics.
func WithPartitions(n int) Option {
	return func(c *RedisClient) {
		if n > 0 {
			c.partitions = n
		}
	}
}

// WithPrefix sets the stream key prefix.
func WithPrefix(prefix string) Option {
	return func(c *RedisClient) {
		if prefix != "" {
			c.prefix = prefix
		}
	}
}

func (c *RedisClient) streamKey(topic string, partition int) string {
	return c.prefix + topic + ":" + strconv.Itoa(partition)
}

func hashKey(key string, numParts int) int {
	if numParts == 0 || key == "" {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32()) % numParts
}
