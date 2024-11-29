package stream

import (
	"reflect"
	"sort"
	"sync"

	lynx "github.com/Tangerg/lynx/pkg/sync"
	"github.com/karosown/katool/algorithm"
	"github.com/karosown/katool/collect/lists"
	"github.com/karosown/katool/container/optional"
)

type Stream[T any, Slice ~[]T] struct {
	options  *Options[T]
	source   *Slice
	parallel bool
}

func NewStream(source *[]any) *Stream[any, []any] {
	resOptions := make(Options[any], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[any]{opt: (*source)[i]})
	}
	return &Stream[any, []any]{
		options: &resOptions,
		source:  source,
	}
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
	size := len(*s.options)
	resChan := make(chan any, size)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	//if !s.parallel {
	for i := 0; i < size; i++ {
		resSource = append(resSource, <-resChan)
	}
	return ToStream(&resSource)
}
func (s *Stream[T, Slice]) FlatMap(fn func(i T) *Stream[any, []any]) *Stream[any, []any] {
	size := len(*s.options)
	resSource := make([]any, 0)
	resChan := make(chan []any, size)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall.ToList()
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	//if !s.parallel {
	for i := 0; i < size; i++ {
		resSource = append(resSource, <-resChan...)
	}
	return ToStream(&resSource)
}
func (s *Stream[T, Slice]) Distinct(hash algorithm.HashComputeFunction) *Stream[T, Slice] {
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
func (s *Stream[T, Slice]) Reduce(begin any, singleSolve func(cntValue any, nxt T) any, parallelResultSolve func(sum1, sum2 any) any) any {
	size := len(*s.options)
	beginType := reflect.TypeOf(begin)
	lock := &sync.Mutex{}
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		currentBegin := reflect.New(beginType).Elem().Interface()
		for i := 0; i < len(options); i++ {
			currentBegin = singleSolve(currentBegin, options[i].opt)
		}
		lock.Lock()
		defer lock.Unlock()
		begin = parallelResultSolve(begin, currentBegin)
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	return begin
}

func (s *Stream[T, Slice]) Filter(fn func(i T) bool) *Stream[T, Slice] {
	res := make(Slice, 0)
	size := len(*s.options)
	resChan := make(chan T, size)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			if fn((options)[i].opt) {
				resChan <- (options)[i].opt
			}
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	for i := 0; i < len(resChan); i++ {
		res = append(res, <-resChan)
	}
	return ToStream(&res)
}
func (s *Stream[T, Slice]) ToOptionList() Options[T] {
	return *s.options
}
func (s *Stream[T, Slice]) ToList() []T {
	res := make([]T, 0)
	size := len(*s.options)
	resChan := make(chan T, size)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			resChan <- (options)[i].opt
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	for i := 0; i < size; i++ {
		res = append(res, <-resChan)
	}
	return res
}
func (s *Stream[T, Slice]) ToMap(k func(index int, item T) any, v func(i int, item T) any) map[any]any {
	ress := sync.Map{}
	size := len(*s.options)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			ress.Store(k(pos*optional.IsTrue(s.parallel, size<<4^0x1, 1)+i, (options)[i].opt), v(pos*optional.IsTrue(s.parallel, size<<4^0x1, 1)+i, (options)[i].opt))
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	res := make(map[interface{}]interface{})
	ress.Range(func(key, value interface{}) bool {
		res[key] = value
		return true // 继续遍历
	})
	return res
}

func (s *Stream[T, Slice]) GroupBy(groupBy func(item T) any) map[any]Slice {
	res := make(map[any]Slice, 0)
	size := len(*s.options)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			key := groupBy((options)[i].opt)
			if _, ok := res[key]; !ok {
				res[key] = make(Slice, 0)
			}
			res[key] = append(res[key], (*s.source)[i])
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
	return res
}

func (s *Stream[T, Slice]) OrderBy(desc bool, orderBy algorithm.HashComputeFunction) *Stream[T, Slice] {
	sort.SliceStable(*s.options, func(i, j int) bool {
		if desc {
			return orderBy((*s.options)[i].opt) > orderBy((*s.options)[j].opt)
		} else {
			return orderBy((*s.options)[i].opt) < orderBy((*s.options)[j].opt)
		}
	})
	return s
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

func (s *Stream[T, Slice]) ForEach(fn func(item T)) {
	size := len(*s.options)
	lists.Partition(*s.options, optional.IsTrue(s.parallel, size<<4^0x1, 1)).ForEach(func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			fn((options)[i].opt)
		}
		return nil
	}, s.parallel, lynx.NewLimiter(optional.IsTrue(s.parallel, algorithm.NumOfTwoMultiply(size), 1)))
}

func (s *Stream[T, Slice]) Count() int64 {
	return int64(len(*s.options))
}

func (s *Stream[T, Slice]) Parallel() *Stream[T, Slice] {
	s.parallel = true
	return s
}

func (s *Stream[T, Slice]) UnParallel() *Stream[T, Slice] {
	s.parallel = false
	return s
}
