package xredis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetJSON stores a value as JSON.
func (c *Client) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, c.Key(key), b, ttl).Err()
}

// GetJSON loads JSON data into out and returns whether the key exists.
func (c *Client) GetJSON(ctx context.Context, key string, out any) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	b, err := c.Get(ctx, c.Key(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, json.Unmarshal(b, out)
}

// HSetJSON stores a hash field as JSON.
func (c *Client) HSetJSON(ctx context.Context, key, field string, value any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.HSet(ctx, c.Key(key), field, b).Err()
}

// HGetJSON loads a hash field from JSON and returns whether it exists.
func (c *Client) HGetJSON(ctx context.Context, key, field string, out any) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	b, err := c.HGet(ctx, c.Key(key), field).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, json.Unmarshal(b, out)
}
