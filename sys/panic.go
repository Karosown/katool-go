package sys

import (
	"github.com/karosown/katool-go/xlog"
)

func Panic(err any) {
	xlog.KaToolLoggerWrapper.ApplicationDesc(err).Panic()
}

func Warn(err any) {
	xlog.KaToolLoggerWrapper.ApplicationDesc(err).Warn()
}
