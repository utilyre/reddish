package service

import (
	"context"
	"fmt"
	"time"

	"github.com/utilyre/reddish/internal/app/domain"
	"github.com/utilyre/reddish/internal/app/port"
)

type StorageService struct {
	storageRepo port.StorageRepository
}

func NewStorageService(storageRepo port.StorageRepository) *StorageService {
	return &StorageService{storageRepo: storageRepo}
}

func (ss *StorageService) Exists(ctx context.Context, key string) error {
	k, err := domain.NewKey(key)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	if err := ss.storageRepo.Exists(ctx, k); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (ss *StorageService) Get(ctx context.Context, key string) ([]byte, error) {
	k, err := domain.NewKey(key)
	if err != nil {
		return nil, fmt.Errorf("domain: %w", err)
	}

	v, err := ss.storageRepo.Get(ctx, k)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	return []byte(v), nil
}

type SetOption func(opts *setOptions) error

type setOptions struct {
	expiresAt time.Time
}

func SetWithExpiresAt(expiresAt time.Time) SetOption {
	return func(opts *setOptions) error {
		opts.expiresAt = expiresAt
		return nil
	}
}

func (ss *StorageService) Set(ctx context.Context, key string, val []byte, opts ...SetOption) error {
	options := setOptions{
		expiresAt: time.Time{},
	}

	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return err
		}
	}

	k, err := domain.NewKey(key)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	v, err := domain.NewVal(val)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	if err := ss.storageRepo.Set(ctx, k, v, options.expiresAt); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (ss *StorageService) Delete(ctx context.Context, key string) error {
	k, err := domain.NewKey(key)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	if err := ss.storageRepo.Delete(ctx, k); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}
