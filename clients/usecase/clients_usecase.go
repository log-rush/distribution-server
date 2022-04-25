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
}

func NewClientsUseCase(clientsRepo domain.ClientsRepository) domain.ClientsUseCase {
	return &clientsUseCase{
		repo:    clientsRepo,
		decoder: lrp.NewDecoder(),
	}
}

func (u *clientsUseCase) NewClient(ctx context.Context) (domain.Client, error) {
	client, err := u.repo.Create(ctx)
	if err != nil {
		return domain.Client{}, err
	}

	go func() {
		for msg := range client.Receive {
			u.handleMessage(msg)
		}
	}()

	return client, nil
}

func (u *clientsUseCase) DestroyClient(ctx context.Context, id string) error {
	return u.repo.Remove(ctx, id)
}

func (u *clientsUseCase) handleMessage(msg []byte) {
	message, err := u.decoder.Decode(msg)
	if err != nil {
		// do something
		return
	}
	fmt.Printf("received: %b, %s\n", message.OPCode, message.Payload)
}
