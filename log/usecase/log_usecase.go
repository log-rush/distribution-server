package usecase

import (
	"context"
	"time"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
	"golang.org/x/sync/errgroup"
)

type logUseCase struct {
	logsRepo    domain.LogRepository
	streamsRepo domain.LogStreamRepository
	timeout     time.Duration
}

func NewLogUseCase(logsRepo domain.LogRepository, streamsRepo domain.LogStreamRepository, timeout time.Duration) domain.LogUseCase {
	return &logUseCase{
		logsRepo:    logsRepo,
		streamsRepo: streamsRepo,
	}
}

func (u *logUseCase) SendLog(ctx context.Context, streamId string, log *domain.Log) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	errGroup, context := errgroup.WithContext(_ctx)

	stream, err := u.streamsRepo.GetStream(context, streamId)
	if err != nil {
		return err
	}

	errGroup.Go(func() error {
		u.logsRepo.AddLogs(context, streamId, &[]domain.Log{*log})
		return nil
	})

	errGroup.Go(func() error {
		stream.Stream <- *log
		return nil
	})

	err = errGroup.Wait()
	return err
}

func (u *logUseCase) SendLogBatch(ctx context.Context, streamId string, logs *[]domain.Log) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	errGroup, context := errgroup.WithContext(_ctx)

	stream, err := u.streamsRepo.GetStream(context, streamId)
	if err != nil {
		return err
	}

	errGroup.Go(func() error {
		u.logsRepo.AddLogs(context, streamId, logs)
		return nil
	})

	for _, log := range *logs {
		log := log
		errGroup.Go(func() error {
			stream.Stream <- log
			return nil
		})
	}

	err = errGroup.Wait()
	return err
}
