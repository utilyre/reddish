package port

import (
	"context"
	"time"

	"github.com/utilyre/reddish/internal/app/domain"
)

type StorageRepository interface {
	Get(
		ctx context.Context,
		key domain.Key,
	) (domain.Val, error)

	Set(
		ctx context.Context,
		key domain.Key,
		val domain.Val,
		expiresAt time.Time,
	) error

	Delete(
		ctx context.Context,
		key domain.Key,
	) error
}
