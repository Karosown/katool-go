package stream

//TIP 模仿Java的Stream流式处理
// 通常使用ToStream或者Of方法即可（Of是为了保留Java原有的API)，特殊情况构建Any的Stream可以用NewStream方法
// 另外两个Option的Stream方法为内部实现(newOptionsStream 和 newOptionStream)
// 使用泛型的时候注意不要造成泛型循环

// Package stream provides Java-style stream processing capabilities
// Typically use ToStream or Of methods (Of is to preserve Java's original API),
// For special cases building Any Stream, use NewStream method
// The other two Option Stream methods are internal implementations (newOptionsStream and newOptionStream)
// Be careful not to create generic cycles when using generics

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"sync"

	lynx "github.com/Tangerg/lynx/pkg/sync"
	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/collect/lists"
	"github.com/karosown/katool-go/container/cutil"
	"github.com/karosown/katool-go/container/optional"
	"github.com/karosown/katool-go/convert"
	"github.com/karosown/katool-go/sys"
)

// getPageSize 计算分页大小
// getPageSize calculates page size
func getPageSize(size int) int {
	return size >> 2
}

// Stream 流式处理器，模仿Java Stream API
// Stream is a stream processor that mimics Java Stream API
type Stream[T any, Slice ~[]T] struct {
	options  *Options[T]
	source   *Slice
	parallel bool
}

// NewStream 创建any类型的新流
// NewStream creates a new stream of any type
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

// ToStream 从切片创建流
// ToStream creates a stream from a slice
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

// Sub 获取子流（切片）
// Sub gets a sub-stream (slice)
func (s *Stream[T, Slice]) Sub(begin, end int) *Stream[T, Slice] {
	if s.source == nil {
		return &Stream[T, Slice]{
			options: &Options[T]{},
			source:  s.source,
		}
	}
	length := len(*s.source)
	if length == 0 {
		return &Stream[T, Slice]{
			options: &Options[T]{},
			source:  s.source,
		}
	}
	// 处理负数索引
	if begin < 0 {
		begin = length + begin
	}
	if end < 0 {
		end = length + end
	}
	// 边界检查和修正
	if begin < 0 {
		begin = 0
	}
	if end > length {
		end = length
	}
	if begin > end {
		begin = end
	}
	resOptions := make(Options[T], 0, end-begin)
	for i := begin; i < end; i++ {
		resOptions = append(resOptions, Option[T]{opt: (*s.source)[i]})
	}
	return &Stream[T, Slice]{
		options: &resOptions,
		source:  s.source,
	}
}

// Skip 跳过前n个元素
// Skip skips the first n elements
func (s *Stream[T, Slice]) Skip(n int) *Stream[T, Slice] {
	return s.Sub(n, -1)
}

// fromAnySlice 从any切片转换为指定类型切片
// fromAnySlice converts from any slice to specified type slice
func fromAnySlice[T any, Slice ~[]T](source []any) Slice {
	res := make([]T, len(source))
	for i, v := range source {
		res[i], _ = v.(T)
	}
	return res
}

// FromAnySlice 从any切片创建流
// FromAnySlice creates a stream from any slice
func FromAnySlice[T any, Slice ~[]T](source []any) *Stream[T, Slice] {
	slice := fromAnySlice[T, Slice](source)
	return ToStream(&slice)
}

// Cast 类型转换流
// Cast converts stream type
func Cast[T any, Slice []T](source *Stream[any, []any]) *Stream[T, Slice] {
	slice := fromAnySlice[T, Slice](source.ToList())
	return ToStream(&slice)
}

// Of 从给定切片创建流（Java风格API）
// Of creates a stream from the given slice (Java-style API)
func Of[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	return ToStream(source)
}

// newOptionsStream 创建Options流（内部方法）
// newOptionsStream creates Options stream (internal method)
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

// newOptionStream 创建Option流（内部方法）
// newOptionStream creates Option stream (internal method)
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

// ToParallelStream 创建并行流
// ToParallelStream creates a parallel stream
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

// Join 连接两个流
// Join joins two streams
func (s *Stream[T, Slice]) Join(source *Slice) *Stream[T, Slice] {
	list := s.ToList()
	list = append(list, *source...)
	return ToStream(&list)
}

// goRun 并行执行辅助函数
// goRun is a helper function for parallel execution
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

// Map 映射转换元素
// Map transforms elements
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
// FlatMap flattens the stream, requires a function that returns a new Stream
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
// Distinct removes duplicates using default method (default is JSON string comparison)
func (s *Stream[T, Slice]) Distinct() *Stream[T, Slice] {
	return s.DistinctBy(algorithm.HASH_WITH_JSON)
}

