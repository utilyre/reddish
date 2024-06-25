package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/adapters/rpc"
	"github.com/utilyre/reddish/internal/app/service"
	"github.com/utilyre/reddish/internal/config"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))
}

func main() {
	slog.Info("loading config", "path", ".env")
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	storageRepo := mapstorage.NewMapStorage()
	storageSVC := service.NewStorageService(storageRepo)
	storageHandler := rpc.NewStorageHandler(storageSVC)

	srv := &http.Server{
		Addr:    cfg.StorageServerAddr,
		Handler: rpc.NewStorageServer(storageHandler),
	}
	slog.Info("starting storage server", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start storage server", "address", srv.Addr, "error", err)
		os.Exit(1)
	}
}
