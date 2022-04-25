package usecase

import (
	"context"
	"time"

	"github.com/log-rush/simple-server/domain"
)

type logStreamUseCase struct {
	streamsRepo domain.LogStreamRepository
	timeout     time.Duration
	pool        logDistributionWorkerPool
}

func NewLogStreamUseCase(repo domain.LogStreamRepository, timeout time.Duration) domain.LogStreamUseCase {
	u := &logStreamUseCase{
		streamsRepo: repo,
		timeout:     timeout,
		pool:        logDistributionWorkerPool{},
	}
	u.pool.Start()
	return u
}

func (u *logStreamUseCase) RegisterStream(ctx context.Context, alias string) (domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	stream, err := u.streamsRepo.CreateStream(context, alias)
	if err != nil {
		return domain.LogStream{}, err
	}

	go func() {
		for log := range stream.Stream {
			u.pool.PostJob(log, stream.ID)
		}
	}()

	return stream, nil
}

func (u *logStreamUseCase) UnregisterStream(ctx context.Context, id string) error {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	err := u.streamsRepo.DeleteStream(context, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *logStreamUseCase) SubscribeToStream(ctx context.Context, id string) (domain.LogsChannel, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	stream, err := u.streamsRepo.GetStream(context, id)
	if err != nil {
		return nil, err
	}

	return stream.Stream, nil
}

func (u *logStreamUseCase) GetAvailableStreams(ctx context.Context) ([]domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	streams, err := u.streamsRepo.ListStreams(context)
	if err != nil {
		return nil, err
	}

	return streams, nil
}
