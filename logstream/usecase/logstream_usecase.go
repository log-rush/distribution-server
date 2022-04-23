package usecase

import (
	"context"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
)

type logStreamUseCase struct{} // TODO: Add repos

func NewLogStreamUSeCase() domain.LogStreamUseCase {
	return &logStreamUseCase{}
}

func (u logStreamUseCase) RegisterStream(ctx context.Context, alias string) (domain.LogStream, error) {
	return domain.LogStream{
		ID:    generateID(),
		Alias: alias,
	}, nil
}

func (u logStreamUseCase) UnregisterStream(ctx context.Context, id string) error {
	return nil
}

func (u logStreamUseCase) SubscribeToStream(ctx context.Context, id string) error {
	return nil
}

func (u logStreamUseCase) UnsubscribeFromStream(ctx context.Context, id string) error {
	return nil
}
