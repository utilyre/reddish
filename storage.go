package main

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Storage interface {
	Set(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type MapStorage struct {
	m  map[string][]byte
	mu sync.RWMutex
}

func NewMapStorage() *MapStorage {
	return &MapStorage{m: make(map[string][]byte)}
}

func (ms *MapStorage) Set(ctx context.Context, key string, val []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.m[key] = val
	return nil
}

func (ms *MapStorage) Get(ctx context.Context, key string) ([]byte, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	val, ok := ms.m[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return val, nil
}
