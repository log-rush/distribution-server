package distributionServer

import (
	"log"

	"github.com/log-rush/distribution-server/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.SugaredLogger
}

func createLogger() domain.Logger {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:   false,
		DisableCaller: true,
		Encoding:      "console",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "",
			FunctionKey:    "",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	l, err := config.Build()
	defer l.Sync()
	if err != nil {
		log.Fatalln(err)
	}

	logger := logger{l.Sugar()}
	return &logger
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
func (l *logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}
func (l *logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
func (l *logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}
func (l *logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}
func (l *logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}
func (l *logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *logger) Named(name string) domain.Logger {
	return &logger{l.logger.Named(name)}
}
