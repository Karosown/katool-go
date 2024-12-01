package ioc

import (
	"sync"
)

var ic = sync.Map{}

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
func MustRegister(key string, valueFunction func() any) {
	ic.Store(key, valueFunction())
}

func Register(key string, valueFunction func() any) {
	if Get(key) != nil {
		panic("ioc:key already exists")
	} else {
		ic.Store(key, valueFunction())
	}
}
