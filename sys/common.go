package sys

import (
	"runtime"
)

// GetLocalFunctionName 获取调用此函数的函数名称
// GetLocalFunctionName gets the name of the function that calls this function
func GetLocalFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
