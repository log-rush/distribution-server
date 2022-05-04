package memory

import (
	"context"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/logstream/repository"
)

type logStreamRepository struct {
	streams           *map[string]domain.LogStream
	logsChannelBuffer int
}

func NewLogStreamRepository(logsChannelBuffer int) domain.LogStreamRepository {
	return &logStreamRepository{
		streams:           &map[string]domain.LogStream{},
		logsChannelBuffer: logsChannelBuffer,
	}
}

func (repo *logStreamRepository) CreateStream(ctx context.Context, alias, id, key string) (domain.LogStream, error) {
	entityId := id
	if entityId == "" {
		entityId = repository.GenerateID()
	}
	entityKey := key
	if entityKey == "" {
		entityKey = repository.GenerateID()
	}

	entity := domain.LogStream{
		ID:        entityId,
		Alias:     alias,
		Stream:    make(domain.LogsChannel, repo.logsChannelBuffer),
		SecretKey: entityKey,
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
	stream, ok := (*repo.streams)[id]
	if ok {
		close(stream.Stream)
		delete(*repo.streams, id)
		return nil
	}
	return domain.ErrStreamNotFound
}

func (repo *logStreamRepository) ListStreams(ctx context.Context) ([]domain.LogStream, error) {
	streams := make([]domain.LogStream, len(*repo.streams))
	for _, stream := range *repo.streams {
		streams = append(streams, stream)
	}
	return streams, nil
}
