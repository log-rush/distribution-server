package usecase

import (
	"context"
	"time"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/lrp"
	"golang.org/x/sync/errgroup"
)

type clientsUseCase struct {
	cRepo               domain.ClientsRepository
	sRepo               domain.SubscriptionsRepository
	decoder             lrp.LRPDecoder
	encoder             lrp.LRPEncoder
	timeout             time.Duration
	clientCheckInterval time.Duration
	maxResponseLatency  time.Duration
	l                   *domain.Logger
}

type extendedClient struct {
	domain.Client
	lastCheck int64
}

func NewClientsUseCase(
	clientsRepo domain.ClientsRepository,
	subscriptions domain.SubscriptionsRepository,
	clientCheckInterval time.Duration,
	maxResponseLatency time.Duration,
	timeout time.Duration,
	logger domain.Logger,
) domain.ClientsUseCase {
	return &clientsUseCase{
		cRepo:               clientsRepo,
		sRepo:               subscriptions,
		decoder:             lrp.NewDecoder(),
		encoder:             lrp.NewEncoder(),
		clientCheckInterval: clientCheckInterval,
		maxResponseLatency:  maxResponseLatency,
		timeout:             timeout,
		l:                   &logger,
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

	eClient := extendedClient{
		Client:    client,
		lastCheck: time.Now().UnixMilli(),
	}

	timer := time.NewTicker(u.clientCheckInterval)
	go func() {
		(*u.l).Debugf("[%s] started request listener", client.ID)
	outer:
		for {
			select {
			case <-timer.C:
				go func() {
					client.Send <- u.encoder.Encode(lrp.NewMesssage(lrp.OprStillAlive, []byte{}))
					eClient.lastCheck = time.Now().UnixMilli()
					(*u.l).Warnf("[%s] checking if alive", client.ID)
					<-time.After(u.maxResponseLatency)
					if eClient.lastCheck > 0 {
						// client did not respond in time
						(*u.l).Warnf("[%s] client inactive", client.ID)
						// TODO: close ws connection
					}
				}()
			case msg := <-client.Receive:
				if msg != nil {
					(*u.l).Debugf("[%s] received %s ", client.ID, msg)
					u.handleMessage(msg, &eClient)
				} else {
					break outer
				}
			}
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

func (u *clientsUseCase) handleError(err error, from *extendedClient) bool {
	if err != nil {
		(*u.l).Warnf("[%s] received errornous message (result: %s)", from.ID, err.Error())
		from.Send <- u.encoder.Encode(lrp.LRPMessage{OPCode: lrp.OprErr, Payload: []byte(err.Error())})
		return true
	}
	return false
}

func (u *clientsUseCase) handleMessage(msg []byte, from *extendedClient) {
	message, err := u.decoder.Decode(msg)
	if u.handleError(err, from) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if message.OPCode == lrp.OprSubscribe {
		err := u.sRepo.AddSubscription(ctx, string(message.Payload), from.Client)
		u.handleError(err, from)
		(*u.l).Infof("[%s] subscribed %s", from.ID, string(message.Payload))
	} else if message.OPCode == lrp.OprUnsubscribe {
		err := u.sRepo.RemoveSubscription(ctx, string(message.Payload), from.ID)
		u.handleError(err, from)
		(*u.l).Infof("[%s] unsubscribed %s", from.ID, string(message.Payload))
	} else if message.OPCode == lrp.OprAlive {
		from.lastCheck = -1
		(*u.l).Infof("[%s] still alive", from.ID)
	}
}
