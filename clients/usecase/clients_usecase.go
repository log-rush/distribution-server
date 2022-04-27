package usecase

import (
	"context"
	"time"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/lrp"
	"golang.org/x/sync/errgroup"
)

type clientsUseCase struct {
	cRepo   domain.ClientsRepository
	sRepo   domain.SubscriptionsRepository
	decoder lrp.LRPDecoder
	encoder lrp.LRPEncoder
	timeout time.Duration
	l       *domain.Logger
}

func NewClientsUseCase(clientsRepo domain.ClientsRepository, subscriptions domain.SubscriptionsRepository, timeout time.Duration, logger domain.Logger) domain.ClientsUseCase {
	return &clientsUseCase{
		cRepo:   clientsRepo,
		sRepo:   subscriptions,
		decoder: lrp.NewDecoder(),
		encoder: lrp.NewEncoder(),
		timeout: timeout,
		l:       &logger,
	}
}

func (u *clientsUseCase) NewClient(ctx context.Context) (domain.Client, error) {
	context, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	client, err := u.cRepo.Create(context)
	if err != nil {
		(*u.l).Errorf("error while creating client: %s", err.Error())
		return domain.Client{}, err
	}
	(*u.l).Infof("created client %s", client.ID)

	go func() {
		(*u.l).Debugf("[%s] started request listener", client.ID)
		for msg := range client.Receive {
			(*u.l).Debugf("[%s] received %s ", client.ID, msg)
			u.handleMessage(msg, &client)
		}
		(*u.l).Debugf("[%s] stopped request listener", client.ID)
	}()

	return client, nil
}

func (u *clientsUseCase) DestroyClient(ctx context.Context, id string) error {
	_ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	errGroup, context := errgroup.WithContext(_ctx)

	errGroup.Go(func() error {
		err := u.cRepo.Remove(context, id)
		if err != nil {
			(*u.l).Errorf("error while deleting %s client: %s", id, err.Error())
			return err
		}
		(*u.l).Infof("delted client %s", id)
		return nil
	})

	errGroup.Go(func() error {
		return u.sRepo.RemoveClient(context, id)
	})

	return errGroup.Wait()
}

func (u *clientsUseCase) handleError(err error, from *domain.Client) bool {
	if err != nil {
		(*u.l).Warnf("[%s] received errornous message (error: %s)", from.ID, err.Error())
		from.Send <- u.encoder.Encode(lrp.LRPMessage{OPCode: lrp.OprErr, Payload: []byte(err.Error())})
		return true
	}
	return false
}

func (u *clientsUseCase) handleMessage(msg []byte, from *domain.Client) {
	message, err := u.decoder.Decode(msg)
	if u.handleError(err, from) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if message.OPCode == lrp.OprSubscribe {
		err := u.sRepo.AddSubscription(ctx, string(message.Payload), *from)
		u.handleError(err, from)
		(*u.l).Infof("[%s] subscribed %s", from.ID, string(message.Payload))

	} else if message.OPCode == lrp.OprUnsubscribe {
		err := u.sRepo.RemoveSubscription(ctx, string(message.Payload), from.ID)
		u.handleError(err, from)
		(*u.l).Infof("[%s] unsubscribed %s", from.ID, string(message.Payload))
	}
}
