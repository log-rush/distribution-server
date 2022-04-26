package usecase

import (
	"context"
	"fmt"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/lrp"
)

type clientsUseCase struct {
	repo    domain.ClientsRepository
	decoder lrp.LRPDecoder
	l       *domain.Logger
}

func NewClientsUseCase(clientsRepo domain.ClientsRepository, logger domain.Logger) domain.ClientsUseCase {
	return &clientsUseCase{
		repo:    clientsRepo,
		decoder: lrp.NewDecoder(),
		l:       &logger,
	}
}

func (u *clientsUseCase) NewClient(ctx context.Context) (domain.Client, error) {
	client, err := u.repo.Create(ctx)
	if err != nil {
		(*u.l).Errorf("error while creating client: %s", err.Error())
		return domain.Client{}, err
	}
	(*u.l).Infof("created client %s", client.ID)

	go func() {
		(*u.l).Debugf("[%s] started request listener", client.ID)
		for msg := range client.Receive {
			(*u.l).Debugf("[%s] received %s ", client.ID, msg)
			u.handleMessage(msg)
		}
		(*u.l).Debugf("[%s] stopped request listener", client.ID)
	}()

	return client, nil
}

func (u *clientsUseCase) DestroyClient(ctx context.Context, id string) error {
	err := u.repo.Remove(ctx, id)
	if err != nil {
		(*u.l).Errorf("error while deleting %s client: %s", id, err.Error())
		return err
	}
	(*u.l).Infof("delted client %s", id)
	return nil
}

func (u *clientsUseCase) handleMessage(msg []byte) {
	message, err := u.decoder.Decode(msg)
	if err != nil {
		// do something
		return
	}
	fmt.Printf("received: %b, %s\n", message.OPCode, message.Payload)
}
