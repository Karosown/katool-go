package stream

import "github.com/karosown/katool-go/algorithm"

type AbstarctStream[T any, R any, Slice ~[]T, RSlice ~[]R] interface {
	Sub(begin, end int) *Stream[T, Slice]
	Skip(n int) *Stream[T, Slice]
	Join(source *Slice) *Stream[T, Slice]
	SetPageSizeGetFunc(getter func(size int) int) *Stream[T, Slice]
	SetMaxGoroutineNum(num int) *Stream[T, Slice]
	Map(fn func(i T) R) *Stream[R, RSlice]
	FlatMap(fn func(i T) *Stream[R, RSlice]) *Stream[R, RSlice]
	Distinct() *Stream[T, Slice]
	DistinctBy(hash algorithm.HashComputeFunction[T]) *Stream[T, Slice]
	Reduce(begin any, atomicSolveFunction func(cntValue any, nxt T) any, parallelResultSolve func(sum1, sum2 any) any) any
	Filter(fn func(i T) bool) *Stream[T, Slice]
	ToOptionList() Options[T]
	ToList() Slice
	ToMap(k func(index int, item T) any, v func(i int, item T) any) map[any]any
	GroupBy(groupBy func(item T) any) map[any]Slice
	OrderBy(desc bool, orderBy algorithm.HashComputeFunction[T]) *Stream[T, Slice]
	OrderById(desc bool, orderBy algorithm.IDComputeFunction) *Stream[T, Slice]
	Sort(orderBy func(a, b T) bool) *Stream[T, Slice]
	Collect(call func(data Options[T], sourceData Slice) any) any
	Merge(arrOrStream any) *Stream[T, Slice]
	Intersect(arrOrStream any, validEq ...func(a, b T) bool) *Stream[T, Slice]
	IntersectWith(arrOrStream any, validEq func(a, b T) bool, validHash ...algorithm.HashComputeFunction[T]) *Stream[T, Slice]
	Difference(arrOrStream any, validEq ...func(a, b T) bool) *Stream[T, Slice]
	DifferenceWith(arrOrStream any, validEq func(a, b T) bool, validHash ...algorithm.HashComputeFunction[T]) *Stream[T, Slice]
	Clone() *Stream[T, Slice]
	ForEach(fn func(item T)) *Stream[T, Slice]
	Count() int64
	Parallel() *Stream[T, Slice]
	ParallelWithSetting(pageSizeGetter func(size int) int, maxGoroutineNum int) *Stream[T, Slice]
	UnParallel() *Stream[T, Slice]
}
