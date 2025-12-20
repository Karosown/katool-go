package xredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Serializer encodes/decodes values stored in Redis.
type Serializer interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, out any) error
}

// JSONSerializer encodes values as JSON.
type JSONSerializer struct{}

func (JSONSerializer) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONSerializer) Decode(data []byte, out any) error {
	if out == nil {
		return errors.New("xredis: nil output")
	}
	return json.Unmarshal(data, out)
}

// StringSerializer encodes values as string.
type StringSerializer struct{}

func (StringSerializer) Encode(v any) ([]byte, error) {
	switch t := v.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		return []byte(fmt.Sprint(v)), nil
	}
}

func (StringSerializer) Decode(data []byte, out any) error {
	if out == nil {
		return errors.New("xredis: nil output")
	}
	switch t := out.(type) {
	case *string:
		*t = string(data)
		return nil
	case *[]byte:
		*t = append((*t)[:0], data...)
		return nil
	default:
		return errors.New("xredis: StringSerializer requires *string or *[]byte")
	}
}

// BytesSerializer stores raw bytes.
type BytesSerializer struct{}

func (BytesSerializer) Encode(v any) ([]byte, error) {
	switch t := v.(type) {
	case []byte:
		return t, nil
	case string:
		return []byte(t), nil
	default:
		return json.Marshal(v)
	}
}

func (BytesSerializer) Decode(data []byte, out any) error {
	if out == nil {
		return errors.New("xredis: nil output")
	}
	switch t := out.(type) {
	case *[]byte:
		*t = append((*t)[:0], data...)
		return nil
	case *string:
		*t = string(data)
		return nil
	default:
		return json.Unmarshal(data, out)
	}
}

// TemplateOption customizes a RedisTemplate.
type TemplateOption func(*RedisTemplate)

// WithTemplatePrefix sets a key prefix for all operations.
func WithTemplatePrefix(prefix string) TemplateOption {
	return func(t *RedisTemplate) {
		t.prefix = prefix
	}
}

// WithTemplateSerializer sets the value serializer.
func WithTemplateSerializer(s Serializer) TemplateOption {
	return func(t *RedisTemplate) {
		if s != nil {
			t.serializer = s
		}
	}
}

// RedisTemplate is a small, Java-style convenience wrapper.
type RedisTemplate struct {
	client     redis.Cmdable
	prefix     string
	serializer Serializer
}

// NewRedisTemplate creates a template for a redis client.
func NewRedisTemplate(client redis.Cmdable, opts ...TemplateOption) *RedisTemplate {
	t := &RedisTemplate{
		client:     client,
		serializer: JSONSerializer{},
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Key applies prefix to the key.
func (t *RedisTemplate) Key(key string) string {
	if t == nil || t.prefix == "" {
		return key
	}
	return t.prefix + key
}

func (t *RedisTemplate) ensureCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// Set stores a value with optional TTL.
func (t *RedisTemplate) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx = t.ensureCtx(ctx)
	b, err := t.serializer.Encode(value)
	if err != nil {
		return err
	}
	return t.client.Set(ctx, t.Key(key), b, ttl).Err()
}

// Get loads a value into out and reports if the key exists.
func (t *RedisTemplate) Get(ctx context.Context, key string, out any) (bool, error) {
	ctx = t.ensureCtx(ctx)
	b, err := t.client.Get(ctx, t.Key(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, t.serializer.Decode(b, out)
}

// GetBytes returns raw bytes.
func (t *RedisTemplate) GetBytes(ctx context.Context, key string) ([]byte, bool, error) {
	ctx = t.ensureCtx(ctx)
	b, err := t.client.Get(ctx, t.Key(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return b, true, nil
}

// Exists checks if key exists.
func (t *RedisTemplate) Exists(ctx context.Context, key string) (bool, error) {
	ctx = t.ensureCtx(ctx)
	n, err := t.client.Exists(ctx, t.Key(key)).Result()
	return n > 0, err
}

// Del deletes keys.
func (t *RedisTemplate) Del(ctx context.Context, keys ...string) (int64, error) {
	ctx = t.ensureCtx(ctx)
	if len(keys) == 0 {
		return 0, nil
	}
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, t.Key(k))
	}
	return t.client.Del(ctx, out...).Result()
}

// Expire sets key TTL.
func (t *RedisTemplate) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.Expire(ctx, t.Key(key), ttl).Result()
}

// TTL returns key TTL.
func (t *RedisTemplate) TTL(ctx context.Context, key string) (time.Duration, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.TTL(ctx, t.Key(key)).Result()
}

// Incr increments a numeric key.
func (t *RedisTemplate) Incr(ctx context.Context, key string) (int64, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.Incr(ctx, t.Key(key)).Result()
}

// IncrBy increments a numeric key by delta.
func (t *RedisTemplate) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.IncrBy(ctx, t.Key(key), delta).Result()
}

// HSet stores a field in a hash.
func (t *RedisTemplate) HSet(ctx context.Context, key, field string, value any) error {
	ctx = t.ensureCtx(ctx)
	b, err := t.serializer.Encode(value)
	if err != nil {
		return err
	}
	return t.client.HSet(ctx, t.Key(key), field, b).Err()
}

// HGet loads a hash field into out and reports if the field exists.
func (t *RedisTemplate) HGet(ctx context.Context, key, field string, out any) (bool, error) {
	ctx = t.ensureCtx(ctx)
	b, err := t.client.HGet(ctx, t.Key(key), field).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, t.serializer.Decode(b, out)
}

// HDel deletes fields from a hash.
func (t *RedisTemplate) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.HDel(ctx, t.Key(key), fields...).Result()
}

// HExists checks if a hash field exists.
func (t *RedisTemplate) HExists(ctx context.Context, key, field string) (bool, error) {
	ctx = t.ensureCtx(ctx)
	return t.client.HExists(ctx, t.Key(key), field).Result()
}

// HGetAll returns all hash fields as raw bytes.
func (t *RedisTemplate) HGetAll(ctx context.Context, key string) (map[string][]byte, error) {
	ctx = t.ensureCtx(ctx)
	all, err := t.client.HGetAll(ctx, t.Key(key)).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[string][]byte, len(all))
	for k, v := range all {
		out[k] = []byte(v)
	}
	return out, nil
}

// Do passes through a raw redis command.
func (t *RedisTemplate) Do(ctx context.Context, args ...any) *redis.Cmd {
	ctx = t.ensureCtx(ctx)
	if t == nil || t.client == nil {
		cmd := redis.NewCmd(ctx, args...)
		cmd.SetErr(errors.New("xredis: nil client"))
		return cmd
	}
	if d, ok := t.client.(interface {
		Do(context.Context, ...any) *redis.Cmd
	}); ok {
		return d.Do(ctx, args...)
	}
	cmd := redis.NewCmd(ctx, args...)
	cmd.SetErr(errors.New("xredis: underlying client does not support Do"))
	return cmd
}

var _ = JSONSerializer{}
