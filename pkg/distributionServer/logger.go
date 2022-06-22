package distributionServer

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger() zap.SugaredLogger {
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
	return *l.Sugar()
}
