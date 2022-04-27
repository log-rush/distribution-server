package usecase

import (
	"context"
	"time"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/lrp"
	"golang.org/x/sync/errgroup"
)

type logStreamUseCase struct {
	streamsRepo       domain.LogStreamRepository
	subscriptionsRepo domain.SubscriptionsRepository
	pool              logDistributionWorkerPool
	encoder           lrp.LRPEncoder
	timeout           time.Duration
	l                 *domain.Logger
}

func NewLogStreamUseCase(repo domain.LogStreamRepository, supscriptions domain.SubscriptionsRepository, timeout time.Duration, logger domain.Logger) domain.LogStreamUseCase {
	u := &logStreamUseCase{
		streamsRepo:       repo,
		subscriptionsRepo: supscriptions,
		timeout:           timeout,
		pool:              NewPool(5, &supscriptions, logger),
		l:                 &logger,
		encoder:           lrp.NewEncoder(),
	}
	u.pool.Start()
	return u
}

func (u *logStreamUseCase) RegisterStream(ctx context.Context, alias string) (domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	stream, err := u.streamsRepo.CreateStream(context, alias)
	if err != nil {
		(*u.l).Errorf("error while creating stream: %s", err.Error())
		return domain.LogStream{}, err
	}
	(*u.l).Infof("created stream %s", stream.ID)

	go func() {
		(*u.l).Debugf("[%s] starting log listener", stream.ID)
		for log := range stream.Stream {
			(*u.l).Debugf("[%s] received log %s ", stream.ID, log.Message)
			u.pool.PostJob(log, stream.ID)
		}
		(*u.l).Debugf("[%s] stopped log listener", stream.ID)
	}()

	return stream, nil
}

func (u *logStreamUseCase) UnregisterStream(ctx context.Context, id string) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	errGroup, context := errgroup.WithContext(_ctx)
	errGroup.Go(func() error {
		err := u.streamsRepo.DeleteStream(context, id)
		if err != nil {
			(*u.l).Errorf("error while deleting stream %s: %s", id, err.Error())
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		subscribers, err := u.subscriptionsRepo.GetSubscribers(context, id)
		if err != nil {
			return err
		}
		err = u.subscriptionsRepo.RemoveStream(context, id)
		if err != nil {
			return err
		}
		for _, client := range subscribers {
			client := client
			errGroup.Go(func() error {
				client.Send <- u.encoder.Encode(lrp.NewMesssage(lrp.OprUnsubscribe, []byte(id)))
				return nil
			})
		}
		return nil
	})

	err := errGroup.Wait()
	if err != nil {
		(*u.l).Infof("deleted stream %s", id)
	}
	return err
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
		(*u.l).Errorf("error while listing streams: %s", err.Error())
		return nil, err
	}

	return streams, nil
}
