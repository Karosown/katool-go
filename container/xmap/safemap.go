package xmap

import (
	"sync"

	"github.com/karosown/katool-go/container/stream"
)

// SafeMap 是对 sync.Map 的泛型封装，提供线程安全的映射操作
// SafeMap is a generic wrapper for sync.Map, providing thread-safe mapping operations
type SafeMap[K comparable, V any] struct {
	internal sync.Map
}

// NewSafeMap 创建一个新的泛型 SafeMap
// NewSafeMap creates a new generic SafeMap
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{}
}

// CopySafeMap 从现有map复制创建新的SafeMap
// CopySafeMap creates a new SafeMap by copying from an existing map
func CopySafeMap[K comparable, V any, M ~map[K]V](mp M) *SafeMap[K, V] {
	internal := sync.Map{}
	for k, v := range mp {
		internal.Store(k, v)
	}
	return &SafeMap[K, V]{
		internal: internal,
	}
}

// SafeMapFromAny 从any类型的map转换为指定类型的SafeMap
// SafeMapFromAny converts a map of any type to a SafeMap of specified types
func SafeMapFromAny[K comparable, V any, M ~map[any]any](m M) *SafeMap[K, V] {
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
	return CopySafeMap(mp)
}

// Get 获取指定键的值
// Get retrieves the value for a specified key
func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	var zero V
	value, ok := m.internal.Load(key)
	if !ok {
		return zero, false
	}
	return value.(V), true
}

// Set 设置键值对
// Set sets a key-value pair
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.internal.Store(key, value)
}

// Delete 删除指定键值对
// Delete removes the specified key-value pair
func (m *SafeMap[K, V]) Delete(key K) {
	m.internal.Delete(key)
}

// Has 检查键是否存在
// Has checks if a key exists
func (m *SafeMap[K, V]) Has(key K) bool {
	_, ok := m.internal.Load(key)
	return ok
}

// LoadOrStore 尝试获取现有值，如果不存在则存储并返回提供的值
// LoadOrStore attempts to get an existing value, or stores and returns the provided value if it doesn't exist
func (m *SafeMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := m.internal.LoadOrStore(key, value)
	return actual.(V), loaded
}

// LoadAndDelete 加载值并删除键
// LoadAndDelete loads a value and deletes the key
func (m *SafeMap[K, V]) LoadAndDelete(key K) (V, bool) {
	var zero V
	value, loaded := m.internal.LoadAndDelete(key)
	if !loaded {
		return zero, false
	}
	return value.(V), loaded
}

// Range 遍历所有键值对
// Range iterates over all key-value pairs
func (m *SafeMap[K, V]) Range(fn func(K, V) bool) {
	m.internal.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}

// Clear 清空 Map 中的所有元素
// Clear removes all elements from the Map
func (m *SafeMap[K, V]) Clear() {
	m.internal.Clear()
}

// Reset 通过创建新的 sync.Map 实例来清空
// Reset clears by creating a new sync.Map instance
func (m *SafeMap[K, V]) Reset() {
	m.internal = sync.Map{}
}

// ToMap 转换为普通Map
// ToMap converts to a regular Map
func (m *SafeMap[K, V]) ToMap() Map[K, V] {
	mapper := Map[K, V]{}
	m.Range(func(k K, v V) bool {
		mapper[k] = v
		return true
	})
	return mapper
}

// ToStream 将SafeMap转换为Stream条目
// ToStream converts the SafeMap to Stream entries
func (m *SafeMap[K, V]) ToStream() stream.Entries[K, V] {
	return stream.EntrySet(m.ToMap())
}

// Len 返回 Map 中元素的数量
// Len returns the number of elements in the Map
func (m *SafeMap[K, V]) Len() int {
	count := 0
	m.Range(func(_ K, _ V) bool {
		count++
		return true
	})
	return count
}
