package domain

import "context"

type SubscriptionsRepository interface {
	GetSubscribers(ctx context.Context, stream string) ([]Client, error)
	AddSubscription(ctx context.Context, stream string, client Client) error
	RemoveSubscription(ctx context.Context, stream, clientId string) error
	RemoveClient(ctx context.Context, clientId string) error
	RemoveStream(ctx context.Context, stream string) error
}
