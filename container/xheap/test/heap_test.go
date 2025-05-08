package test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/karosown/katool-go/container/xheap"
)

// ------- 小根堆 ------------------------------------------------------------

func TestMinHeap_Int(t *testing.T) {
	h := xheap.NewMinHeap[int]()
	input := []int{5, 3, 9, 1, 4}
	for _, v := range input {
		h.PushVal(v)
	}

	// 弹出顺序应为升序
	want := []int{1, 3, 4, 5, 9}
	var got []int
	for h.Len() > 0 {
		got = append(got, h.PopVal())
	}
	if !slices.Equal(want, got) {
		t.Errorf("min‑heap order wrong: want %v, got %v", want, got)
	}
}

// ------- 大根堆 ------------------------------------------------------------

func TestMaxHeap_Int(t *testing.T) {
	h := xheap.NewMaxHeap[int]()
	input := []int{7, 2, 10, 6}
	for _, v := range input {
		h.PushVal(v)
	}

	want := []int{10, 7, 6, 2}
	for _, w := range want {
		if got := h.PopVal(); got != w {
			t.Fatalf("max‑heap popped %d, want %d", got, w)
		}
	}
}

// ------- 自定义比较器 ------------------------------------------------------

func TestCustomHeap_StringLen(t *testing.T) {
	// “字符串越短优先级越高”
	h := xheap.NewHeap(func(a, b string) bool { return len(a) < len(b) },
		"golang", "is", "fun", "!")
	order := []string{"!", "is", "fun", "golang"}
	for _, w := range order {
		if got := h.PopVal(); got != w {
			t.Fatalf("custom heap popped %q, want %q", got, w)
		}
	}
}

// ------- Fix 与 RemoveIdx --------------------------------------------------

func TestFixAndRemove(t *testing.T) {
	h := xheap.NewMinHeap[int](9, 8, 7, 6) // 初始建堆
	// 手动把索引 2 的元素改成更小的 0，并调用 Fix
	h.Data()[2] = 0
	h.Fix(2)

	if top, _ := h.Peek(); top != 0 {
		t.Fatalf("before Fix, top=%d, want 0", top)
	}

	// 删除索引 1（任意位置）并检测数量
	before := h.Len()
	removed := h.RemoveIdx(1)
	if h.Len() != before-1 {
		t.Fatalf("RemoveIdx length wrong, got %d, want %d", h.Len(), before-1)
	}
	if removed == 0 {
		t.Fatalf("RemoveIdx removed unexpected element 0 (should be non‑zero)")
	}
}

// ------- Benchmark ---------------------------------------------------------

func BenchmarkHeapPushPop(b *testing.B) {
	const N = 1_000
	r := rand.New(rand.NewSource(42))

	b.Run("MinHeap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := xheap.NewMinHeap[int]()
			for j := 0; j < N; j++ {
				h.PushVal(r.Int())
			}
			for h.Len() > 0 {
				_ = h.PopVal()
			}
		}
	})
}
