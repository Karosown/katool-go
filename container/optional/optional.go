package optional

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// Optional 可选值容器，用于安全处理可能为空的值
// Optional is a container for values that may or may not be present
type Optional[T any] struct {
	value   T    // 存储的值 / Stored value
	present bool // 是否存在值 / Whether value is present
}

// Of 创建一个包含指定值的Optional
// Of creates an Optional containing the specified value
func Of[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}

// Empty 创建一个空的Optional
// Empty creates an empty Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

// OfNullable 根据值是否为零值创建Optional
// OfNullable creates Optional based on whether value is zero value
func OfNullable[T comparable](value T) Optional[T] {
	var zero T
	if value == zero {
		return Empty[T]()
	}
	return Of(value)
}

// IsPresent 检查Optional是否包含值
// IsPresent checks if Optional contains a value
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty 检查Optional是否为空
// IsEmpty checks if Optional is empty
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get 获取Optional中的值，如果为空则panic
// Get retrieves the value in Optional, panics if empty
func (o Optional[T]) Get() T {
	if !o.present {
		panic("Optional is empty")
	}
	return o.value
}

// OrElse 如果Optional为空，则返回默认值
// OrElse returns the default value if Optional is empty
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet 如果Optional为空，则调用函数获取默认值
// OrElseGet calls function to get default value if Optional is empty
func (o Optional[T]) OrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// OrElsePanic 如果Optional为空，则panic并显示消息
// OrElsePanic panics with message if Optional is empty
func (o Optional[T]) OrElsePanic(message string) T {
	if o.present {
		return o.value
	}
	panic(message)
}

// IfPresent 如果Optional包含值，则执行指定函数
// IfPresent executes the specified function if Optional contains a value
func (o Optional[T]) IfPresent(consumer func(T)) {
	if o.present {
		consumer(o.value)
	}
}

// IfPresentOrElse 如果Optional包含值则执行第一个函数，否则执行第二个函数
// IfPresentOrElse executes first function if value present, otherwise second function
func (o Optional[T]) IfPresentOrElse(consumer func(T), emptyAction func()) {
	if o.present {
		consumer(o.value)
	} else {
		emptyAction()
	}
}

// Map 如果Optional包含值，则应用映射函数（链式调用版本）
// Map applies mapping function if Optional contains a value (for method chaining)
func (o Optional[T]) Map(mapper func(T) any) Optional[any] {
	if !o.present {
		return Empty[any]()
	}
	return Of[any](mapper(o.value))
}

// MapTyped 对Optional进行类型安全的映射转换
// MapTyped performs type-safe mapping transformation on Optional
func MapTyped[T, R any](o Optional[T], mapper func(T) R) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return Of(mapper(o.value))
}

// FlatMap 如果Optional包含值，则应用返回Optional的映射函数
// FlatMap applies mapping function that returns Optional if value present
func FlatMap[T, R any](o Optional[T], mapper func(T) Optional[R]) Optional[R] {
	if !o.present {
		return Empty[R]()
	}
	return mapper(o.value)
}

// Filter 如果Optional包含值且满足条件，则返回该Optional，否则返回空Optional
// Filter returns the Optional if it contains value and satisfies condition, otherwise empty
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if !o.present || !predicate(o.value) {
		return Empty[T]()
	}
	return o
}

// String 返回Optional的字符串表示
// String returns string representation of Optional
func (o Optional[T]) String() string {
	if o.present {
		return "Optional[" + fmt.Sprintf("%v", o.value) + "]"
	}
	return "Optional.empty"
}

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

func Must[T any](value T, err error) T {
	return value
}
