package links

import "context"

type Storage interface {
	Create(ctx context.Context, link Link) (string, error)
	FindOne(ctx context.Context, id string) (Link, error)
	Delete(ctx context.Context, id string) error
}
