package stream

type MapEnhancer[T any, R any, Slice []T, RSlice []R] struct {
	*Stream[T, Slice]
}

func (e *MapEnhancer[T, R, Slice, RSlice]) Map(fn func(i T) R) *Stream[R, RSlice] {
	return Cast[R, RSlice](e.Stream.Map(func(i T) any {
		return fn(i)
	}))
}

func Eh[T any, R any, Slice []T, RSlice []R](source *Slice) *MapEnhancer[T, R, Slice, RSlice] {
	return &MapEnhancer[T, R, Slice, RSlice]{
		ToStream(source),
	}
}
func Em[T any, R any, Slice []T, RSlice []R](source *Stream[T, Slice]) *MapEnhancer[T, R, Slice, RSlice] {
	return &MapEnhancer[T, R, Slice, RSlice]{
		source,
	}
}
