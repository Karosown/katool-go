package stream

import (
	"sort"

	"github.com/karosown/katool/algorithm"
)

type Option[T any] struct {
	opt T
}
type Options[T any] []Option[T]
type Stream[T any, Slice ~[]T] struct {
	options *Options[T]
	source  *Slice
	//parallel bool
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
	size := len(*s.options)
	//if !s.parallel {
	resSource := make([]any, 0)
	for i := 0; i < size; i++ {
		runCall := fn((*s.options)[i].opt)
		resSource = append(resSource, runCall)
	}
	return ToStream(&resSource)
	//}
	//resChan := make(chan any, size)
	//group := &sync.WaitGroup{}
	//group.Add(size)
	//for i := 0; i < size; i++ {
	//	go func(wt *sync.WaitGroup, i int) {
	//		defer wt.Done()
	//		resChan <- fn((*s.options)[i].opt)
	//	}(group, i)
	//}
	//group.Wait()
	//resSource := make([]any, 0)
	//for i := 0; i < size; i++ {
	//	resSource = append(resSource, <-resChan)
	//}
	//return ToStream(&resSource)
}

func (s *Stream[T, Slice]) Distinct(hash func(cnt T) algorithm.HashType) *Stream[T, Slice] {
	//if !s.parallel {
	res := make(Slice, 0)
	size := len(*s.options)
	if size < 1e10+5 {
		sort.SliceStable(*s.options, func(i, j int) bool {
			return hash((*s.options)[i].opt) < hash((*s.options)[j].opt)
		})
		for i := 0; i < size; i++ {
			if i == 0 {
				res = append(res, (*s.options)[i].opt)
			} else if hash((*s.options)[i-1].opt) != hash((*s.options)[i].opt) {
				res = append(res, (*s.options)[i].opt)
			}
		}
	} else {
		//  if large data, use map
		m := make(map[algorithm.HashType]bool, 0)
		for i := 0; i < size; i++ {
			if _, ok := m[hash((*s.options)[i].opt)]; !ok {
				m[hash((*s.options)[i].opt)] = true
				res = append(res, (*s.options)[i].opt)
			}
		}
	}
	return ToStream(&res)
	//}
	//return nil
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
func (s *Stream[T, Slice]) Sort(orderBy func(a, b T) bool) *Stream[T, Slice] {
	sort.SliceStable(*s.options, func(i, j int) bool {
		return orderBy((*s.options)[i].opt, (*s.options)[j].opt)
	})
	return s
}
func (s *Stream[T, Slice]) Collect(call func(data Options[T], sourceData Slice) any) any {
	res := call(*s.options, *s.source)
	return res
}

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
func (s *Stream[T, Slice]) ForEach(fn func(item T)) {
	for i := 0; i < len(*s.options); i++ {
		fn((*s.options)[i].opt)
	}
}

func (s *Stream[T, Slice]) Count() int64 {
	return int64(len(*s.options))
}