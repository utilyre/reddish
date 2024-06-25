package mapstorage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/app"
)

func TestGet_norecord(t *testing.T) {
	ctx := context.Background()
	ms := mapstorage.NewMapStorage()

	_, err := ms.Get(ctx, "sample")
	if !errors.Is(err, app.ErrNoRecord) {
		t.Errorf("err = '%v'; want '%v'", err, app.ErrNoRecord)
	}
}

func TestDelete_norecord(t *testing.T) {
	ctx := context.Background()
	ms := mapstorage.NewMapStorage()

	err := ms.Delete(ctx, "sample")
	if !errors.Is(err, app.ErrNoRecord) {
		t.Errorf("err = '%v'; want '%v'", err, app.ErrNoRecord)
	}
}
