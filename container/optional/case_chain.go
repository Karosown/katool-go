package optional

import (
	"errors"

	"github.com/karosown/katool-go/container/cutil"
)

type OptSwitch[T any] struct {
	isRun      bool
	isBreak    bool
	lastStatus bool
	lastResult *T
	lastError  error
}

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

func (t *OptSwitch[T]) CaseFunc(flag func() bool, run ...func() (*T, error)) *OptSwitch[T] {
	return t.Case(flag(), run...)
}
func (t *OptSwitch[T]) Break() *OptSwitch[T] {
	if t.lastStatus {
		t.isBreak = true
	}
	return t
}
func (t *OptSwitch[T]) Default(run ...func(res *T, err error) (*T, error)) *OptSwitch[T] {
	if !t.isRun {
		for _, fn := range run {
			t.lastResult, t.lastError = fn(t.lastResult, t.lastError)
		}
	}
	return t
}

func (t *OptSwitch[T]) Submit() (*T, error) {
	return t.lastResult, t.lastError
}
