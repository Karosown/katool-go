package xmap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

// Codec defines value encoding/decoding for RedisMap.
type Codec[V any] interface {
	Marshal(V) ([]byte, error)
	Unmarshal([]byte, *V) error
}

// JSONCodec is the default codec for RedisMap.
type JSONCodec[V any] struct{}

func (JSONCodec[V]) Marshal(v V) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec[V]) Unmarshal(b []byte, out *V) error {
	return json.Unmarshal(b, out)
}

// KeyCodec defines key encoding/decoding for RedisMap.
type KeyCodec[K any] interface {
	Encode(K) (string, error)
	Decode(string) (K, error)
}

// RedisMapOption customizes RedisMap.
type RedisMapOption[K comparable, V any] func(*RedisMap[K, V])

// WithRedisMapTTL sets a key TTL for the backing Redis hash.
func WithRedisMapTTL[K comparable, V any](ttl time.Duration) RedisMapOption[K, V] {
	return func(m *RedisMap[K, V]) {
		m.ttl = ttl
	}
}

// WithRedisMapCodec sets a custom value codec.
func WithRedisMapCodec[K comparable, V any](codec Codec[V]) RedisMapOption[K, V] {
	return func(m *RedisMap[K, V]) {
		if codec != nil {
			m.codec = codec
		}
	}
}

// WithRedisMapKeyCodec sets a custom key codec.
func WithRedisMapKeyCodec[K comparable, V any](codec KeyCodec[K]) RedisMapOption[K, V] {
	return func(m *RedisMap[K, V]) {
		if codec != nil {
			m.keyCodec = codec
		}
	}
}

// RedisMap is a Redis-backed map using a hash key.
type RedisMap[K comparable, V any] struct {
	client   redis.Cmdable
	key      string
	ttl      time.Duration
	codec    Codec[V]
	keyCodec KeyCodec[K]
}

// NewRedisMap creates a Redis-backed map stored as a hash key.
func NewRedisMap[K comparable, V any](client redis.Cmdable, key string, opts ...RedisMapOption[K, V]) *RedisMap[K, V] {
	m := &RedisMap[K, V]{
		client:   client,
		key:      key,
		codec:    defaultValueCodec[V](),
		keyCodec: defaultKeyCodec[K](),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Get retrieves a value for a key.
func (m *RedisMap[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	var zero V
	if ctx == nil {
		ctx = context.Background()
	}
	field, err := m.encodeKey(key)
	if err != nil {
		return zero, false, err
	}
	b, err := m.client.HGet(ctx, m.key, field).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, false, nil
		}
		return zero, false, err
	}
	if err := m.codec.Unmarshal(b, &zero); err != nil {
		return zero, false, err
	}
	return zero, true, nil
}

// Set stores a key-value pair.
func (m *RedisMap[K, V]) Set(ctx context.Context, key K, value V) error {
	if ctx == nil {
		ctx = context.Background()
	}
	field, err := m.encodeKey(key)
	if err != nil {
		return err
	}
	b, err := m.codec.Marshal(value)
	if err != nil {
		return err
	}
	if err := m.client.HSet(ctx, m.key, field, b).Err(); err != nil {
		return err
	}
	if m.ttl > 0 {
		_ = m.client.Expire(ctx, m.key, m.ttl).Err()
	}
	return nil
}

// Delete removes a key-value pair.
func (m *RedisMap[K, V]) Delete(ctx context.Context, key K) error {
	if ctx == nil {
		ctx = context.Background()
	}
	field, err := m.encodeKey(key)
	if err != nil {
		return err
	}
	return m.client.HDel(ctx, m.key, field).Err()
}

// Has checks if a key exists.
func (m *RedisMap[K, V]) Has(ctx context.Context, key K) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	field, err := m.encodeKey(key)
	if err != nil {
		return false, err
	}
	return m.client.HExists(ctx, m.key, field).Result()
}

// Len returns the number of elements in the map.
func (m *RedisMap[K, V]) Len(ctx context.Context) (int64, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return m.client.HLen(ctx, m.key).Result()
}

