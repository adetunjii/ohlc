package logging

import "go.uber.org/zap"

type Logger struct {
	instance SugarLogger
}

func NewLogger(sugarLogger *zap.SugaredLogger) *Logger {
	return &Logger{sugarLogger}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.instance.Infof(msg, args)
}

func (l *Logger) Error(msg string, err error) {
	if err != nil {
		l.instance.Errorf(msg, err.Error())
	} else {
		l.instance.Errorf(msg, "error", "unknown error")
	}
}

func (l *Logger) Fatal(msg string, err error) {
	if err != nil {
		l.instance.Fatalf(msg, "error", err.Error())
	} else {
		l.instance.Errorf(msg, "error", "unknown error")
	}
}
