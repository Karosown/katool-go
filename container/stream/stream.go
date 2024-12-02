package stream

import (
	"fmt"
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

func NewOptionsStream[Opt any, Opts Options[Opt]](source *[]Opts) *Stream[Options[any], []Options[any]] {
	resOptions := make([]Option[Options[any]], 0)
	sourceList := make([]Options[any], 0)
	optionsAdpter := func(opts Options[Opt]) Options[any] {
		res := make(Options[any], 0)
		for i := 0; i < len(opts); i++ {
			convertOption := &Option[any]{opt: any(opts[i].opt)}
			res = append(res, *convertOption)
		}
		return res
	}
	for i := 0; i < len(*source); i++ {
		t := Options[Opt]((*source)[i])
		resOptions = append(resOptions, Option[Options[any]]{opt: optionsAdpter(t)})
		sourceList = append(sourceList, optionsAdpter(t))
	}
	return &Stream[Options[any], []Options[any]]{
		options: (*Options[Options[any]])(&resOptions),
		source:  &sourceList,
	}
}

func NewOptionStream[Opt any, T Options[Opt]](source *T) *Stream[Option[any], []Option[any]] {
	resOptions := make(Options[Option[any]], 0)
	resSource := make([]Option[any], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[Option[any]]{
			Option[any]{opt: any((*source)[i].opt)},
		})
		resSource = append(resSource, Option[any]{opt: any((*source)[i].opt)})
	}
	return &Stream[Option[any], []Option[any]]{
		options: (&resOptions),
		source:  &resSource,
	}
}
func ToParallelStream[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	resOptions := make(Options[T], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[T]{opt: (*source)[i]})
	}
	return &Stream[T, Slice]{
		options:  &resOptions,
		source:   source,
		parallel: true,
	}
}
func goRun[T any](datas []T, parallel bool, solve func(pos int, automicDatas []T) error) {
	size := len(datas)
	err := lists.Partition(datas, optional.IsTrue(parallel, (size>>2)^1, 1)).ForEach(solve, parallel, lynx.NewLimiter(optional.IsTrue(parallel, algorithm.NumOfTwoMultiply(size), 1)))
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (s *Stream[T, Slice]) Map(fn func(i T) any) *Stream[any, []any] {
	resSource := make([]any, 0)
	size := len(*s.options)
	resChan := make(chan any, size)
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall
		}
		return nil
	})
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
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall.ToList()
		}
		return nil
	})
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
func (s *Stream[T, Slice]) Reduce(begin any, atomicSolveFunction func(cntValue any, nxt T) any, parallelResultSolve func(sum1, sum2 any) any) any {
	if atomicSolveFunction == nil {
		panic("atomicSolveFunction must not nil")
	}
	if s.parallel && parallelResultSolve == nil {
		panic("parallelResultSolve must not be nil where parallelResult")
	}
	//size := len(*s.options)
	beginType := reflect.TypeOf(begin)
	lock := &sync.Mutex{}
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		if s.parallel {
			currentBegin := reflect.New(beginType).Elem().Interface()
			for i := 0; i < len(options); i++ {
				currentBegin = atomicSolveFunction(currentBegin, options[i].opt)
			}
			lock.Lock()
			defer lock.Unlock()

			if parallelResultSolve != nil {
				begin = parallelResultSolve(begin, currentBegin)
			}
		} else {
			for i := 0; i < len(options); i++ {
				begin = atomicSolveFunction(begin, options[i].opt)
			}
		}
		return nil
	})
	return begin
}

func (s *Stream[T, Slice]) Filter(fn func(i T) bool) *Stream[T, Slice] {
	res := make(Slice, 0)
	size := len(*s.options)
	resChan := make(chan T, size)
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			if fn((options)[i].opt) {
				resChan <- (options)[i].opt
			}
		}
		return nil
	})
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
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			resChan <- (options)[i].opt
		}
		return nil
	})
	for i := 0; i < size; i++ {
		res = append(res, <-resChan)
	}
	return res
}
func (s *Stream[T, Slice]) ToMap(k func(index int, item T) any, v func(i int, item T) any) map[any]any {
	ress := sync.Map{}
	size := len(*s.options)
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			ress.Store(k(pos*optional.IsTrue(s.parallel, size<<4^0x1, 1)+i, (options)[i].opt), v(pos*optional.IsTrue(s.parallel, size<<4^0x1, 1)+i, (options)[i].opt))
		}
		return nil
	})
	res := make(map[any]any)
	ress.Range(func(key, value any) bool {
		res[key] = value
		return true // 继续遍历
	})
	return res
}

func (s *Stream[T, Slice]) GroupBy(groupBy func(item T) any) map[any]Slice {
	res := make(map[any]Slice, 0)
	//size := len(*s.options)
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			key := groupBy((options)[i].opt)
			if _, ok := res[key]; !ok {
				res[key] = make(Slice, 0)
			}
			res[key] = append(res[key], (*s.source)[i])
		}
		return nil
	})
	return res
}

func (s *Stream[T, Slice]) OrderBy(desc bool, orderBy algorithm.HashComputeFunction) *Stream[T, Slice] {
	if !s.parallel {
		sort.SliceStable(*s.options, func(i, j int) bool {
			if desc {
				return orderBy((*s.options)[i].opt) > orderBy((*s.options)[j].opt)
			} else {
				return orderBy((*s.options)[i].opt) < orderBy((*s.options)[j].opt)
			}
		})
		return s
	}

	//size := len(*s.options)
	data := make([]Options[T], 0)
	// opt opt opt opt -> opts opts
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		data = append(data, options)
		return nil
	})
	optionsStream := NewOptionsStream[T, Options[T]](&data)
	optionsStream.parallel = s.parallel
	res := optionsStream.Map(func(options Options[any]) any {
		sort.SliceStable(options, func(i, j int) bool {
			if desc {
				return orderBy((options)[i].opt) > orderBy((options)[j].opt)
			} else {
				return orderBy((options)[i].opt) < orderBy((options)[j].opt)
			}
		})
		return options
	}).Map(func(v any) any {
		i := v.(Options[any])
		ress := NewOptionStream(&i).Map(func(item Option[any]) any {
			return item.opt
		}).ToList()
		return ress
	}).ToList()
	re := make([]any, 0)
	mergeSorted := ToStream(&res).Reduce(re, func(cntValue any, nxt any) any {
		ts := cntValue.([]any)
		nxts := nxt.([]any)
		if len(nxts) == 0 {
			return ts
		}
		lenRe := len(ts)
		lenNxt := len(nxts)
		ress := make([]any, 0)
		l := 0
		r := 0
		for l < lenRe && r < lenNxt {
			current := nxts[r].(T)
			total := ts[l].(T)
			if orderBy(total) > orderBy(current) {
				if desc {
					ress = append(ress, total)
					l++
				} else {
					ress = append(ress, current)
					r++
				}
			} else {
				if desc {
					ress = append(ress, current)
					r++
				} else {
					ress = append(ress, total)
					l++
				}
			}
		}

		if r < lenNxt {
			ress = append(ress, nxts[r:lenNxt]...)
		}
		if l < lenRe {
			ress = append(ress, ts[l:lenRe]...)
		}
		return ress
	}, nil).([]any)
	result := make(Slice, 0)
	for i := 0; i < len(mergeSorted); i++ {
		result = append(result, mergeSorted[i].(T))
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
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
	//size := len(*s.options)
	goRun(*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			fn((options)[i].opt)
		}
		return nil
	})
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
