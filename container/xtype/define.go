package xtype

import (
	"github.com/karosown/katool/container/xmap"
)

type AbstractMap[K comparable, V any] interface {
	map[K]V | xmap.Map[K, V]
}
