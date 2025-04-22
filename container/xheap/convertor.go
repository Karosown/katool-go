package xheap

import (
	"cmp"
	"container/heap"
)

/*
   ┌──────────────┐     ToSafe()     ┌─────────────────┐
   │  *Heap[T]    │ ───────────────► │ *SafeHeap[T]    │
   └──────────────┘                  └─────────────────┘
            ▲           ToUnsafe()            │
            └─────────────────────────────────┘

最/大根堆互转（仅 cmp.Ordered 可用）
   NewMaxFromMin / NewMinFromMax
*/

// ---------------- Heap <=> SafeHeap ------------------------------------

// ToSafe 生成一个线程安全包装；内部与原 Heap **共享** 底层数据。
// 若想得到独立副本可用 Clone() 再 ToSafe()。
func (h *Heap[T]) ToSafe() *SafeHeap[T] {
	return &SafeHeap[T]{raw: h}
}

// ToUnsafe 将 SafeHeap 的底层 *Heap 暴露出来。
// ⚠️ 仅在所有并发操作已停止、且你确信自己会外部加锁时使用。
func (sh *SafeHeap[T]) ToUnsafe() *Heap[T] {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return sh.raw
}

// ---------------- 复制 --------------------------------------------------

// Clone 深拷贝一个全新的 Heap（比较器保持不变）
func (h *Heap[T]) Clone() *Heap[T] {
	dup := append([]T(nil), h.data...) // copy
	return NewHeap(h.less, dup...)
}

// ---------------- 最小堆 <→ 最大堆 --------------------------------------

// NewMaxFromMin 把“最小堆”复制为“最大堆”。
// 仅限元素满足 cmp.Ordered；否则请用 CloneWithLess。
func NewMaxFromMin[T cmp.Ordered](src *Heap[T]) *Heap[T] {
	cp := src.Clone()
	cp.less = func(a, b T) bool { return a > b } // 反转
	heapify(cp)
	return cp
}

// NewMinFromMax 相反方向
func NewMinFromMax[T cmp.Ordered](src *Heap[T]) *Heap[T] {
	cp := src.Clone()
	cp.less = func(a, b T) bool { return a < b }
	heapify(cp)
	return cp
}

// CloneWithLess：复制并替换比较器（适用于完全自定义顺序）
func CloneWithLess[T any](src *Heap[T], less func(a, b T) bool) *Heap[T] {
	cp := src.Clone()
	cp.less = less
	heapify(cp)
	return cp
}

// 内部工具：重新建堆
func heapify[T any](h *Heap[T]) { heap.Init(h) }
