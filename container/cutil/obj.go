package cutil

import (
	"reflect"
	"strconv"
	"unicode"
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
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func IsAllDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
func IsNumericAdvanced(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// 检查是否是空字符串
		if s == "" {
			return false
		}
		// 检查是否包含除数字和'.'的其他字符
		for _, r := range s {
			if !unicode.IsDigit(r) && r != '.' {
				return false
			}
		}
		return true
	}
	return true
}
