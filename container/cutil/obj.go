package cutil

import "reflect"

func IsBlank[T comparable](obj T) bool {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
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
