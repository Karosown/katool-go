package stream

type Option[T any] struct {
	opt T
}
type Options[T any] []Option[T]
type Stream[T any, Slice ~[]T] struct {
	options *Options[T]
	source  *Slice
}

func ToStream[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	resOptions := make(Options[T], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[T]{opt: (*source)[i]})
	}
	return &Stream[T, Slice]{
		options: &resOptions,
		source:  source,
	}
}

func (s *Stream[T, Slice]) Map(fn func(i T) any) *Stream[any, []any] {
	resSource := make([]any, 0)
	for i := 0; i < len(*s.options); i++ {
		runCall := fn((*s.options)[i].opt)
		resSource = append(resSource, runCall)
	}
	return ToStream(&resSource)
}
func (s *Stream[T, Slice]) Distinct(validFunc func(cnt, nxt T) bool) *Stream[T, Slice] {
	res := make(Slice, 0)
	for i := 0; i < len(*s.options); i++ {
		if i == 0 {
			res = append(res, (*s.source)[i])
		} else if !validFunc((*s.options)[i-1].opt, (*s.options)[i].opt) {
			res = append(res, (*s.source)[i])
		}
	}
	return ToStream(&res)
}
func (s *Stream[T, Slice]) Reduce(begin any, fn func(cntValue any, nxt T) any) any {
	for i := 0; i < len(*s.options); i++ {
		begin = fn(begin, (*s.options)[i].opt)
	}
	return begin
}

func (s *Stream[T, Slice]) Filter(fn func(i T) bool) *Stream[T, Slice] {
	res := make(Slice, 0)
	for i := 0; i < len(*s.options); i++ {
		if fn((*s.options)[i].opt) {
			res = append(res, (*s.source)[i])
		}
	}
	return ToStream(&res)
}
func (s *Stream[T, Slice]) ToOptionList() Options[T] {
	return *s.options
}
func (s *Stream[T, Slice]) ToList() []T {
	options := s.options
	res := make([]T, 0)
	for i := 0; i < len(*options); i++ {
		res = append(res, (*options)[i].opt)
	}
	return res
}
func (s *Stream[T, Slice]) ToMap(k func(index int, item T) any, v func(i int, item T) any) map[any]any {
	res := make(map[any]any, 0)
	for i := 0; i < len(*s.options); i++ {
		res[k(i, (*s.options)[i].opt)] = v(i, (*s.options)[i].opt)
	}
	return res
}
func (s *Stream[T, Slice]) Collect(call func(data Options[T], sourceData Slice) any) any {
	res := call(*s.options, *s.source)
	return res
}

func (s *Options[T]) ForEach(fn func(item T)) {
	for i := 0; i < len(*s); i++ {
		fn((*s)[i].opt)
	}
}

func (s *Stream[T, Slice]) Count() int64 {
	return int64(len(*s.options))
}
