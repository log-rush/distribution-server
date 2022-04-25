package domain

import "context"

type SubscriptionsRepository interface {
	GetSubscribers(ctx context.Context, stream string) ([]Client, error)
	RemoveSubscription(ctx context.Context, stream string, client string) ([]Client, error)
	RemoveClient(ctx context.Context, client string) error
	RemoveStream(ctx context.Context, stream string) error
}
