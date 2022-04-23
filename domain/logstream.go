package domain

import "context"

type LogStream struct {
	ID     string      `json:"id"`
	Alias  string      `json:"alias"`
	Stream LogsChannel `json:"-"`
}

type LogsChannel chan Log

type LogStreamUseCase interface {
	RegisterStream(ctx context.Context, alias string) (LogStream, error)
	UnregisterStream(ctx context.Context, id string) error
	SubscribeToStream(ctx context.Context, id string) (LogsChannel, error)
	GetAvailableStreams(ctx context.Context) ([]LogStream, error)
}

type LogStreamRepository interface {
	CreateStream(ctx context.Context, alias string) (LogStream, error)
	DeleteStream(ctx context.Context, alias string) error
	GetStream(ctx context.Context, id string) (LogStream, error)
	ListStreams(ctx context.Context) ([]LogStream, error)
}
