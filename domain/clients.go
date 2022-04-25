package domain

import "context"

type Client struct {
	ID      string
	Send    chan []byte
	Receive chan []byte
}

type ClientsUseCase interface {
	NewClient(ctx context.Context) (Client, error)
	DestroyClient(ctx context.Context, id string) error
}

type ClientsRepository interface {
	Create(ctx context.Context) (Client, error)
	GetClient(ctx context.Context, id string) (Client, error)
	Remove(ctx context.Context, id string) error
}