// Keys returns all keys when a key decoder is available.
func (m *RedisMap[K, V]) Keys(ctx context.Context) ([]K, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if m.keyCodec == nil {
		return nil, errors.New("xmap: key codec not set")
	}
	fields, err := m.client.HKeys(ctx, m.key).Result()
	if err != nil {
		return nil, err
	}
	keys := make([]K, 0, len(fields))
	for _, field := range fields {
		k, err := m.keyCodec.Decode(field)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// Values returns all values.
func (m *RedisMap[K, V]) Values(ctx context.Context) ([]V, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	vals, err := m.client.HVals(ctx, m.key).Result()
	if err != nil {
		return nil, err
	}
	out := make([]V, 0, len(vals))
	for _, v := range vals {
		var decoded V
		if err := m.codec.Unmarshal([]byte(v), &decoded); err != nil {
			return nil, err
		}
		out = append(out, decoded)
	}
	return out, nil
}

// GetAll loads all key-value pairs.
func (m *RedisMap[K, V]) GetAll(ctx context.Context) (map[K]V, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if m.keyCodec == nil {
		return nil, errors.New("xmap: key codec not set")
	}
	all, err := m.client.HGetAll(ctx, m.key).Result()
	if err != nil {
		return nil, err
	}
	out := make(map[K]V, len(all))
	for field, v := range all {
		k, err := m.keyCodec.Decode(field)
		if err != nil {
			return nil, err
		}
		var decoded V
		if err := m.codec.Unmarshal([]byte(v), &decoded); err != nil {
			return nil, err
		}
		out[k] = decoded
	}
	return out, nil
}

// Clear removes the backing Redis hash.
func (m *RedisMap[K, V]) Clear(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return m.client.Del(ctx, m.key).Err()
}

func (m *RedisMap[K, V]) encodeKey(key K) (string, error) {
	if m.keyCodec == nil {
		return "", errors.New("xmap: key codec not set")
	}
	return m.keyCodec.Encode(key)
}

func defaultKeyCodec[K any]() KeyCodec[K] {
	var zero K
	t := reflect.TypeOf(zero)
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.String {
		return stringKeyCodec[K]{}
	}
	return nil
}

type stringKeyCodec[K any] struct{}

func (stringKeyCodec[K]) Encode(k K) (string, error) {
	return fmt.Sprint(k), nil
}

func (stringKeyCodec[K]) Decode(s string) (K, error) {
	var zero K
	t := reflect.TypeOf(zero)
	if t == nil {
		return zero, errors.New("xmap: key type is nil")
	}
	v := reflect.New(t).Elem()
	if v.Kind() != reflect.String {
		return zero, errors.New("xmap: key decoder not set")
	}
	v.SetString(s)
	return v.Interface().(K), nil
}

func defaultValueCodec[V any]() Codec[V] {
	var zero V
	t := reflect.TypeOf(zero)
	if t == nil {
		return JSONCodec[V]{}
	}
	if t.Kind() == reflect.String {
		return stringValueCodec[V]{}
	}
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return bytesValueCodec[V]{}
	}
	return JSONCodec[V]{}
}

type stringValueCodec[V any] struct{}

func (stringValueCodec[V]) Marshal(v V) ([]byte, error) {
	return []byte(fmt.Sprint(v)), nil
}

func (stringValueCodec[V]) Unmarshal(b []byte, out *V) error {
	if out == nil {
		return errors.New("xmap: nil output for string codec")
	}
	t := reflect.TypeOf(*out)
	if t == nil || t.Kind() != reflect.String {
		return errors.New("xmap: string codec used with non-string type")
	}
	v := reflect.New(t).Elem()
	v.SetString(string(b))
	*out = v.Interface().(V)
	return nil
}

type bytesValueCodec[V any] struct{}

func (bytesValueCodec[V]) Marshal(v V) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice || rv.Type().Elem().Kind() != reflect.Uint8 {
		return nil, errors.New("xmap: bytes codec used with non-bytes type")
	}
	if rv.IsNil() {
		return nil, nil
	}
	out := make([]byte, rv.Len())
	reflect.Copy(reflect.ValueOf(out), rv)
	return out, nil
}

func (bytesValueCodec[V]) Unmarshal(b []byte, out *V) error {
	if out == nil {
		return errors.New("xmap: nil output for bytes codec")
	}
	t := reflect.TypeOf(*out)
	if t == nil || t.Kind() != reflect.Slice || t.Elem().Kind() != reflect.Uint8 {
		return errors.New("xmap: bytes codec used with non-bytes type")
	}
	v := reflect.MakeSlice(t, len(b), len(b))
	reflect.Copy(v, reflect.ValueOf(b))
	*out = v.Interface().(V)
	return nil
}
