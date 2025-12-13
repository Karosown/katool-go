package stream

import (
	"fmt"
	lynx "github.com/Tangerg/lynx/pkg/sync"
	"github.com/karosown/katool-go/algorithm"
	"github.com/karosown/katool-go/collect/lists"
	"github.com/karosown/katool-go/container/optional"
)

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
	stream := ToStream(&slice).SetPageSizeGetFunc(source.getPageSize).SetMaxGoroutineNum(source.maxGoroutineNum)
	stream.parallel = source.parallel
	return stream
}

// Of 从给定切片创建流（Java风格API）
// Of creates a stream from the given slice (Java-style API)
func Of[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	return ToStream(source)
}

// newOptionsStream 创建Options流（内部方法）
// newOptionsStream creates Options stream (internal method)
func newOptionsStream[Opt any, Opts Options[Opt]](source *[]Opts, getPageSize func(int) int, goNum int) *Stream[Options[any], []Options[any]] {
	resOptions := make([]Option[Options[any]], 0)
	sourceList := make([]Options[any], 0)
	optionsAdpter := func(opts Options[Opt]) Options[any] {
		res := make(Options[any], 0)
		for i := 0; i < len(opts); i++ {
			convertOption := &Option[any]{opt: any((opts[i].opt))}
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
		options:         (*Options[Options[any]])(&resOptions),
		source:          &sourceList,
		getPageSize:     getPageSize,
		maxGoroutineNum: goNum,
	}
}

// newOptionStream 创建Option流（内部方法）,这个方法不用并行，用于内部转换
// newOptionStream creates Option stream (internal method)
func newOptionStream[Opt any, T Options[Opt]](source *T) *Stream[Option[any], []Option[any]] {
	resOptions := make(Options[Option[any]], 0)
	resSource := make([]Option[any], 0)
	size := len(*source)
	for i := 0; i < size; i++ {
		resOptions = append(resOptions, Option[Option[any]]{
			Option[any]{opt: any(((*source)[i].opt))},
		})
		resSource = append(resSource, Option[any]{opt: any((*source)[i].opt)})
	}
	return &Stream[Option[any], []Option[any]]{
		options:         &resOptions,
		source:          &resSource,
		getPageSize:     getPageSize,
		maxGoroutineNum: algorithm.NumOfTwoMultiply(size),
	}
}

// ToParallelStream 创建并行流
// ToParallelStream creates a parallel stream
func ToParallelStream[T any, Slice ~[]T](source *Slice) *Stream[T, Slice] {
	resOptions := make(Options[T], 0)
	for i := 0; i < len(*source); i++ {
		resOptions = append(resOptions, Option[T]{opt: ((*source)[i])})
	}
	return &Stream[T, Slice]{
		options:     &resOptions,
		source:      source,
		parallel:    true,
		getPageSize: getPageSize,
	}
}

// goRun 并行执行辅助函数
// goRun is a helper function for parallel execution
func goRun[T any](getPageSize func(int) int, maxGoroutineNum int, datas []T, parallel bool, solve func(pos int, automicDatas []T) error) {
	size := len(datas)
	pageSize := optional.IsTrue((getPageSize(size)) == 0, 1, getPageSize(size))
	goNum := optional.IsTrue(maxGoroutineNum == 0, algorithm.NumOfTwoMultiply(size), maxGoroutineNum)
	err := lists.Partition(datas, optional.IsTrue(parallel, pageSize, 1)).ForEach(solve, parallel, lynx.NewLimiter(optional.IsTrue(parallel, goNum, 1)))
	if err != nil {
		fmt.Println(err)
	}
	return
}
