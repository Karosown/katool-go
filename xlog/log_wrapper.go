package xlog

import (
	"runtime"

	"github.com/karosown/katool-go/xlog/xtype"
)

// LogWrapper 支持方法链式调用的日志消息包装器
// LogWrapper is a wrapper for LogMessage with method chaining support
type LogWrapper struct {
	message        xtype.LogMessage
	functionByFunc func(int) string
	callLayer      int
}

// NewLogWrapper 创建新的空日志包装器
// NewLogWrapper creates a new empty LogWrapper
func NewLogWrapper() *LogWrapper {
	return &LogWrapper{
		message: xtype.LogMessage{
			Type: xtype.INFO, // Default type
		},
	}
}

// KaToolLoggerWrapper KaTool默认日志包装器实例
// KaToolLoggerWrapper is the default logger wrapper instance for KaTool
var (
	KaToolLoggerWrapper = NewLogWrapper().
		Header("KaTool").FunctionByFunc(func(layerCount int) string {
		pc, _, _, _ := runtime.Caller(layerCount)
		return runtime.FuncForPC(pc).Name()
	}).Error()
)

// Header 设置头部信息并返回包装器以支持链式调用
// Header sets the Header field and returns the LogWrapper for chaining
func (lw *LogWrapper) Header(header string) *LogWrapper {
	lw.message.Header = header
	return lw
}

// Function 设置函数名并返回包装器以支持链式调用
// Function sets the Function field and returns the LogWrapper for chaining
func (lw *LogWrapper) Function(function string) *LogWrapper {
	lw.message.Function = function
	return lw
}

// FunctionByFunc 通过函数设置函数名并返回包装器以支持链式调用
// FunctionByFunc sets the Function field by function and returns the LogWrapper for chaining
func (lw *LogWrapper) FunctionByFunc(function func(int) string) *LogWrapper {
	lw.callLayer = 1
	lw.functionByFunc = function
	return lw
}

// ApplicationDesc 设置应用描述并返回包装器以支持链式调用
// ApplicationDesc sets the ApplicationDesc field and returns the LogWrapper for chaining
func (lw *LogWrapper) ApplicationDesc(desc any) *LogWrapper {
	lw.message.ApplicationDesc = desc
	return lw
}

// Format 设置格式化函数并返回包装器以支持链式调用
// Format sets the format function and returns the LogWrapper for chaining
func (lw *LogWrapper) Format(format func(l xtype.LogMessage) string) *LogWrapper {
	lw.message.Format = format
	return lw
}

// Type 设置日志类型并返回包装器以支持链式调用
// Type sets the Type field and returns the LogWrapper for chaining
func (lw *LogWrapper) Type(errType xtype.LogType) *LogWrapper {
	lw.message.Type = errType
	return lw
}

// Error 设置日志类型为错误并返回包装器以支持链式调用
// Error sets the Type field to ERROR and returns the LogWrapper for chaining
func (lw *LogWrapper) Error() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.ERROR
	return wrapper
}

// Warn 设置日志类型为警告并返回包装器以支持链式调用
// Warn sets the Type field to WARN and returns the LogWrapper for chaining
func (lw *LogWrapper) Warn() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.WARN
	return wrapper
}

// Info 设置日志类型为信息并返回包装器以支持链式调用
// Info sets the Type field to INFO and returns the LogWrapper for chaining
func (lw *LogWrapper) Info() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.INFO
	return wrapper
}

// Debug 设置日志类型为调试并返回包装器以支持链式调用
// Debug sets the Type field to DEBUG and returns the LogWrapper for chaining
func (lw *LogWrapper) Debug() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.DEBUG
	return wrapper
}

// Build 构建并返回底层的日志消息
// Build returns the underlying LogMessage
func (lw *LogWrapper) Build() xtype.LogMessage {
	lw.callLayer++
	if lw.message.Function == "" && lw.functionByFunc != nil {
		lw.message.Function = lw.functionByFunc(lw.callLayer)
	}
	return lw.message
}

// String 转换为字符串表示
// String converts to string representation
func (lw *LogWrapper) String() string {
	lw.callLayer++
	return lw.Build().String()
}

// Panic 触发panic并返回日志错误
// Panic triggers panic and returns log error
func (lw *LogWrapper) Panic() *xtype.LogError {
	lw.callLayer++
	return lw.Build().Panic()
}

// Throws 抛出错误并返回日志错误
// Throws throws error and returns log error
func (lw *LogWrapper) Throws() *xtype.LogError {
	lw.callLayer++
	return lw.message.Error()
}
