package xmap

import (
	"fmt"

	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/container/stream"
	"github.com/spf13/cast"
)

type KV[K any, V any] struct {
	Key   K
	Value V
}
type KMap[K any, V any] struct {
	// 使用这个Map来做Hash映射
	internalMap Map[uint32, *KV[K, V]]
}

func NewKMap[K any, V any]() *KMap[K, V] {
	return &KMap[K, V]{}
}
func hash[K any](key *K) uint32 {
	return cast.ToUint32(fmt.Sprintf("%p", key))
}
func (m *KMap[K, V]) Set(k K, v V) {
	u := hash(&k)
	m.internalMap.Set(u, &KV[K, V]{k, v})
}

func (m *KMap[K, V]) Get(k K) (V, bool) {
	u, ok := m.internalMap.Get(hash(&k))
	if !ok {
		return *new(V), false
	}
	return u.Value, ok
}

func (m *KMap[K, V]) Delete(k K) {
	u, ok := m.internalMap.Get(hash(&k))
	if !ok || nil == optional.ToNillable(u) {
		return
	}
	m.internalMap.Delete(hash(&k))
}

func (m *KMap[K, V]) Has(k K) bool {
	u, ok := m.internalMap.Get(hash(&k))
	if !ok {
		return false
	}
	return nil != u
}
func (m *KMap[K, V]) Len() int {
	return len(m.internalMap)
}
func (m *KMap[K, V]) ToStream() *stream.Stream[*KV[K, V], []*KV[K, V]] {
	values := m.internalMap.Values()
	return stream.Of(&values)
}

func (m *KMap[K, V]) Foreach(f func(k K, v V)) {
	m.ToStream().ForEach(func(item *KV[K, V]) {
		f(item.Key, item.Value)
	})
}

func (m *KMap[K, V]) Keys() []K {
	return stream.Cast[K](m.ToStream().Map(func(i *KV[K, V]) any {
		return i.Key
	})).ToList()
}

func (m *KMap[K, V]) Values() []V {
	return stream.Cast[V](m.ToStream().Map(func(i *KV[K, V]) any {
		return i.Value
	})).ToList()
}

func (m *KMap[K, V]) Clear() {
	for k := range m.internalMap {
		delete(m.internalMap, k)
	}
}
func (m *KMap[K, V]) Reset() *KMap[K, V] {
	m.internalMap = NewMap[uint32, *KV[K, V]]()
	return m
}
