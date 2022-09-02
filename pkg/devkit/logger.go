package devkit

import "github.com/log-rush/distribution-server/domain"

type LogLevel string

const (
	LogLevelDPanic LogLevel = "dpanic"
	LogLevelPanic  LogLevel = "panic"
	LogLevelDebug  LogLevel = "debug"
	LogLevelInfo   LogLevel = "info"
	LogLevelWarn   LogLevel = "warn"
	LogLevelError  LogLevel = "error"
	LogLevelFatal  LogLevel = "fatal"
)

type mockLogger struct {
	prefix  string
	handler func(level LogLevel, log string, args ...interface{})
}

func NewLogger(handler func(level LogLevel, log string, args ...interface{})) domain.Logger {
	return &mockLogger{"", handler}
}

func (l *mockLogger) Debug(args ...interface{}) {
	l.handler(LogLevelDebug, l.prefix+"%s", args...)
}
func (l *mockLogger) Debugf(template string, args ...interface{}) {
	l.handler(LogLevelDebug, l.prefix+template, args...)
}

func (l *mockLogger) Error(args ...interface{}) {
	l.handler(LogLevelError, l.prefix+"%s", args...)
}
func (l *mockLogger) Errorf(template string, args ...interface{}) {
	l.handler(LogLevelError, l.prefix+template, args...)
}

func (l *mockLogger) Fatal(args ...interface{}) {
	l.handler(LogLevelFatal, l.prefix+"%s", args...)
}
func (l *mockLogger) Fatalf(template string, args ...interface{}) {
	l.handler(LogLevelFatal, l.prefix+template, args...)
}

func (l *mockLogger) Info(args ...interface{}) {
	l.handler(LogLevelInfo, l.prefix+"%s", args...)
}
func (l *mockLogger) Infof(template string, args ...interface{}) {
	l.handler(LogLevelInfo, l.prefix+template, args...)
}

func (l *mockLogger) Panic(args ...interface{}) {
	l.handler(LogLevelPanic, l.prefix+"%s", args...)
}
func (l *mockLogger) Panicf(template string, args ...interface{}) {
	l.handler(LogLevelPanic, l.prefix+template, args...)
}

func (l *mockLogger) Warn(args ...interface{}) {
	l.handler(LogLevelWarn, l.prefix+"%s", args...)
}
func (l *mockLogger) Warnf(template string, args ...interface{}) {
	l.handler(LogLevelWarn, l.prefix+template, args...)
}

func (l *mockLogger) Named(name string) domain.Logger {
	return &mockLogger{l.prefix + name + " ", l.handler}
}
