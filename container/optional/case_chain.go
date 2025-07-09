package optional

import (
	"errors"

	"github.com/karosown/katool-go/container/cutil"
)

// OptSwitch 可选择链式操作结构
// OptSwitch is a structure for optional chaining operations
type OptSwitch[T any] struct {
	isRun      bool
	isBreak    bool
	lastStatus bool
	lastResult *T
	lastError  error
}

// Case 根据条件执行操作
// Case executes operations based on a condition
func (t *OptSwitch[T]) Case(flag bool, run ...func() (*T, error)) *OptSwitch[T] {
	if t.isBreak {
		return t
	}
	if !cutil.IsEmpty(run) {
		if flag {
			t.lastStatus = true
			t.lastResult, t.lastError = run[0]()
			t.isRun = true
		} else {
			t.lastStatus = false
			if len(run) >= 2 {
				for i := 1; i < len(run); i++ {
					t.lastResult, t.lastError = run[i]()
					if errors.Is(t.lastError, nil) {
						return t
					}
				}
			}
		}
	}
	return t
}

// CaseFunc 根据函数条件执行操作
// CaseFunc executes operations based on a function condition
func (t *OptSwitch[T]) CaseFunc(flag func() bool, run ...func() (*T, error)) *OptSwitch[T] {
	return t.Case(flag(), run...)
}

// Break 中断链式操作
// Break interrupts the chaining operations
func (t *OptSwitch[T]) Break() *OptSwitch[T] {
	if t.lastStatus {
		t.isBreak = true
	}
	return t
}

// Default 设置默认操作
// Default sets default operations
func (t *OptSwitch[T]) Default(run ...func(res *T, err error) (*T, error)) *OptSwitch[T] {
	if !t.isRun {
		for _, fn := range run {
			t.lastResult, t.lastError = fn(t.lastResult, t.lastError)
		}
	}
	return t
}

// Submit 提交并返回结果
// Submit submits and returns the result
func (t *OptSwitch[T]) Submit() (*T, error) {
	return t.lastResult, t.lastError
}
