package xmap

import (
	"github.com/karosown/katool-go/container/stream"
)

// Map 是原生 map 的类型别名
// Map is a type alias for native map
type Map[K comparable, V any] map[K]V

// NewMap 创建并返回一个新的 Map 实例
// NewMap creates and returns a new Map instance
func NewMap[K comparable, V any]() Map[K, V] {
	return make(Map[K, V])
}

// CopyMap 从现有map复制创建新的Map
// CopyMap creates a new Map by copying from an existing map
func CopyMap[K comparable, V any, M ~map[K]V](m M) Map[K, V] {
	return Map[K, V](m)
}

// MapFromAny 从any类型的map转换为指定类型的Map
// MapFromAny converts a map of any type to a Map of specified types
func MapFromAny[K comparable, V any, M ~map[any]any](m M) Map[K, V] {
	mp := map[K]V{}
	for k, v := range m {
		k2, ok := k.(K)
		if !ok {
			continue
		}
		v2, ok := v.(V)
		if !ok {
			continue
		}
		mp[k2] = v2
	}
	return mp
}

// Get 获取指定键的值
// Get retrieves the value for a specified key
func (m Map[K, V]) Get(key K) (V, bool) {
	val, exists := m[key]
	return val, exists
}

// Set 设置键值对
// Set sets a key-value pair
func (m Map[K, V]) Set(key K, value V) {
	m[key] = value
}

// Delete 删除指定键值对
// Delete removes the specified key-value pair
func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

// Has 检查键是否存在
// Has checks if a key exists
func (m Map[K, V]) Has(key K) bool {
	_, exists := m[key]
	return exists
}

// Len 返回map中元素数量
// Len returns the number of elements in the map
func (m Map[K, V]) Len() int {
	return len(m)
}

// Keys 返回所有键的切片
// Keys returns a slice of all keys
func (m Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有值的切片
// Values returns a slice of all values
func (m Map[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Clear 清空map
// Clear empties the map
func (m Map[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
}

// Reset 重置map为新的空实例
// Reset resets the map to a new empty instance
func (m Map[K, V]) Reset() Map[K, V] {
	m = NewMap[K, V]()
	return m
}

// ForEach 遍历map中的每个键值对
// ForEach iterates over each key-value pair in the map
func (m Map[K, V]) ForEach(fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// ToStream 将map转换为Stream条目
// ToStream converts the map to Stream entries
func (m Map[K, V]) ToStream() stream.Entries[K, V] {
	return stream.EntrySet(m)
}
