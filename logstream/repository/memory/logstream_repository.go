package repository

import (
	"context"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
	"github.com/fabiankachlock/log-rush-simple-server/logstream/repository"
)

type logStreamRepository struct {
	streams *map[string]domain.LogStream
}

func NewLogStreamRepository() domain.LogStreamRepository {
	return &logStreamRepository{
		streams: &map[string]domain.LogStream{},
	}
}

func (repo *logStreamRepository) CreateStream(ctx context.Context, alias string) (domain.LogStream, error) {
	entityId := repository.GenerateID()
	entity := domain.LogStream{
		ID:     entityId,
		Alias:  alias,
		Stream: make(domain.LogsChannel, 16), // TODO: move to config?
	}

	(*repo.streams)[entityId] = entity
	return entity, nil
}

func (repo *logStreamRepository) GetStream(ctx context.Context, id string) (domain.LogStream, error) {
	entity, ok := (*repo.streams)[id]
	if !ok {
		return domain.LogStream{}, domain.ErrStreamNotFound
	}
	return entity, nil
}

func (repo *logStreamRepository) DeleteStream(ctx context.Context, id string) error {
	_, ok := (*repo.streams)[id]
	if !ok {
		delete(*repo.streams, id)
		return nil
	}
	return domain.ErrStreamNotFound
}

func (repo *logStreamRepository) ListStreams(ctx context.Context) ([]domain.LogStream, error) {
	streams := make([]domain.LogStream, 0)
	for _, stream := range *repo.streams {
		streams = append(streams, stream)
	}
	return streams, nil
}
