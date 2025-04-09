package xtype

import (
	"fmt"
)

type LogMessage struct {
	Header          string
	Function        string
	ApplicationDesc string
	Type            LogType
	Format          func(message LogMessage, format string) string
}

// String returns the formatted error message
func (l LogMessage) String() string {
	if l.Format != nil {
		return l.Format(l, l.ApplicationDesc)
	}
	return fmt.Sprintf("%s ==> [%s:%s] == %s",
		l.Header,
		l.Function,
		l.Type,
		l.ApplicationDesc)
}

// Panic returns the formatted error message
func (l LogMessage) Panic() *LogError {
	s := l.String()
	if l.Type.Is(ERROR) {
		panic(s)
	} else {
		fmt.Errorf(s)
	}
	return l.Error()
}

func (l LogMessage) Error() *LogError {
	return NewLogError(l.String(), l.Type)
}
