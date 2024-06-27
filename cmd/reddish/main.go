package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/utilyre/reddish/internal/adapters/hashmap"
	"github.com/utilyre/reddish/internal/adapters/rpc"
	"github.com/utilyre/reddish/internal/app/service"
	"github.com/utilyre/reddish/internal/config"
	"github.com/utilyre/reddish/pkg/resp"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: cfg.LogLevel},
	)))

	slog.Info("initializing services")
	storageRepo := hashmap.New()
	storageSVC := service.NewStorageService(storageRepo)
	storageHandler := rpc.NewStorageHandler(storageSVC)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		StartRESPServer(ctx, cfg)
	}()

	go func() {
		defer wg.Done()
		StartGRPCServer(ctx, cfg, storageHandler)
	}()

	wg.Wait()
}

func StartRESPServer(ctx context.Context, cfg *config.Config) {
	if len(cfg.RESPServerAddr) == 0 {
		return
	}

	srv := &resp.Server{
		Addr:    cfg.RESPServerAddr,
		Handler: resp.HandlerFunc(func(args []string) { fmt.Println(args) }),
	}

	go func() {
		slog.Info("starting resp server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, resp.ErrServerClosed) {
			slog.Error("failed to start resp server", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("closing resp server")
	if err := srv.Close(); err != nil {
		slog.Error("failed to close resp server", "error", err)
	}
}

func StartGRPCServer(ctx context.Context, cfg *config.Config, storage rpc.Storage) {
	if len(cfg.GRPCServerAddr) == 0 {
		return
	}

	srv := &http.Server{
		Addr:    cfg.GRPCServerAddr,
		Handler: rpc.NewStorageServer(storage),
	}

	go func() {
		slog.Info("starting grpc server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start grpc server", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("closing grpc server")
	if err := srv.Close(); err != nil {
		slog.Error("failed to close grpc server", "error", err)
	}
}
