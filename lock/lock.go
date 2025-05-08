package lock

import (
	"sync"

	"github.com/karosown/katool-go/container/xmap"
)

type LockSupport struct {
	wt    chan bool
	state chan bool
}

func NewLockSupport() *LockSupport {
	return &LockSupport{
		wt:    make(chan bool),
		state: make(chan bool),
	}
}

func (l *LockSupport) Park() bool {
	if len(l.state) >= 1 {
		return false
	}
	l.state <- true
	return <-l.wt
}

func (l *LockSupport) Unpark() (err error) {
	defer func() error {
		e := recover()
		if e != nil {
			err = e.(error)
		}
		if err != nil {
			return err
		}
		return nil
	}()
	<-l.state
	if len(l.wt) >= 1 {
		return nil
	}
	l.wt <- true
	return nil
}

func Synchronized(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}
func SynchronizedErr(locker sync.Locker, f func() error) error {
	locker.Lock()
	defer locker.Unlock()
	return f()
}
func SynchronizedT[T any](locker sync.Locker, f func() T) T {
	locker.Lock()
	defer locker.Unlock()
	return f()
}
func SynchronizedTErr[T any](locker sync.Locker, f func() (T, error)) (T, error) {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

type LockMap xmap.SafeMap[string, sync.Locker]
