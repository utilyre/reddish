package service_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/utilyre/reddish/internal/adapters/hashmap"
	"github.com/utilyre/reddish/internal/app/domain"
	"github.com/utilyre/reddish/internal/app/service"
)

func TestStorageService_Set_empty(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	err := svc.Set(ctx, "", nil)
	if !errors.Is(err, domain.ErrEmpty) {
		t.Errorf("err = '%v'; want '%v'", err, domain.ErrEmpty)
	}
}

func TestStorageService_workflow(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	const key = "sample_key"
	const val = "sample_val"

	t.Log("running StorageService.Set")
	if err := svc.Set(ctx, key, []byte(val)); err != nil {
		t.Fatalf("err = '%v'; want '<nil>'", err)
	}

	t.Log("running StorageService.Get")
	v, err := svc.Get(ctx, key)
	if err != nil {
		t.Fatalf("err = '%v'; want '<nil>'", err)
	}
	if !slices.Equal(v, []byte(val)) {
		t.Fatalf("v = '%s'; want '%s'", v, val)
	}

	t.Log("running StorageService.Delete")
	if err := svc.Delete(ctx, key); err != nil {
		t.Fatalf("err = '%v'; want '<nil>'", err)
	}
}

func newStorageService() *service.StorageService {
	repo := hashmap.New()
	return service.NewStorageService(repo)
}
