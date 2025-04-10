package xtype

type AbstractMap[K comparable, V any] interface {
	map[K]V
}
