package test

import (
	"testing"

	"katool/container/stream"
	"katool/convert"
	"katool/lock"
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
			println("即将进入阻塞，等待异步唤醒" + convert.ConvertToString(i))
			support.Park()
			println("唤醒成功" + convert.ConvertToString(i))
		}(i, lockss[i])
	}
	stream.ToStream(&lockss).ToOptionList().ForEach(func(support *lock.LockSupport) {
		println("准备唤醒")
		support.Unpark()
	})
}
