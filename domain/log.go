package domain

import (
	"context"

	logRush "github.com/log-rush/server-devkit"
)

type Log = logRush.Log

type LogUseCase interface {
	SendLog(ctx context.Context, stream string, log *Log) error
	SendLogBatch(ctx context.Context, stream string, logs *[]Log) error
}

type LogRepository interface {
	AddLogs(ctx context.Context, stream string, logs *[]Log) error
	FetchLogs(ctx context.Context, stream string) ([]Log, error)
}
