package service_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/app"
	"github.com/utilyre/reddish/internal/app/domain"
	"github.com/utilyre/reddish/internal/app/service"
)

func TestStorageService_Get_norecord(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	_, err := svc.Get(ctx, "sample")
	if !errors.Is(err, app.ErrNoRecord) {
		t.Errorf("err = '%v'; want '%v'", err, app.ErrNoRecord)
	}
}

func TestStorageService_Set_empty(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	err := svc.Set(ctx, "", nil)
	if !errors.Is(err, domain.ErrEmpty) {
		t.Errorf("err = '%v'; want '%v'", err, domain.ErrEmpty)
	}
}

func TestStorageService_Delete_norecord(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	err := svc.Delete(ctx, "sample")
	if !errors.Is(err, app.ErrNoRecord) {
		t.Errorf("err = '%v'; want '%v'", err, app.ErrNoRecord)
	}
}

func TestStorageService_workflow(t *testing.T) {
	ctx := context.Background()
	svc := newStorageService()

	const key = "sample_key"
	const val = "sample_val"

	t.Log("running StorageService.Set")
	if err := svc.Set(ctx, key, []byte(val)); err != nil {
		t.Errorf("err = '%v'; want '<nil>'", err)
	}

	t.Log("running StorageService.Get")
	v, err := svc.Get(ctx, key)
	if err != nil {
		t.Errorf("err = '%v'; want '<nil>'", err)
	}
	if !slices.Equal(v, []byte(val)) {
		t.Errorf("v = '%s'; want '%s'", v, val)
	}

	t.Log("running StorageService.Delete")
	if err := svc.Delete(ctx, key); err != nil {
		t.Errorf("err = '%v'; want '<nil>'", err)
	}
}

func newStorageService() *service.StorageService {
	repo := mapstorage.NewMapStorage()
	return service.NewStorageService(repo)
}
