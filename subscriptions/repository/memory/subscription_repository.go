package memory

import (
	"context"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/app"
	"github.com/log-rush/distribution-server/subscriptions/repository"
)

type subscriptionsRepository struct {
	subscribers   map[string]*[]domain.Client
	isSubscribing map[string]*[]string
	lsRepo        domain.LogStreamRepository
}

func NewSubscriptionsRepository(context *app.Context) domain.SubscriptionsRepository {
	return &subscriptionsRepository{
		lsRepo:        context.Repos.LogStream,
		subscribers:   map[string]*[]domain.Client{},
		isSubscribing: map[string]*[]string{},
	}
}

func (repo *subscriptionsRepository) GetSubscribers(ctx context.Context, streamId string) ([]domain.Client, error) {
	subscribers, ok := repo.subscribers[streamId]
	if !ok {
		return []domain.Client{}, domain.ErrStreamNotFound
	}

	return *subscribers, nil
}

func (repo *subscriptionsRepository) AddSubscription(ctx context.Context, streamId string, client domain.Client) error {
	subscribers, ok1 := repo.subscribers[streamId]
	subscriptions, ok2 := repo.isSubscribing[client.ID]
	if !ok1 {
		_, err := repo.lsRepo.GetStream(ctx, streamId)
		if err != nil {
			return err
		}
		subscribers = &[]domain.Client{}
	}
	if !ok2 {
		subscriptions = &[]string{}
	}

	repo.subscribers[streamId] = repository.AppendUniqueToSlice(subscribers, client, repo.clientComperator(client.ID))
	repo.isSubscribing[streamId] = repository.AppendUniqueToSlice(subscriptions, streamId, repo.stringComperator(streamId))
	return nil
}

func (repo *subscriptionsRepository) RemoveSubscription(ctx context.Context, streamId, clientId string) error {
	subscribers, ok1 := repo.subscribers[streamId]
	subscriptions, ok2 := repo.isSubscribing[streamId]
	if !ok1 || !ok2 {
		return domain.ErrStreamNotFound
	}

	repo.subscribers[streamId] = repository.RemoveFromSlice(subscribers, repo.clientComperator(clientId))
	repo.isSubscribing[streamId] = repository.RemoveFromSlice(subscriptions, repo.stringComperator(streamId))
	return nil
}

func (repo *subscriptionsRepository) RemoveStream(ctx context.Context, streamId string) error {
	subscribers, ok := repo.subscribers[streamId]
	if !ok {
		return domain.ErrStreamNotFound
	}

	delete(repo.subscribers, streamId)
	for _, subscriber := range *subscribers {
		client, ok := repo.isSubscribing[subscriber.ID]
		if !ok {
			continue
		}
		repo.isSubscribing[subscriber.ID] = repository.RemoveFromSlice(client, repo.stringComperator(streamId))
	}
	return nil
}

func (repo *subscriptionsRepository) RemoveClient(ctx context.Context, clientId string) error {
	streams, ok := repo.isSubscribing[clientId]
	if !ok {
		return domain.ErrClientNotFound
	}

	delete(repo.isSubscribing, clientId)
	for _, stream := range *streams {
		clients, ok := repo.subscribers[stream]
		if !ok {
			return domain.ErrStreamNotFound
		}
		repo.subscribers[stream] = repository.RemoveFromSlice(clients, repo.clientComperator(clientId))
	}
	return nil
}

func (repo *subscriptionsRepository) clientComperator(clientId string) func(v domain.Client) bool {
	return func(v domain.Client) bool {
		return v.ID == clientId
	}
}

func (repo *subscriptionsRepository) stringComperator(str string) func(v string) bool {
	return func(v string) bool {
		return str == v
	}
}
