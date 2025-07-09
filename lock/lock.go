package lock

import (
	"sync"

	"github.com/karosown/katool-go/container/xmap"
)

// LockSupport 类似Java的LockSupport，提供park/unpark机制
// LockSupport provides park/unpark mechanism similar to Java's LockSupport
type LockSupport struct {
	wt    chan bool
	state chan bool
}

// NewLockSupport 创建新的LockSupport实例
// NewLockSupport creates a new LockSupport instance
func NewLockSupport() *LockSupport {
	return &LockSupport{
		wt:    make(chan bool),
		state: make(chan bool),
	}
}

// Park 阻塞当前协程直到被unpark唤醒
// Park blocks the current goroutine until it's unparked
func (l *LockSupport) Park() bool {
	if len(l.state) >= 1 {
		return false
	}
	l.state <- true
	return <-l.wt
}

// Unpark 唤醒被park阻塞的协程
// Unpark wakes up the goroutine that was parked
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

// Synchronized 在锁保护下执行函数
// Synchronized executes a function under lock protection
func Synchronized(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}

// SynchronizedErr 在锁保护下执行带错误返回的函数
// SynchronizedErr executes a function with error return under lock protection
func SynchronizedErr(locker sync.Locker, f func() error) error {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

// SynchronizedT 在锁保护下执行带泛型返回值的函数
// SynchronizedT executes a function with generic return value under lock protection
func SynchronizedT[T any](locker sync.Locker, f func() T) T {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

// SynchronizedTErr 在锁保护下执行带泛型返回值和错误的函数
// SynchronizedTErr executes a function with generic return value and error under lock protection
func SynchronizedTErr[T any](locker sync.Locker, f func() (T, error)) (T, error) {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

// LockMap 线程安全的锁映射
// LockMap is a thread-safe map for locks
type LockMap = xmap.SafeMap[string, sync.Locker]
