package stream

//TIP 模仿Java的Stream流式处理
// 通常使用ToStream或者Of方法即可（Of是为了保留Java原有的API)，特殊情况构建Any的Stream可以用NewStream方法
// 另外两个Option的Stream方法为内部实现(newOptionsStream 和 newOptionStream)
// 使用泛型的时候注意不要造成泛型循环

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	lynx "github.com/Tangerg/lynx/pkg/sync"
	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/collect/lists"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/sys"
)

func getPageSize(size int) int {
	return size >> 2
}

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

// Of creates a stream from the given slice.
func Of[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	return ToStream(source)
}

func newOptionsStream[Opt any, Opts Options[Opt]](source *[]Opts) *Stream[Options[any], []Options[any]] {
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

func newOptionStream[Opt any, T Options[Opt]](source *T) *Stream[Option[any], []Option[any]] {
	resOptions := make(Options[Option[any]], 0)
	resSource := make([]Option[any], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[Option[any]]{
			Option[any]{opt: any((*source)[i].opt)},
		})
		resSource = append(resSource, Option[any]{opt: any((*source)[i].opt)})
	}
	return &Stream[Option[any], []Option[any]]{
		options: &resOptions,
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
func (s *Stream[T, Slice]) Join(source *Slice) *Stream[T, []T] {
	list := s.ToList()
	list = append(list, *source...)
	return ToStream(&list)
}
func goRun[T any](datas []T, parallel bool, solve func(pos int, automicDatas []T) error) {
	size := len(datas)
	pageSize := optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size))
	goNum := algorithm.NumOfTwoMultiply(size)
	err := lists.Partition(datas, optional.IsTrue(parallel, pageSize, 1)).ForEach(solve, parallel, lynx.NewLimiter(optional.IsTrue(parallel, goNum, 1)))
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (s *Stream[T, Slice]) Map(fn func(i T) any) *Stream[any, []any] {
	size := len(*s.options)
	resChan := make(chan any, size)
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall
		}
		return nil
	})
	resSource := convert.ChanToArray(resChan)
	return ToStream(&resSource)
}

// FlatMap 扁平化处理，需要放入一个返回新的Stream流的函数
// 参考：https://blog.csdn.net/feinifi/article/details/128980814
func (s *Stream[T, Slice]) FlatMap(fn func(i T) *Stream[any, []any]) *Stream[any, []any] {
	size := len(*s.options)
	resChan := make(chan []any, size)
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			runCall := fn(options[i].opt)
			resChan <- runCall.ToList()
		}
		return nil
	})
	//if !s.parallel {
	resSource := convert.ChanToFlatArray(resChan)
	return ToStream(&resSource)
}

// Distinct 按照默认方法去重(默认是json化之后来进行字符串的比对)
func (s *Stream[T, Slice]) Distinct() *Stream[T, Slice] {
	return s.DistinctBy(algorithm.HASH_WITH_JSON)
}

// DistinctBy 按照指定方法去重
func (s *Stream[T, Slice]) DistinctBy(hash algorithm.HashComputeFunction) *Stream[T, Slice] {
	res := make(Slice, 0)
	size := len(*s.options)
	if size < 1e10+5 {
		options := s.Sort(func(a, b T) bool { return hash(a) < hash(b) }).ToOptionList()
		for i := 0; i < len(options); i++ {
			if i == 0 {
				res = append(res, (options)[i].opt)
			} else if hash((options)[i-1].opt) != hash((options)[i].opt) {
				res = append(res, (options)[i].opt)
			}
		}
	} else {
		//  if large data, use map
		m := make(map[algorithm.HashType]bool)
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

// Reduce 求和计算
func (s *Stream[T, Slice]) Reduce(begin any, atomicSolveFunction func(cntValue any, nxt T) any, parallelResultSolve func(sum1, sum2 any) any) any {
	if atomicSolveFunction == nil {
		sys.Panic("atomicSolveFunction must not nil")
	}
	if s.parallel && parallelResultSolve == nil {
		sys.Panic("parallelResultSolve must not be nil where parallelResult")
	}
	//size := len(*s.options)
	beginType := reflect.TypeOf(begin)
	lock := &sync.Mutex{}
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
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
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			if fn((options)[i].opt) {
				resChan <- (options)[i].opt
			}
		}
		return nil
	})
	chanSize := len(resChan)
	for i := 0; i < chanSize; i++ {
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
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
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
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			index := pos*optional.IsTrue(s.parallel,
				optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size)), 1) + i
			ress.Store(k(index, (options)[i].opt), v(index, (options)[i].opt))
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

// GroupBy 分工单一原则，保证GroupBy无法修改options
func (s *Stream[T, Slice]) GroupBy(groupBy func(item T) any) map[any]Slice {
	res := &sync.Map{}
	size := len(*s.options)
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			key := groupBy((options)[i].opt)
			if _, ok := res.Load(key); !ok {
				res.Store(key, make(Slice, 0))
			}
			value, ok := res.Load(key)
			if ok {
				index := pos*optional.IsTrue(s.parallel,
					optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size)), 1) + i
				res.Store(key, append(value.(Slice), (*s.source)[index]))
			}

		}
		return nil
	})
	result := make(map[any]Slice)
	res.Range(func(key, value any) bool {
		result[key] = value.(Slice)
		return true // 继续遍历
	})
	return result
}

