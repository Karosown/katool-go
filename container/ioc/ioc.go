package ioc

import (
	"sync"

	"github.com/karosown/katool-go/sys"
)

// ic IOC容器的全局实例
// ic is the global instance of the IOC container
var ic = sync.Map{}

// GetDef 获取组件，如果不存在则注册默认值
// GetDef gets a component, registers a default value if it doesn't exist
func GetDef[V any](key string, value V) V {
	get := Get(key)
	if get == nil {
		RegisterValue(key, value)
		get = value
	}
	return get.(V)
}

// GetDefFunc 获取组件，如果不存在则使用函数创建
// GetDefFunc gets a component, creates it using a function if it doesn't exist
func GetDefFunc[V any](key string, value func() V) V {
	get := Get(key)
	if get == nil {
		ForceRegister(key, value)
		get = Get(key)
	}
	return get.(V)
}

// Get 获取组件
// Get retrieves a component
func Get(key string) any {
	if v, ok := ic.Load(key); ok {
		return v
	}
	return nil
}

// MustRegisterValue 强制注册值（覆盖已存在的）
// MustRegisterValue forcibly registers a value (overwrites existing)
func MustRegisterValue(key string, value any) {
	ic.Store(key, value)
}

// RegisterValue 注册值（如果已存在则抛出异常）
// RegisterValue registers a value (panics if already exists)
func RegisterValue(key string, value any) {
	if Get(key) != nil {
		sys.Panic("ioc:key already exists")
	} else {
		ic.Store(key, value)
	}
}

// MustRegister 强制注册工厂函数（覆盖已存在的）
// MustRegister forcibly registers a factory function (overwrites existing)
func MustRegister[V any](key string, valueFunction func() V) {
	ic.Store(key, valueFunction())
}

// ForceRegister 强制注册工厂函数（删除已存在的）
// ForceRegister forcibly registers a factory function (deletes existing)
func ForceRegister[V any](key string, valueFunction func() V) {
	if Get(key) != nil {
		ic.Delete(key)
	}
	ic.Store(key, valueFunction())
}

// Register 注册工厂函数（如果已存在则抛出异常）
// Register registers a factory function (panics if already exists)
func Register(key string, valueFunction func() any) {
	if Get(key) != nil {
		sys.Panic("ioc:key already exists")
	} else {
		ic.Store(key, valueFunction())
	}
}
