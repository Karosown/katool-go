package test

import (
	"fmt"
	"testing"

	"github.com/karosown/katool-go/xlog"
)

func Test_logger_wrapper(t *testing.T) {
	fmt.Println(xlog.KaToolLoggerWrapper.Error().ApplicationDesc("测试").Build())
	xlog.KaToolLoggerWrapper.Error().ApplicationDesc("测试")
}
