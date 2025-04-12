package sys

import (
	"runtime"
)

func GetLocalFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
