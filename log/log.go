package log

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Warnln(arg ...interface{})
	Infoln(arg ...interface{})
	Errorln(arg ...interface{})
	Warn(arg ...interface{})
	Info(arg ...interface{})
	Error(arg ...interface{})
}
