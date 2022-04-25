package usecase

import (
	"context"

	"github.com/log-rush/simple-server/domain"
)

type clientsUseCase struct {
	repo domain.ClientsRepository
}

func NewClientsUseCase(clientsRepo domain.ClientsRepository) domain.ClientsUseCase {
	return &clientsUseCase{
		repo: clientsRepo,
	}
}

func (u *clientsUseCase) NewClient(ctx context.Context) (domain.Client, error) {
	return u.repo.Create(ctx)
}

func (u *clientsUseCase) DestroyClient(ctx context.Context, id string) error {
	return u.repo.Remove(ctx, id)
}
