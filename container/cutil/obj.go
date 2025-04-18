package cutil

import (
	"reflect"
)

func IsBlank[T comparable](obj T) bool {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return true
		}
		for i := 0; i < v.Len(); i++ {
			if !IsBlank(v.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Map:
		return v.Len() == 0
	case reflect.Chan:
		return v.IsNil()
	case reflect.Ptr:
		return v.IsNil() || IsBlank(v.Elem().Interface()) // Calls IsBlank recursively for pointer types
	default:
		return obj == *new(T) // Check for zero value of the type
	}
}
