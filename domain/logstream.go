package domain

import "context"

type LogStream struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

type LogStreamUseCase interface {
	RegisterStream(ctx context.Context, alias string) (LogStream, error)
	UnregisterStream(ctx context.Context, id string) error
	SubscribeToStream(ctx context.Context, id string) error
	UnsubscribeFromStream(ctx context.Context, id string) error
}

type LogStreamRepository interface {
	CreateStream(ctx context.Context, alias string) (LogStream, error)
	GetStream(ctx context.Context, id string) (LogStream, error)
}
