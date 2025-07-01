package optional

import (
	"cmp"
	"errors"
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
func IdentityErr[T any](obj T, errs ...error) func() (T, error) {
	return func() (T, error) {
		return obj, errors.Join(errs...)
	}
}

func IdentityOnlyErr[T any](errs ...error) func() (T, error) {
	return func() (T, error) {
		t := new(T)
		return *t, errors.Join(errs...)
	}
}
func In[T cmp.Ordered](item T, arr ...T) bool {
	search := slices.Contains(arr, item)
	return search
}
