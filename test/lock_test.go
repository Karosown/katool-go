package test

import (
	"testing"

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
