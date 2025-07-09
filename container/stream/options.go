package stream

// Option 可选项包装器
// Option is an optional value wrapper
type Option[T any] struct {
	opt T
}

// Options 可选项集合
// Options is a collection of optional values
type Options[T any] []Option[T]

// ForEach you can't change the options data with you use the function,if you want to change the options data,use the Stream.ForEach
// ForEach 遍历每个选项（不能修改数据，如需修改请使用Stream.ForEach）
// ForEach iterates over each option (cannot modify data, use Stream.ForEach if modification needed)
func (s Options[T]) ForEach(fn func(item T)) Options[T] {
	for i := 0; i < len(s); i++ {
		fn((s)[i].opt)
	}
	return s
}

// Stream 转换为流
// Stream converts to stream
func (s Options[T]) Stream() *Stream[T, []T] {
	res := make([]T, 0)
	for i := 0; i < len(s); i++ {
		res = append(res, s[i].opt)
	}
	return ToStream(&res)
}
