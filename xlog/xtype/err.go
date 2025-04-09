package xtype

import (
	"errors"
)

type LogError struct {
	error
	LogType
}

func NewLogError(err string, logType LogType) *LogError {
	return &LogError{errors.New(err), logType}
}
