package xredis

import "github.com/redis/go-redis/v9"

// Client wraps go-redis client with small helpers.
type Client struct {
	*redis.Client
	prefix string
}

// Option customizes the redis client wrapper.
type Option func(*Client)

// WithPrefix sets a key prefix for helper methods.
func WithPrefix(prefix string) Option {
	return func(c *Client) {
		c.prefix = prefix
	}
}

// NewClient creates a wrapped redis client.
func NewClient(opts *redis.Options, opt ...Option) *Client {
	c := &Client{
		Client: redis.NewClient(opts),
	}
	for _, o := range opt {
		o(c)
	}
	return c
}

// NewClientFromURL parses a redis URL and creates a wrapped client.
func NewClientFromURL(url string, opt ...Option) (*Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return NewClient(opts, opt...), nil
}

// MustNewClientFromURL panics when URL parsing fails.
func MustNewClientFromURL(url string, opt ...Option) *Client {
	c, err := NewClientFromURL(url, opt...)
	if err != nil {
		panic(err)
	}
	return c
}

// Key applies the configured prefix to the key.
func (c *Client) Key(key string) string {
	if c == nil || c.prefix == "" {
		return key
	}
	return c.prefix + key
}

// Template returns a RedisTemplate using this client.
func (c *Client) Template(opts ...TemplateOption) *RedisTemplate {
	if c == nil {
		return NewRedisTemplate(nil, opts...)
	}
	if c.prefix != "" {
		opts = append(opts, WithTemplatePrefix(c.prefix))
	}
	return NewRedisTemplate(c.Client, opts...)
}
