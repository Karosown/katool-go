package xheap

import (
	"cmp"
	"container/heap"

	"github.com/karosown/katool-go/sys"
)

// Heap 是可同时作为小根堆或大根堆使用的泛型优先级队列。
// 默认比较器为小根堆 (less = a < b)。
type Heap[T any] struct {
	data []T
	less func(a, b T) bool // 返回 true 表示 a 的优先级“更高”（在堆顶）
}

// --- 公共构造函数 ---------------------------------------------------------

// NewHeap 允许自定义比较器；比较器必须是严格弱序（同样的要求见 container/heap 文档）
func NewHeap[T any](less func(a, b T) bool, init ...T) *Heap[T] {
	if less == nil {
		sys.Panic("xtype.Heap: less function must not be nil")
	}
	h := &Heap[T]{data: append([]T(nil), init...), less: less}
	heap.Init(h)
	return h
}

// NewMinHeap[T cmp.Ordered] 生成小根堆 (最小值优先)
func NewMinHeap[T cmp.Ordered](init ...T) *Heap[T] {
	return NewHeap(func(a, b T) bool { return a < b }, init...)
}

// NewMaxHeap[T cmp.Ordered] 生成大根堆 (最大值优先)
func NewMaxHeap[T cmp.Ordered](init ...T) *Heap[T] {
	return NewHeap(func(a, b T) bool { return a > b }, init...)
}

// --- heap.Interface 实现 ---------------------------------------------------

func (h *Heap[T]) Len() int           { return len(h.data) }
func (h *Heap[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *Heap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *Heap[T]) Push(x any)         { h.data = append(h.data, x.(T)) }
func (h *Heap[T]) Pop() any           { n := len(h.data) - 1; v := h.data[n]; h.data = h.data[:n]; return v }

// --- 更友好的方法 -----------------------------------------------------------

// PushVal 插入元素
func (h *Heap[T]) PushVal(v T) { heap.Push(h, v) }

// PopVal 弹出堆顶元素；堆空时 panic，与 container/heap 一致
func (h *Heap[T]) PopVal() T { return heap.Pop(h).(T) }

// Peek 查看堆顶，但不移除。若为空，第二个返回值为 false。
func (h *Heap[T]) Peek() (v T, ok bool) {
	if len(h.data) == 0 {
		return v, false
	}
	return h.data[0], true
}

// Fix 在外部直接修改某个索引的值后，重新调整堆序。
func (h *Heap[T]) Fix(i int) { heap.Fix(h, i) }

// RemoveIdx 删除并返回索引 i 处的元素。
func (h *Heap[T]) RemoveIdx(i int) T { return heap.Remove(h, i).(T) }

// Data 返回内部切片的引用，主要用于调试或 Fix 场景。
// ***注意***：这是同一块底层内存，外部改动会直接影响堆。
// 如果希望只读，可自行改成返回拷贝。
func (h *Heap[T]) Data() []T { return h.data }
