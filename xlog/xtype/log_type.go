package xtype

type LogType int

const (
	ERROR LogType = -1
	INFO  LogType = 0
	DEBUG LogType = 1
	WARN  LogType = 2
)

func (lt LogType) Is(logType LogType) bool {
	return lt == logType
}
func (lt LogType) String() string {
	switch lt {
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case WARN:
		return "WARNING"
	default:
		return "UNKNOWN"
	}
}
