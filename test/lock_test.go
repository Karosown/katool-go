package test

import (
	"testing"

	"github.com/karosown/katool/container/stream"
	"github.com/karosown/katool/convert"
	"github.com/karosown/katool/lock"
)

func TestLockSupport(t *testing.T) {
	support := lock.NewLockSupport()
	go func() {
		println("即将进入阻塞，等待异步唤醒")
		support.Park()
		println("唤醒成功")
	}()
	println("准备唤醒")
	support.Unpark()
}

func TestMultLockSupport(t *testing.T) {
	lockss := make([]*lock.LockSupport, 10)
	for i := 0; i < 10; i++ {
		lockss[i] = lock.NewLockSupport()
	}
	for i := 0; i < 10; i++ {
		go func(i int, support *lock.LockSupport) {
			println("即将进入阻塞，等待异步唤醒" + convert.ToString(i))
			support.Park()
			println("唤醒成功" + convert.ToString(i))
		}(i, lockss[i])
	}
	stream.ToStream(&lockss).ToOptionList().ForEach(func(support *lock.LockSupport) {
		println("准备唤醒")
		support.Unpark()
	})

}
