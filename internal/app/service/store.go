package service

import (
	"context"
	"fmt"

	"github.com/utilyre/reddish/internal/app/domain"
	"github.com/utilyre/reddish/internal/app/port"
)

type StorageService struct {
	storageRepo port.StorageRepo
}

func NewStoreService(storageRepo port.StorageRepo) *StorageService {
	return &StorageService{storageRepo: storageRepo}
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

func (ss *StorageService) Set(ctx context.Context, key string, val []byte) error {
	k, err := domain.NewKey(key)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	v, err := domain.NewVal(val)
	if err != nil {
		return fmt.Errorf("domain: %w", err)
	}

	if err := ss.storageRepo.Set(ctx, k, v); err != nil {
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