func (s *Stream[T, Slice]) OrderBy(desc bool, orderBy algorithm.HashComputeFunction) *Stream[T, Slice] {
	if !s.parallel {
		sort.SliceStable(*s.options, func(i, j int) bool {
			a := orderBy((*s.options)[i].opt)
			b := orderBy((*s.options)[j].opt)
			if desc {
				return a > b
			} else {
				return a < b
			}
		})
		return s
	}

	size := len(*s.options)
	data := make([]Options[T], 0, optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size)))
	// opt opt opt opt -> opts opts
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		//println(pos, options)
		data = append(data, options)
		return nil
	})

	optionsStream := newOptionsStream[T, Options[T]](&data)

	optionsStream.parallel = s.parallel
	sortedMap := optionsStream.Map(func(options Options[any]) any {
		sort.SliceStable(options, func(i, j int) bool {
			a := orderBy(options[i].opt)
			b := orderBy(options[j].opt)
			if desc {
				return a > b
			} else {
				return a < b
			}
		})
		return options
	})
	sortedMap.parallel = s.parallel
	res := sortedMap.Map(func(v any) any {
		i := v.(Options[any])
		ress := newOptionStream(&i).Map(func(item Option[any]) any {
			return item.opt
		}).ToList()
		return ress
	}).ToList()
	re := make([]any, 0)
	toStream := ToStream(&res)
	toStream.parallel = s.parallel
	mergeSorted := toStream.Reduce(re, algorithm.MergeSortedArrayWithPrimaryData[T](desc, orderBy), algorithm.MergeSortedArrayWithPrimaryData[T](desc, orderBy)).([]any)
	result := make(Slice, 0)
	for i := 0; i < len(mergeSorted); i++ {
		result = append(result, mergeSorted[i].(T))
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
}

func (s *Stream[T, Slice]) OrderById(desc bool, orderBy algorithm.IDComputeFunction) *Stream[T, Slice] {
	if !s.parallel {
		sort.SliceStable(*s.options, func(i, j int) bool {
			a := orderBy((*s.options)[i].opt)
			b := orderBy((*s.options)[j].opt)
			if desc {
				return a > b
			} else {
				return a < b
			}
		})
		return s
	}

	size := len(*s.options)
	data := make([]Options[T], 0, optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size)))
	// opt opt opt opt -> opts opts
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		//println(pos, options)
		data = append(data, options)
		return nil
	})
	optionsStream := newOptionsStream[T, Options[T]](&data)
	optionsStream.parallel = s.parallel
	sortedMap := optionsStream.Map(func(options Options[any]) any {
		sort.SliceStable(options, func(i, j int) bool {
			a := orderBy(options[i].opt)
			b := orderBy(options[j].opt)
			if desc {
				return a > b
			} else {
				return a < b
			}
		})
		return options
	})
	sortedMap.parallel = s.parallel
	res := sortedMap.Map(func(v any) any {
		i := v.(Options[any])
		ress := newOptionStream(&i).Map(func(item Option[any]) any {
			return item.opt
		}).ToList()
		return ress
	}).ToList()
	re := make([]any, 0)
	toStream := ToStream(&res)
	toStream.parallel = s.parallel
	mergeSorted := toStream.Reduce(re, algorithm.MergeSortedArrayWithPrimaryId[T](desc, orderBy), algorithm.MergeSortedArrayWithPrimaryId[T](desc, orderBy)).([]any)
	result := make(Slice, 0)
	for i := 0; i < len(mergeSorted); i++ {
		result = append(result, mergeSorted[i].(T))
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
}

func (s *Stream[T, Slice]) Sort(orderBy func(a, b T) bool) *Stream[T, Slice] {

	if !s.parallel {
		sort.SliceStable(*s.options, func(i, j int) bool {
			return orderBy((*s.options)[i].opt, (*s.options)[j].opt)
		})
		return s
	}

	size := len(*s.options)
	data := make([]Options[T], 0, optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size)))
	// opt opt opt opt -> opts opts
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		//println(pos, options)
		data = append(data, options)
		return nil
	})
	optionsStream := newOptionsStream[T, Options[T]](&data)
	optionsStream.parallel = s.parallel
	sortedMap := optionsStream.Map(func(options Options[any]) any {
		sort.SliceStable(options, func(i, j int) bool {
			return orderBy(options[i].opt.(T), options[j].opt.(T))
		})
		return options
	})
	sortedMap.parallel = s.parallel
	res := sortedMap.Map(func(v any) any {
		i := v.(Options[any])
		ress := newOptionStream(&i).Map(func(item Option[any]) any {
			return item.opt
		}).ToList()
		return ress
	}).ToList()
	re := make([]any, 0)
	toStream := ToStream(&res)
	toStream.parallel = s.parallel
	mergeSorted := toStream.Reduce(re, algorithm.MergeSortedArrayWithEntity[T](orderBy), algorithm.MergeSortedArrayWithEntity[T](orderBy)).([]any)
	result := make(Slice, 0)
	for i := 0; i < len(mergeSorted); i++ {
		result = append(result, mergeSorted[i].(T))
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
}
func (s *Stream[T, Slice]) Collect(call func(data Options[T], sourceData Slice) any) any {
	res := call(*s.options, *s.source)
	return res
}

func (s *Stream[T, Slice]) ForEach(fn func(item T)) *Stream[T, Slice] {
	//size := len(*s.options)
	goRun[Option[T]](*s.options, s.parallel, func(pos int, options []Option[T]) error {
		for i := 0; i < len(options); i++ {
			fn((options)[i].opt)
		}
		return nil
	})
	return s
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
