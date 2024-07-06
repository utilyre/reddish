package rpc

import (
	"context"
	"errors"

	"github.com/utilyre/reddish/internal/app"
	"github.com/utilyre/reddish/internal/app/service"
)

type StorageHandler struct {
	storageSVC *service.StorageService
}

func NewStorageHandler(storageSVC *service.StorageService) *StorageHandler {
	return &StorageHandler{storageSVC: storageSVC}
}

func (sh *StorageHandler) Exists(ctx context.Context, r *ExistsReq) (*ExistsResp, error) {
	var existed int64
	var errs []error

	for _, key := range r.Keys {
		if err := sh.storageSVC.Exists(ctx, key); err != nil && !errors.Is(err, app.ErrNoRecord) {
			errs = append(errs, err)
		}

		existed++
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &ExistsResp{NumExisted: existed}, nil
}

func (sh *StorageHandler) Get(ctx context.Context, r *GetReq) (*GetResp, error) {
	val, err := sh.storageSVC.Get(ctx, r.Key)
	if err != nil {
		return nil, err
	}

	return &GetResp{Val: val}, nil
}

func (sh *StorageHandler) Set(ctx context.Context, r *SetReq) (*SetResp, error) {
	if err := sh.storageSVC.Set(
		ctx,
		r.Key,
		r.Val,
		service.SetWithExpiresAt(r.ExpiresAt.AsTime()),
	); err != nil {
		return nil, err
	}

	return &SetResp{}, nil
}

func (sh *StorageHandler) Del(ctx context.Context, r *DelReq) (*DelResp, error) {
	var deleted int64
	var errs []error

	for _, key := range r.Keys {
		if err := sh.storageSVC.Delete(ctx, key); err != nil {
			errs = append(errs, err)
			continue
		}

		deleted++
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &DelResp{NumDeleted: deleted}, nil
}
