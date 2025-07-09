package xtype

// LogType 日志类型定义
// LogType defines the log type
type LogType int

// 日志级别常量定义
// Log level constants definition
const (
	ERROR LogType = -1 // 错误级别 / Error level
	INFO  LogType = 0  // 信息级别 / Info level
	DEBUG LogType = 1  // 调试级别 / Debug level
	WARN  LogType = 2  // 警告级别 / Warning level
)

// Is 检查日志类型是否匹配
// Is checks if the log type matches
func (lt LogType) Is(logType LogType) bool {
	return lt == logType
}

// String 转换为字符串表示
// String converts to string representation
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
