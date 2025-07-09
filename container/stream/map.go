package stream

import (
	"github.com/karosown/katool-go/convert"
)

// Entry 键值对条目
// Entry represents a key-value pair entry
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Entries 键值对条目集合
// Entries represents a collection of key-value pair entries
type Entries[K comparable, V any] []Entry[K, V]

// EntrySet 从映射创建条目集合
// EntrySet creates entry set from map
func EntrySet[K comparable, V any, Map ~map[K]V](m Map) Entries[K, V] {
	var entries []Entry[K, V]
	for k, v := range m {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return entries
}

// Identity 获取条目的身份副本
// Identity gets an identity copy of entries
func (e Entries[K, V]) Identity() *[]Entry[K, V] {
	convert := make([]Entry[K, V], len(e))
	for i := 0; i < len(e); i++ {
		convert[i] = Entry[K, V]{Key: e[i].Key, Value: e[i].Value}
	}
	return &convert
}

// KeySet 获取所有键的集合
// KeySet gets a set of all keys
func (e Entries[K, V]) KeySet() []K {
	keyset := e.ToStream().Map(func(i Entry[K, V]) any {
		return i.Key
	}).ToList()
	return convert.FromAnySlice[K](keyset)
}

// Values 获取所有值的集合
// Values gets a set of all values
func (e Entries[K, V]) Values() []V {
	keyset := e.ToStream().Map(func(i Entry[K, V]) any {
		return i.Value
	}).ToList()
	return convert.FromAnySlice[V](keyset)
}

// KeySetStream 获取键的流
// KeySetStream gets a stream of keys
func (e Entries[K, V]) KeySetStream() *Stream[K, []K] {
	ks := e.KeySet()
	return ToStream(&ks)
}

// ValuesStream 获取值的流
// ValuesStream gets a stream of values
func (e Entries[K, V]) ValuesStream() *Stream[V, []V] {
	ks := e.Values()
	return ToStream(&ks)
}

// ToStream 转换为流
// ToStream converts to stream
func (e Entries[K, V]) ToStream() *Stream[Entry[K, V], []Entry[K, V]] {
	return ToStream[Entry[K, V], []Entry[K, V]](e.Identity())
}

// ToParallelStream 转换为并行流
// ToParallelStream converts to parallel stream
func (e Entries[K, V]) ToParallelStream() *Stream[Entry[K, V], []Entry[K, V]] {
	return ToParallelStream[Entry[K, V], []Entry[K, V]](e.Identity())
}
