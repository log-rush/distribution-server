package usecase

import (
	"context"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
)

type LogStreamUseCase struct{} // TODO: Add repos

func NewLogStreamUSeCase() domain.LogStreamUseCase {
	return &LogStreamUseCase{}
}

func (u LogStreamUseCase) RegisterStream(ctx context.Context, alias string) (domain.LogStream, error) {
	return domain.LogStream{
		ID:    generateID(),
		Alias: alias,
	}, nil
}

func (u LogStreamUseCase) UnregisterStream(ctx context.Context, id string) error {
	return nil
}

func (u LogStreamUseCase) SubscribeToStream(ctx context.Context, id string) error {
	return nil
}

func (u LogStreamUseCase) UnsubscribeFromStream(ctx context.Context, id string) error {
	return nil
}
