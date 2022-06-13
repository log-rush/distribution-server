package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/commons"
	"github.com/log-rush/simple-server/pkg/lrp"
	"golang.org/x/sync/errgroup"
)

type logStreamUseCase struct {
	lsRepo  domain.LogStreamRepository
	sRepo   domain.SubscriptionsRepository
	pool    logDistributionWorkerPool
	encoder lrp.LRPEncoder
	timeout time.Duration
	l       *domain.Logger
}

func NewLogStreamUseCase(repo domain.LogStreamRepository, supscriptions domain.SubscriptionsRepository, maxAmountOfWorkers int, timeout time.Duration, logger domain.Logger) domain.LogStreamUseCase {
	u := &logStreamUseCase{
		lsRepo:  repo,
		sRepo:   supscriptions,
		timeout: timeout,
		pool:    NewPool(maxAmountOfWorkers, &supscriptions, logger),
		l:       &logger,
		encoder: lrp.NewEncoder(),
	}
	u.pool.Start()
	return u
}

func (u *logStreamUseCase) RegisterStream(ctx context.Context, alias, id, key string) (domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	existingStream, err := u.lsRepo.GetStream(context, id)
	if err != nil && !errors.Is(err, domain.ErrStreamNotFound) {
		(*u.l).Warnf("error while fetching stream %s: %s", id, err.Error())
		return domain.LogStream{}, err
	}

	if err == nil && existingStream.SecretKey != key {
		return domain.LogStream{}, domain.ErrRegisterNotAllowed
	} else if err == nil {
		return existingStream, nil
	}

	stream, err := u.lsRepo.CreateStream(context, alias, id, key)
	if err != nil {
		(*u.l).Errorf("error while creating stream: %s", err.Error())
		return domain.LogStream{}, err
	}
	(*u.l).Infof("created stream %s", stream.ID)

	go func() {
		defer commons.RecoverRoutine(u.l)
		defer (*u.l).Debugf("[%s] stopped log listener", stream.ID)
		(*u.l).Debugf("[%s] starting log listener", stream.ID)
		for logs := range stream.Stream {
			(*u.l).Debugf("[%s] received logs (%d) ", stream.ID, len(logs))
			u.pool.PostJob(logs, stream.ID)
		}
	}()

	return stream, nil
}

func (u *logStreamUseCase) UnregisterStream(ctx context.Context, id, key string) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	stream, err := u.lsRepo.GetStream(_ctx, id)
	if err != nil {
		(*u.l).Errorf("error while fetching stream %s: %s", id, err.Error())
		return err
	}

	if stream.SecretKey != key {
		return domain.ErrUnregisterNotAllowed
	}

	errGroup, context := errgroup.WithContext(_ctx)
	errGroup.Go(func() error {
		err := u.lsRepo.DeleteStream(context, id)
		if err != nil {
			(*u.l).Errorf("error while deleting stream %s: %s", id, err.Error())
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		subscribers, err := u.sRepo.GetSubscribers(context, id)
		if err != nil {
			(*u.l).Warnf("error gettings stream subscribers %s: %s", id, err.Error())
			// discard error since the stream might not had any subscribers
			return nil
		}
		err = u.sRepo.RemoveStream(context, id)
		if err != nil {
			(*u.l).Errorf("error while delteting stream subscriptions %s: %s", id, err.Error())
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

	err = errGroup.Wait()
	if err != nil {
		(*u.l).Infof("deleted stream %s", id)
	}
	return err
}

func (u *logStreamUseCase) GetAvailableStreams(ctx context.Context) ([]domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	streams, err := u.lsRepo.ListStreams(context)
	if err != nil {
		(*u.l).Errorf("error while listing streams: %s", err.Error())
		return nil, err
	}

	return streams, nil
}

func (u *logStreamUseCase) GetStream(ctx context.Context, streamId string) (domain.LogStream, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	stream, err := u.lsRepo.GetStream(context, streamId)
	if err != nil {
		(*u.l).Warnf("error while fetching stream %s: %s", streamId, err.Error())
		return domain.LogStream{}, err
	}

	return stream, nil
}
