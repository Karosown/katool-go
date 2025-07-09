package sys

import (
	"github.com/karosown/katool-go/xlog"
)

// Panic 记录错误信息并触发panic
// Panic logs an error message and triggers a panic
func Panic(err any) {
	xlog.KaToolLoggerWrapper.ApplicationDesc(err).Panic()
}

// Warn 记录警告信息
// Warn logs a warning message
func Warn(err any) {
	xlog.KaToolLoggerWrapper.ApplicationDesc(err).Warn()
}
