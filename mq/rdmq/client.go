package rdmq

import (
	"context"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

const (
	defaultPartitions = 4
	defaultPrefix     = "mq:"
)

// RedisClient implements mq.Client on top of Redis Streams.
type RedisClient struct {
	rdb        *redis.Client
	partitions int
	prefix     string

	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool
}

// NewRedisClient creates a Redis-backed mq client using go-redis options.
func NewRedisClient(addr, password string, db int, opts ...Option) *RedisClient {
	redisOpts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	return NewRedisClientWithOptions(redisOpts, opts...)
}

// NewRedisClientWithOptions creates a Redis-backed mq client using go-redis options.
func NewRedisClientWithOptions(redisOpts *redis.Options, opts ...Option) *RedisClient {
	ctx, cancel := context.WithCancel(context.Background())
	c := &RedisClient{
		rdb:        redis.NewClient(redisOpts),
		partitions: defaultPartitions,
		prefix:     defaultPrefix,
		ctx:        ctx,
		cancel:     cancel,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Close stops background work and closes the underlying Redis client.
func (c *RedisClient) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil
	}
	c.cancel()
	return c.rdb.Close()
}
