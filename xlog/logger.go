package xlog

// Logger 日志记录器接口
// Logger is a logging interface
type Logger interface {
	Infof(string, ...any)
	Errorf(string, ...any)
	Warnf(string, ...any)
	Warnln(arg ...any)
	Infoln(arg ...any)
	Errorln(arg ...any)
	Warn(arg ...any)
	Info(arg ...any)
	Error(arg ...any)
}
