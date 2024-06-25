package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/adapters/rpc"
	"github.com/utilyre/reddish/internal/app/service"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)
}

func main() {
	storageRepo := mapstorage.NewMapStorage()
	storageSVC := service.NewStorageService(storageRepo)
	storageHandler := rpc.NewStorageHandler(storageSVC)

	srv := &http.Server{
		Addr:    ":5000",
		Handler: rpc.NewStorageServer(storageHandler),
	}
	slog.Info("starting storage server", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start storage server", "address", srv.Addr, "error", err)
		os.Exit(1)
	}
}
