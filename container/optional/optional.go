package optional

import (
	"cmp"
	"errors"
	"slices"
)

// IsTrue 根据条件返回不同的值
// IsTrue returns different values based on a condition
func IsTrue[T any](valid bool, isTrue T, isFalse T) T {
	if valid {
		return isTrue
	} else {
		return isFalse
	}
}

// IsTrueByFunc 根据条件调用不同的函数
// IsTrueByFunc calls different functions based on a condition
func IsTrueByFunc[T any](valid bool, isTrue func() T, isFalse func() T) T {
	if valid {
		return isTrue()
	} else {
		return isFalse()
	}
}

// FuncIsTrueByFunc 根据函数条件调用不同的函数
// FuncIsTrueByFunc calls different functions based on a function condition
func FuncIsTrueByFunc[T any](valid func() bool, isTrue, isFalse func() T) T {
	if valid() {
		return isTrue()
	} else {
		return isFalse()
	}
}

// EmptyStringFunc 返回空字符串的函数
// EmptyStringFunc returns a function that returns an empty string
func EmptyStringFunc() string {
	return ""
}

// Identity 返回一个返回指定对象的函数
// Identity returns a function that returns the specified object
func Identity[T any](obj T) func() T {
	return func() T {
		return obj
	}
}

// IdentityErr 返回一个返回指定对象和错误的函数
// IdentityErr returns a function that returns the specified object and error
func IdentityErr[T any](obj T, errs ...error) func() (T, error) {
	return func() (T, error) {
		return obj, errors.Join(errs...)
	}
}

// IdentityOnlyErr 返回一个只返回错误的函数
// IdentityOnlyErr returns a function that only returns an error
func IdentityOnlyErr[T any](errs ...error) func() (T, error) {
	return func() (T, error) {
		t := new(T)
		return *t, errors.Join(errs...)
	}
}

// In 检查元素是否在数组中
// In checks if an element is in an array
func In[T cmp.Ordered](item T, arr ...T) bool {
	search := slices.Contains(arr, item)
	return search
}
