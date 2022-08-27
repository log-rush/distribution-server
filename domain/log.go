package domain

import (
	"context"
)

type Log struct {
	Message   string `json:"message"`
	TimeStamp int    `json:"timestamp"`
	Stream    string `json:"stream"`
}
type LogUseCase interface {
	SendLog(ctx context.Context, stream string, log *Log) error
	SendLogBatch(ctx context.Context, stream string, logs *[]Log) error
}

type LogRepository interface {
	AddLogs(ctx context.Context, stream string, logs *[]Log) error
	FetchLogs(ctx context.Context, stream string) ([]Log, error)
}
