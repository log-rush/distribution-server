package distributionServer

import (
	"github.com/log-rush/distribution-server/domain"
)

type MultiLogger struct {
	Loggers *[]domain.Logger
}

func NewMultiLogger(loggers ...domain.Logger) domain.Logger {
	return &MultiLogger{Loggers: &loggers}
}

func (l *MultiLogger) Debug(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Debug(args...)
	}
}
func (l *MultiLogger) Debugf(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Debugf(template, args...)
	}
}

func (l *MultiLogger) Error(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Error(args...)
	}
}
func (l *MultiLogger) Errorf(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Errorf(template, args...)
	}
}

func (l *MultiLogger) Fatal(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Fatal(args...)
	}
}
func (l *MultiLogger) Fatalf(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Fatalf(template, args...)
	}
}

func (l *MultiLogger) Info(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Info(args...)
	}
}
func (l *MultiLogger) Infof(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Infof(template, args...)
	}
}

func (l *MultiLogger) Panic(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Panic(args...)
	}
}
func (l *MultiLogger) Panicf(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Panicf(template, args...)
	}
}

func (l *MultiLogger) Warn(args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Warn(args...)
	}
}
func (l *MultiLogger) Warnf(template string, args ...interface{}) {
	for _, logger := range *l.Loggers {
		logger.Warnf(template, args...)
	}
}

func (l *MultiLogger) Named(name string) domain.Logger {
	newLoggers := []domain.Logger{}
	for _, logger := range *l.Loggers {
		newLoggers = append(newLoggers, logger.Named(name))
	}
	return &MultiLogger{&newLoggers}
}
