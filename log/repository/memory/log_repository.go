package memory

import (
	"context"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/log/repository"
)

type logRepository struct {
	logs            map[string]*repository.MaxLenQueue[domain.Log]
	maxStoredAmount int
}

func NewLogRepository(amountOfStoredLogs int) domain.LogRepository {
	return &logRepository{
		logs:            map[string]*repository.MaxLenQueue[domain.Log]{},
		maxStoredAmount: amountOfStoredLogs,
	}
}

func (repo *logRepository) AddLogs(ctx context.Context, streamId string, logs *[]domain.Log) error {
	err := repo.ensureStream(ctx, streamId)
	if err != nil {
		return err
	}

	for idx := range *logs {
		(*logs)[idx].Stream = streamId
		repo.logs[streamId].Enqueue((*logs)[idx])
	}

	return nil
}

func (repo *logRepository) FetchLogs(ctx context.Context, streamId string) ([]domain.Log, error) {
	logs, ok := repo.logs[streamId]
	if !ok {
		return []domain.Log{}, domain.ErrStreamNotFound
	}

	return (*logs).GetAll(), nil
}

func (repo *logRepository) ensureStream(ctx context.Context, streamId string) error {
	_, ok := repo.logs[streamId]
	if ok {
		return nil
	}

	queue := repository.NewMaxLenQueue(repo.maxStoredAmount, func() domain.Log {
		return domain.Log{}
	})
	repo.logs[streamId] = &queue
	return nil
}
