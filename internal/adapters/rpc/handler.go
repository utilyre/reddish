package rpc

import (
	"context"

	"github.com/utilyre/reddish/internal/app/service"
)

type StorageHandler struct {
	storageSVC *service.StorageService
}

func NewStorageHandler(storageSVC *service.StorageService) *StorageHandler {
	return &StorageHandler{storageSVC: storageSVC}
}

func (sh *StorageHandler) Get(ctx context.Context, r *GetReq) (*GetResp, error) {
	val, err := sh.storageSVC.Get(ctx, r.Key)
	if err != nil {
		return nil, err
	}

	return &GetResp{Val: val}, nil
}

func (sh *StorageHandler) Set(ctx context.Context, r *SetReq) (*SetResp, error) {
	if err := sh.storageSVC.Set(ctx, r.Key, r.Val); err != nil {
		return nil, err
	}

	return &SetResp{}, nil
}

func (sh *StorageHandler) Delete(ctx context.Context, r *DeleteReq) (*DeleteResp, error) {
	if err := sh.storageSVC.Delete(ctx, r.Key); err != nil {
		return nil, err
	}

	return &DeleteResp{}, nil
}