// DistinctBy 按照指定方法去重
// DistinctBy removes duplicates using specified method
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
}

// Reduce 求和计算
// Reduce performs aggregation calculation
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

// Filter 过滤元素
// Filter filters elements
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

// ToOptionList 转换为Option列表
// ToOptionList converts to Option list
func (s *Stream[T, Slice]) ToOptionList() Options[T] {
	return *s.options
}

// ToList 转换为列表
// ToList converts to list
func (s *Stream[T, Slice]) ToList() Slice {
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

// ToMap 转换为映射
// ToMap converts to map
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
// GroupBy groups elements by key, ensures GroupBy cannot modify options
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

// OrderBy 按照哈希函数排序
// OrderBy sorts by hash function
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
		itv, _ := mergeSorted[i].(T)
		result = append(result, itv)
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
}

// OrderById 按照ID函数排序
// OrderById sorts by ID function
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
		itv, _ := mergeSorted[i].(T)
		result = append(result, itv)
	}
	//stream := ToStream(&result)
	stream := ToStream(&result)
	return stream
}

// Sort 按照自定义比较函数排序
// Sort sorts by custom comparison function
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
			itv, _ := options[i].opt.(T)
			jtv, _ := options[j].opt.(T)
			return orderBy(itv, jtv)
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
		itv, _ := mergeSorted[i].(T)
		result = append(result, itv)
	}
	stream := ToStream(&result)
	return stream
}

// Collect 收集处理结果
// Collect processes and collects results
func (s *Stream[T, Slice]) Collect(call func(data Options[T], sourceData Slice) any) any {
	res := call(*s.options, *s.source)
	return res
}

// Merge 合并流或数组
// Merge merges streams or arrays
func (s *Stream[T, Slice]) Merge(arrOrStream any) *Stream[T, Slice] {
	var res Slice
	switch arrOrStream.(type) {
	case *Stream[T, Slice]:
		res = append(s.ToList(), arrOrStream.(*Stream[T, Slice]).ToList()...)
	case Slice:
		res = append(s.ToList(), arrOrStream.(Slice)...)
	}
	return ToStream(&res)
}

// Intersect 求交集
// Intersect finds intersection
func (s *Stream[T, Slice]) Intersect(arrOrStream any, validEq ...func(a, b T) bool) *Stream[T, Slice] {
	var temp Slice
	switch arrOrStream.(type) {
	case *Stream[T, Slice]:
		temp = arrOrStream.(*Stream[T, Slice]).ToList()
	case Slice:
		temp = arrOrStream.(Slice)
	}
	return s.Filter(func(i T) bool {
		if cutil.IsEmpty(validEq) {
			return slices.ContainsFunc(temp, func(t T) bool {
				return algorithm.HASH_WITH_JSON(i) == algorithm.HASH_WITH_JSON(t)
			})
		}
		return slices.ContainsFunc(temp, func(t T) bool {
			return validEq[0](i, t)
		})
	})
}

// Difference 求差集
// Difference finds difference
func (s *Stream[T, Slice]) Difference(arrOrStream any, validEq ...func(a, b T) bool) *Stream[T, Slice] {
	var temp Slice
	switch arrOrStream.(type) {
	case *Stream[T, Slice]:
		temp = arrOrStream.(*Stream[T, Slice]).ToList()
	case Slice:
		temp = arrOrStream.(Slice)
	}
	return s.Filter(func(i T) bool {
		if cutil.IsEmpty(validEq) {
			return !slices.ContainsFunc(temp, func(t T) bool {
				return algorithm.HASH_WITH_JSON(i) == algorithm.HASH_WITH_JSON(t)
			})
		}
		return !slices.ContainsFunc(temp, func(t T) bool {
			return validEq[0](i, t)
		})
	})
}

// ForEach 遍历每个元素
// ForEach iterates over each element
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

// Count 计算元素数量
// Count counts the number of elements
func (s *Stream[T, Slice]) Count() int64 {
	return int64(len(*s.options))
}

// Parallel 设置为并行模式
// Parallel sets to parallel mode
func (s *Stream[T, Slice]) Parallel() *Stream[T, Slice] {
	s.parallel = true
	return s
}

// UnParallel 设置为串行模式
// UnParallel sets to serial mode
func (s *Stream[T, Slice]) UnParallel() *Stream[T, Slice] {
	s.parallel = false
	return s
}
