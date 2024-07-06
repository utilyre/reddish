package hashmap

import (
	"context"
	"sync"
	"time"

	"github.com/utilyre/reddish/internal/app"
	"github.com/utilyre/reddish/internal/app/domain"
)

type Hashmap struct {
	dict map[domain.Key]domain.Val
	exp  map[domain.Key]time.Time
	mu   sync.RWMutex
}

func New() *Hashmap {
	return &Hashmap{dict: make(map[domain.Key]domain.Val)}
}

func (ms *Hashmap) Exists(ctx context.Context, key domain.Key) error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	exp, ok := ms.exp[key]
	if ok && time.Now().After(exp) {
		if err := ms.Delete(ctx, key); err != nil {
			return err
		}

		return app.ErrNoRecord
	}

	if _, ok := ms.dict[key]; !ok {
		return app.ErrNoRecord
	}

	return nil
}

func (ms *Hashmap) Get(ctx context.Context, key domain.Key) (domain.Val, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	exp, ok := ms.exp[key]
	if ok && time.Now().After(exp) {
		if err := ms.Delete(ctx, key); err != nil {
			return nil, err
		}

		return nil, app.ErrNoRecord
	}

	val, ok := ms.dict[key]
	if !ok {
		return nil, app.ErrNoRecord
	}

	return val, nil
}

func (ms *Hashmap) Set(ctx context.Context, key domain.Key, val domain.Val, expiresAt time.Time) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if !expiresAt.IsZero() {
		ms.exp[key] = expiresAt
	}

	ms.dict[key] = val
	return nil
}

func (ms *Hashmap) Delete(ctx context.Context, key domain.Key) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.dict[key]; !ok {
		return app.ErrNoRecord
	}

	delete(ms.dict, key)
	return nil
}
