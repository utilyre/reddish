package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

	var wg sync.WaitGroup

	if len(cfg.RESPServerAddr) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			srv := &resp.Server{
				Addr:    cfg.RESPServerAddr,
				Handler: nil,
			}

			slog.Info("starting resp server", "address", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				slog.Error("failed to start resp server", "error", err)
				os.Exit(1)
			}
		}()
	}

	if len(cfg.GRPCServerAddr) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			srv := &http.Server{
				Addr:    cfg.GRPCServerAddr,
				Handler: rpc.NewStorageServer(storageHandler),
			}

			slog.Info("starting http server", "address", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start http server", "error", err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}
