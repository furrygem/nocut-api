package links

import "context"

type Storage interface {
	Create(ctx context.Context, link Link) (string, bool, error)
	FindOne(ctx context.Context, id string) (Link, error)
	FindOneBySource(ctx context.Context, source string) (Link, error)
	Delete(ctx context.Context, id string) error
}
