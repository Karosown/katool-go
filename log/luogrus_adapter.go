package log

import (
	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
}

func (l LogrusAdapter) Infof(s string, i ...any) {
	logrus.Infof(s, i...)
}

func (l LogrusAdapter) Errorf(s string, i ...any) {
	logrus.Errorf(s, i...)
}

func (l LogrusAdapter) Warnf(s string, i ...any) {
	logrus.Warnf(s, i...)
}

func (l LogrusAdapter) Warnln(arg ...any) {
	logrus.Warnln(arg...)
}

func (l LogrusAdapter) Infoln(arg ...any) {
	logrus.Infoln(arg...)
}

func (l LogrusAdapter) Errorln(arg ...any) {
	logrus.Errorln(arg...)
}

func (l LogrusAdapter) Warn(arg ...any) {
	logrus.Warn(arg...)
}

func (l LogrusAdapter) Info(arg ...any) {
	logrus.Info(arg...)
}

func (l LogrusAdapter) Error(arg ...any) {
	logrus.Error(arg...)
}
