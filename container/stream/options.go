package stream

type Option[T any] struct {
	opt T
}

type Options[T any] []Option[T]

// you can't change the options data with you use the function,if you want to change the options data,use the Stream.ForEach
func (s Options[T]) ForEach(fn func(item T)) Options[T] {
	for i := 0; i < len(s); i++ {
		fn((s)[i].opt)
	}
	return s
}

func (s Options[T]) Stream() *Stream[T, []T] {
	res := make([]T, 0)
	for i := 0; i < len(s); i++ {
		res = append(res, s[i].opt)
	}
	return ToStream(&res)
}
