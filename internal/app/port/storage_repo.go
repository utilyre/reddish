package port

import (
	"context"

	"github.com/utilyre/reddish/internal/app/domain"
)

type StorageRepo interface {
	Get(ctx context.Context, key domain.Key) (domain.Val, error)
	Set(ctx context.Context, key domain.Key, val domain.Val) error
	Delete(ctx context.Context, key domain.Key) error
}
