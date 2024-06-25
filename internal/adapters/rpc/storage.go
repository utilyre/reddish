package rpc

import (
	"context"

	"github.com/utilyre/reddish/internal/app/service"
	"github.com/utilyre/reddish/rpc/storage"
)

type StorageHandler struct {
	storageSVC *service.StorageService
}

func NewStorageHandler(storageSVC *service.StorageService) *StorageHandler {
	return &StorageHandler{storageSVC: storageSVC}
}

func (sh *StorageHandler) Get(
	ctx context.Context,
	r *storage.GetReq,
) (*storage.GetResp, error) {
	val, err := sh.storageSVC.Get(ctx, r.Key)
	if err != nil {
		return nil, err
	}

	return &storage.GetResp{Val: val}, nil
}

func (sh *StorageHandler) Set(
	ctx context.Context,
	r *storage.SetReq,
) (*storage.SetResp, error) {
	if err := sh.storageSVC.Set(ctx, r.Key, r.Val); err != nil {
		return nil, err
	}

	return &storage.SetResp{}, nil
}

func (sh *StorageHandler) Delete(
	ctx context.Context,
	r *storage.DeleteReq,
) (*storage.DeleteResp, error) {
	if err := sh.storageSVC.Delete(ctx, r.Key); err != nil {
		return nil, err
	}

	return &storage.DeleteResp{}, nil
}
