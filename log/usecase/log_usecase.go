package usecase

import (
	"context"
	"time"

	"github.com/log-rush/simple-server/domain"
	"golang.org/x/sync/errgroup"
)

type logUseCase struct {
	logsRepo    domain.LogRepository
	streamsRepo domain.LogStreamRepository
	timeout     time.Duration
	l           *domain.Logger
}

func NewLogUseCase(logsRepo domain.LogRepository, streamsRepo domain.LogStreamRepository, timeout time.Duration, logger domain.Logger) domain.LogUseCase {
	return &logUseCase{
		logsRepo:    logsRepo,
		streamsRepo: streamsRepo,
		l:           &logger,
	}
}

func (u *logUseCase) SendLog(ctx context.Context, streamId string, log *domain.Log) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	errGroup, context := errgroup.WithContext(_ctx)

	stream, err := u.streamsRepo.GetStream(context, streamId)
	if err != nil {
		(*u.l).Errorf("error while addding log: %s", err.Error())
		return err
	}

	errGroup.Go(func() error {
		return u.logsRepo.AddLogs(context, streamId, &[]domain.Log{*log})
	})

	errGroup.Go(func() error {
		stream.Stream <- []domain.Log{*log}
		return nil
	})

	return errGroup.Wait()
}

func (u *logUseCase) SendLogBatch(ctx context.Context, streamId string, logs *[]domain.Log) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	errGroup, context := errgroup.WithContext(_ctx)

	stream, err := u.streamsRepo.GetStream(context, streamId)
	if err != nil {
		(*u.l).Errorf("error while batching logs: %s", err.Error())
		return err
	}

	errGroup.Go(func() error {
		return u.logsRepo.AddLogs(context, streamId, logs)
	})

	errGroup.Go(func() error {
		stream.Stream <- *logs
		return nil
	})

	return errGroup.Wait()
}
