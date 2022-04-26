package memory

import (
	"context"

	"github.com/log-rush/simple-server/clients/repository"
	"github.com/log-rush/simple-server/domain"
)

type clientsMemoryRepository struct {
	clients map[string]domain.Client
}

func NewClientsMemoryrepository() domain.ClientsRepository {
	return &clientsMemoryRepository{
		clients: map[string]domain.Client{},
	}
}

func (u *clientsMemoryRepository) Create(ctx context.Context) (domain.Client, error) {
	id := repository.GenerateID()
	client := domain.Client{
		ID:      id,
		Send:    make(chan []byte),
		Receive: make(chan []byte),
	}
	u.clients[id] = client

	return client, nil
}

func (u *clientsMemoryRepository) GetClient(ctx context.Context, id string) (domain.Client, error) {
	client, ok := u.clients[id]
	if !ok {
		return domain.Client{}, domain.ErrClientNotFound
	}
	return client, nil
}

func (u *clientsMemoryRepository) Remove(ctx context.Context, id string) error {
	client, ok := u.clients[id]
	if !ok {
		return domain.ErrClientNotFound
	}
	close(client.Send)
	close(client.Receive)
	delete(u.clients, id)
	return nil
}
