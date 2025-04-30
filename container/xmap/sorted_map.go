package xmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/karosown/katool-go/container/stream"
	"golang.org/x/exp/constraints"
)

// SortedMap 是对 Map 类型的包装，同时提供有序遍历的能力。
// 要求键类型满足 constraints.Ordered 才能排序。
type SortedMap[K constraints.Ordered, V any] struct {
	mapper Map[K, V]
}

// NewSortedMap 创建一个新的 SortedMap 实例
func NewSortedMap[K constraints.Ordered, V any]() *SortedMap[K, V] {
	return &SortedMap[K, V]{
		mapper: NewMap[K, V](),
	}
}
func CopySortedMap[K constraints.Ordered, V any, M ~map[K]V](mp M) *SortedMap[K, V] {
	return &SortedMap[K, V]{
		mapper: Map[K, V](mp),
	}
}
func SortedMapFromAny[K constraints.Ordered, V any, M ~map[any]any](m M) *SortedMap[K, V] {
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
	return CopySortedMap(mp)
}

// Set 添加或更新一个键值对
func (sm *SortedMap[K, V]) Set(key K, value V) {
	sm.mapper.Set(key, value)
}

// Get 根据键获取值
func (sm *SortedMap[K, V]) Get(key K) (V, bool) {
	return sm.mapper.Get(key)
}

// Delete 删除指定键的键值对
func (sm *SortedMap[K, V]) Delete(key K) {
	sm.mapper.Delete(key)
}

// Len 返回 SortedMap 中元素数量
func (sm *SortedMap[K, V]) Len() int {
	return sm.mapper.Len()
}

// SortedKeys 返回排序后的所有键（按照自然顺序排序）
func (sm *SortedMap[K, V]) SortedKeys() []K {
	keys := sm.mapper.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
func (sm *SortedMap[K, V]) ToStream() *stream.Stream[stream.Entry[K, V], []stream.Entry[K, V]] {
	return stream.EntrySet(sm.mapper).ToStream()
}

// MarshalJSON 实现 json.Marshaler 接口，确保按照 SortedKeys 顺序输出
// 输出格式为 {"key": value, ...}
func (sm *SortedMap[K, V]) MarshalJSON() ([]byte, error) {
	keys := sm.SortedKeys()
	var buf bytes.Buffer

	buf.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}

		// 将 key 转换为字符串，
		// 如果 key 的底层类型是 string 则直接使用，否则使用 fmt.Sprintf 转换
		var keyStr string
		if ks, ok := any(key).(string); ok {
			keyStr = ks
		} else {
			keyStr = fmt.Sprintf("%v", key)
		}

		// 序列化 keyStr 确保正确转义
		kBytes, err := json.Marshal(keyStr)
		if err != nil {
			return nil, err
		}
		buf.Write(kBytes)
		buf.WriteByte(':')

		// 获取并序列化 value
		value, ok := sm.mapper.Get(key)
		if !ok {
			// 如果不存在 value，写入 null
			buf.WriteString("null")
		} else {
			vBytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			buf.Write(vBytes)
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
