package ioc

import (
	"sync"
)

var ic = sync.Map{}

func GetDef[V any](key string, value V) V {
	get := Get(key)
	if get == nil {
		RegisterValue(key, value)
		get = value
	}
	return get.(V)
}

func GetDefFunc[V any](key string, value func() V) V {
	get := Get(key)
	if get == nil {
		ForceRegister(key, value)
		get = Get(key)
	}
	return get.(V)
}
func Get(key string) any {
	if v, ok := ic.Load(key); ok {
		return v
	}
	return nil
}
func MustRegisterValue(key string, value any) {
	ic.Store(key, value)
}

func RegisterValue(key string, value any) {
	if Get(key) != nil {
		panic("ioc:key already exists")
	} else {
		ic.Store(key, value)
	}
}
func MustRegister[V any](key string, valueFunction func() V) {
	ic.Store(key, valueFunction())
}
func ForceRegister[V any](key string, valueFunction func() V) {
	if Get(key) != nil {
		ic.Delete(key)
	}
	ic.Store(key, valueFunction())
}
func Register(key string, valueFunction func() any) {
	if Get(key) != nil {
		panic("ioc:key already exists")
	} else {
		ic.Store(key, valueFunction())
	}
}
