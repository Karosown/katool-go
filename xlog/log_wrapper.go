package xlog

import (
	"runtime"

	"github.com/karosown/katool-go/xlog/xtype"
)

// LogWrapper is a wrapper for LogMessage with method chaining support
type LogWrapper struct {
	message        xtype.LogMessage
	functionByFunc func(int) string
	callLayer      int
}

// NewLogWrapper creates a new empty LogWrapper
func NewLogWrapper() *LogWrapper {
	return &LogWrapper{
		message: xtype.LogMessage{
			Type: xtype.INFO, // Default type
		},
	}
}

var (
	KaToolLoggerWrapper = NewLogWrapper().
		Header("KaTool").FunctionByFunc(func(layerCount int) string {
		pc, _, _, _ := runtime.Caller(layerCount)
		return runtime.FuncForPC(pc).Name()
	}).Error()
)

// Header sets the Header field and returns the LogWrapper for chaining
func (lw *LogWrapper) Header(header string) *LogWrapper {
	lw.message.Header = header
	return lw
}

// Function sets the Function field and returns the LogWrapper for chaining
func (lw *LogWrapper) Function(function string) *LogWrapper {
	lw.message.Function = function
	return lw
}

// Function sets the Function field and returns the LogWrapper for chaining
func (lw *LogWrapper) FunctionByFunc(function func(int) string) *LogWrapper {
	lw.callLayer = 1
	lw.functionByFunc = function
	return lw
}

// ApplicationDesc sets the ApplicationDesc field and returns the LogWrapper for chaining
func (lw *LogWrapper) ApplicationDesc(desc any) *LogWrapper {
	lw.message.ApplicationDesc = desc
	return lw
}
func (lw *LogWrapper) Format(format func(l xtype.LogMessage) string) *LogWrapper {
	lw.message.Format = format
	return lw
}

// Type sets the Type field and returns the LogWrapper for chaining
func (lw *LogWrapper) Type(errType xtype.LogType) *LogWrapper {
	lw.message.Type = errType
	return lw
}

// Error sets the Type field to ERROR and returns the LogWrapper for chaining
func (lw *LogWrapper) Error() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.ERROR
	return wrapper
}

func (lw *LogWrapper) Warn() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.WARN
	return wrapper
}

// Info sets the Type field to INFO and returns the LogWrapper for chaining
func (lw *LogWrapper) Info() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.INFO
	return wrapper
}

// Debug sets the Type field to DEBUG and returns the LogWrapper for chaining
func (lw *LogWrapper) Debug() *LogWrapper {
	wrapper := NewLogWrapper()
	wrapper.message = lw.message
	wrapper.message.Type = xtype.DEBUG
	return wrapper
}

// Build returns the underlying LogMessage
func (lw *LogWrapper) Build() xtype.LogMessage {
	lw.callLayer++
	if lw.message.Function == "" && lw.functionByFunc != nil {
		lw.message.Function = lw.functionByFunc(lw.callLayer)
	}
	return lw.message
}

func (lw *LogWrapper) String() string {
	lw.callLayer++
	return lw.Build().String()
}

func (lw *LogWrapper) Panic() *xtype.LogError {
	lw.callLayer++
	return lw.Build().Panic()
}

func (lw *LogWrapper) Throws() *xtype.LogError {
	lw.callLayer++
	return lw.message.Error()
}
