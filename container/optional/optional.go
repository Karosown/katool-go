package optional

import (
	"cmp"
	"slices"
)

func IsTrue[T any](valid bool, isTrue T, isFalse T) T {
	if valid {
		return isTrue
	} else {
		return isFalse
	}
}

func IsTrueByFunc[T any](valid bool, isTrue func() T, isFalse func() T) T {
	if valid {
		return isTrue()
	} else {
		return isFalse()
	}
}

func FuncIsTrueByFunc[T any](valid func() bool, isTrue, isFalse func() T) T {
	if valid() {
		return isTrue()
	} else {
		return isFalse()
	}
}

func EmptyStringFunc() string {
	return ""
}

func Identity[T any](obj T) func() T {
	return func() T {
		return obj
	}
}
func In[T cmp.Ordered](item T, arr ...T) bool {
	search := slices.Contains(arr, item)
	return search
}
