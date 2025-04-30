package xmap

import (
	"github.com/karosown/katool-go/container/stream"
)

// Map 是原生 map 的类型别名
type Map[K comparable, V any] map[K]V

// NewMap 创建并返回一个新的 Map 实例
func NewMap[K comparable, V any]() Map[K, V] {
	return make(Map[K, V])
}
func CopyMap[K comparable, V any, M ~map[K]V](m M) Map[K, V] {
	return Map[K, V](m)
}
func MapFromAny[K comparable, V any, M ~map[any]any](m M) Map[K, V] {
	mp := map[K]V{}
	for k, v := range m {
		mp[k] = v
	}
	return mp
}

// Get 获取指定键的值
func (m Map[K, V]) Get(key K) (V, bool) {
	val, exists := m[key]
	return val, exists
}

// Set 设置键值对
func (m Map[K, V]) Set(key K, value V) {
	m[key] = value
}

// Delete 删除指定键值对
func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

// Has 检查键是否存在
func (m Map[K, V]) Has(key K) bool {
	_, exists := m[key]
	return exists
}

// Len 返回map中元素数量
func (m Map[K, V]) Len() int {
	return len(m)
}

// Keys 返回所有键的切片
func (m Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有值的切片
func (m Map[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Clear 清空map
func (m Map[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
}
func (m Map[K, V]) Reset() Map[K, V] {
	m = NewMap[K, V]()
	return m
}

// ForEach 遍历map中的每个键值对
func (m Map[K, V]) ForEach(fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

func (m Map[K, V]) ToStream() stream.Entries[K, V] {
	return stream.EntrySet(m)
}
