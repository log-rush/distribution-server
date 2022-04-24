package memory

import (
	"context"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
)

type logRepository struct {
	logs map[string]*[]domain.Log
}

func NewLogRepository() domain.LogRepository {
	return &logRepository{
		logs: map[string]*[]domain.Log{},
	}
}

func (repo *logRepository) AddLogs(ctx context.Context, streamId string, logs *[]domain.Log) error {
	err := repo.ensureStream(ctx, streamId)
	if err != nil {
		return err
	}

	for idx := range *logs {
		(*logs)[idx].Stream = streamId
	}

	newSlice := append(*(repo.logs[streamId]), (*logs)...)
	repo.logs[streamId] = &newSlice
	return nil
}

func (repo *logRepository) SetLogs(ctx context.Context, streamId string, logs *[]domain.Log) error {
	err := repo.ensureStream(ctx, streamId)
	if err != nil {
		return err
	}

	for idx := range *logs {
		(*logs)[idx].Stream = streamId
	}

	repo.logs[streamId] = logs
	return nil
}

func (repo *logRepository) FetchLogs(ctx context.Context, streamId string) ([]domain.Log, error) {
	logs, ok := repo.logs[streamId]
	if !ok {
		return make([]domain.Log, 0), domain.ErrStreamNotFound
	}

	return *logs, nil
}

func (repo *logRepository) ensureStream(ctx context.Context, streamId string) error {
	_, ok := repo.logs[streamId]
	if ok {
		return nil
	}
	newStream := make([]domain.Log, 0)
	repo.logs[streamId] = &newStream
	return nil
}
