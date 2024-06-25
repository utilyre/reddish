package mapstorage

import (
	"context"
	"sync"

	"github.com/utilyre/reddish/internal/app"
	"github.com/utilyre/reddish/internal/app/domain"
)

type MapStorage struct {
	m  map[domain.Key]domain.Val
	mu sync.RWMutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{m: make(map[domain.Key]domain.Val)}
}

func (ms *MapStorage) Get(ctx context.Context, key domain.Key) (domain.Val, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	val, ok := ms.m[key]
	if !ok {
		return nil, app.ErrNoRecord
	}

	return val, nil
}

func (ms *MapStorage) Set(ctx context.Context, key domain.Key, val domain.Val) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.m[key] = val
	return nil
}

func (ms *MapStorage) Delete(ctx context.Context, key domain.Key) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.m[key]; !ok {
		return app.ErrNoRecord
	}

	delete(ms.m, key)
	return nil
}
