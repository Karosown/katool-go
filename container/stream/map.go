package stream

import (
	"github.com/karosown/katool-go/convert"
)

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
type Entries[K comparable, V any] []Entry[K, V]

func EntrySet[K comparable, V any, Map ~map[K]V](m Map) Entries[K, V] {
	var entries []Entry[K, V]
	for k, v := range m {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return entries
}
func (e Entries[K, V]) Identity() *[]Entry[K, V] {
	convert := make([]Entry[K, V], len(e))
	for i := 0; i < len(e); i++ {
		convert[i] = Entry[K, V]{Key: e[i].Key, Value: e[i].Value}
	}
	return &convert
}
func (e Entries[K, V]) KeySet() []K {
	keyset := e.ToStream().Map(func(i Entry[K, V]) any {
		return i.Key
	}).ToList()
	return convert.FromAnySlice[K](keyset)
}
func (e Entries[K, V]) Values() []V {
	keyset := e.ToStream().Map(func(i Entry[K, V]) any {
		return i.Value
	}).ToList()
	return convert.FromAnySlice[V](keyset)
}
func (e Entries[K, V]) KeySetStream() *Stream[K, []K] {
	ks := e.KeySet()
	return ToStream(&ks)
}
func (e Entries[K, V]) ValuesStream() *Stream[V, []V] {
	ks := e.Values()
	return ToStream(&ks)
}
func (e Entries[K, V]) ToStream() *Stream[Entry[K, V], []Entry[K, V]] {
	return ToStream[Entry[K, V], []Entry[K, V]](e.Identity())
}

func (e Entries[K, V]) ToParallelStream() *Stream[Entry[K, V], []Entry[K, V]] {
	return ToParallelStream[Entry[K, V], []Entry[K, V]](e.Identity())
}
