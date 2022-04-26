package memory

type subscriptionsRepository struct {
	subscribers map[string]*[]string
	isSubscribing map[string]*[]string
}

func NewSubscriptionsRepository() {
	return &subscriptionsRepository{}
}

func (repo *subscriptionsRepository) AddSubscription(ctx context.Context, streamId, clientId string) error {
	subscribers, ok1 := repo.subscribers[streamId]
	subscriptions, ok2 := repo.isSubscribing[streamId]
	if !ok1 || !ok2 {
		return domain.ErrStreamNotFound
	}

	repo.subscribers[streamId] = repository.AppendUniqueToSlice(subscribers, client)
	repo.isSubscribing[streamId] = repository.AppendUniqueToSlice(subscriptions, streamId)
	return nil
}


func (repo *subscriptionsRepository) RemoveSubscription(ctx context.Context, streamId, clientId string) error {
	subscribers, ok1 := repo.subscribers[streamId]
	subscriptions, ok2 := repo.isSubscribing[streamId]
	if !ok1 || !ok2 {
		return domain.ErrStreamNotFound
	}
	
	repo.subscribers[streamId] = repository.RemoveFromSlice(subscribers, client)
	repo.isSubscribing[streamId] = repository.RemoveFromSlice(subscriptions, streamId)
	return nil
}

func (repo *subscriptionsRepository) RemoveStream(ctx context.Context, streamId string) error {
	subscribers, ok := repo.subscribers[streamId]
	if !ok {
		return domain.ErrStreamNotFound
	}
	
	delete(repo.subscribers, streamId)
	for _, subscriber := range subscribers {
		client, ok := repo.isSubscribing[subscriber]
		if !ok {
			return domain.ErrClientNotFound
		}	
		repo.isSubscribing[subscriber] = repository.RemoveFromSlice(client, streamId)
	}
	return nil
}

func (repo *subscriptionsRepository) RemoveClient(ctx context.Context, clientId string) error {
	streams, ok := repo.isSubscribing[clientId]
	if !ok {
		return domain.ErrClientNotFound
	}
	
	delete(repo.isSubscribing, clientId)
	for _, stream := range streams {
		clients, ok := repo.subscribers[stream]
		if !ok {
			return domain.ErrStreamNotFound
		}	
		repo.isSubscribing[stream] = repository.RemoveFromSlice(clients, clientId)
	}
	return nil
}