package xheap

import (
	"cmp"
	"sync"

	"github.com/karosown/katool-go/lock"
)

// SafeHeap 与 Heap API 基本一致，但内部用 RWMutex 保证并发安全。
type SafeHeap[T any] struct {
	mu  sync.RWMutex
	raw *Heap[T]
}

// ---------- 构造函数 ------------------------------------------------------

// NewSafeHeap 创建线程安全堆，需自定义 less。
func NewSafeHeap[T any](less func(a, b T) bool, init ...T) *SafeHeap[T] {
	return &SafeHeap[T]{raw: NewHeap(less, init...)}
}

// NewSafeMinHeap / NewSafeMaxHeap
func NewSafeMinHeap[T cmp.Ordered](init ...T) *SafeHeap[T] {
	return &SafeHeap[T]{raw: NewMinHeap(init...)}
}
func NewSafeMaxHeap[T cmp.Ordered](init ...T) *SafeHeap[T] {
	return &SafeHeap[T]{raw: NewMaxHeap(init...)}
}

// ---------- 只读操作：RLock ----------------------------------------------

func (h *SafeHeap[T]) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.raw.Len()
}

func (h *SafeHeap[T]) Peek() (v T, ok bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.raw.Peek()
}

// Data 返回底层切片的拷贝（避免外部篡改）
func (h *SafeHeap[T]) Data() []T {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return append([]T(nil), h.raw.Data()...) // 拷贝以保持只读语义
}

// ---------- 写操作：Lock ---------------------------------------------------

func (h *SafeHeap[T]) PushVal(v T) {
	lock.Synchronized(&h.mu, func() {
		h.raw.PushVal(v)
	})
}

func (h *SafeHeap[T]) PopVal() T {
	var idx T
	lock.Synchronized(&h.mu, func() {
		idx = h.raw.PopVal()
	})
	return idx
}

func (h *SafeHeap[T]) Fix(i int) {
	lock.Synchronized(&h.mu, func() {
		h.raw.Fix(i)
	})
}

func (h *SafeHeap[T]) RemoveIdx(i int) T {
	var idx T
	lock.Synchronized(&h.mu, func() {
		idx = h.raw.RemoveIdx(i)
	})
	return idx
}
