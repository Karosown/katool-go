package stream

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
type Entries[K comparable, V any] []Entry[K, V]

func EntrySet[K comparable, V any](m map[K]V) Entries[K, V] {
	var entries []Entry[K, V]
	for k, v := range m {
		entries = append(entries, Entry[K, V]{Key: k, Value: v})
	}
	return entries
}

func (e Entries[K, V]) ToStream() *Stream[Entry[K, V], []Entry[K, V]] {
	convert := make([]Entry[K, V], len(e))
	for i := 0; i < len(e); i++ {
		convert[i] = Entry[K, V]{Key: e[i].Key, Value: e[i].Value}
	}
	return ToStream[Entry[K, V], []Entry[K, V]](&convert)
}

func (e Entries[K, V]) ToParallelStream() *Stream[Entry[K, V], []Entry[K, V]] {
	convert := make([]Entry[K, V], len(e))
	for i := 0; i < len(e); i++ {
		convert[i] = Entry[K, V]{Key: e[i].Key, Value: e[i].Value}
	}
	return ToParallelStream[Entry[K, V], []Entry[K, V]](&convert)
}
