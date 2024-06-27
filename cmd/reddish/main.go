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

	respSRV := &resp.Server{
		Addr:    cfg.RESPServerAddr,
		Handler: resp.HandlerFunc(func(args []string) { fmt.Println(args) }),
	}
	httpSRV := &http.Server{
		Addr:    cfg.GRPCServerAddr,
		Handler: rpc.NewStorageServer(storageHandler),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()

		slog.Info("closing resp server")
		if err := respSRV.Close(); err != nil {
			slog.Error("failed to close resp server", "error", err)
		}
	}()
	go func() {
		<-ctx.Done()

		slog.Info("closing http server")
		if err := httpSRV.Close(); err != nil {
			slog.Error("failed to close http server", "error", err)
		}
	}()

	var wg sync.WaitGroup

	if len(cfg.RESPServerAddr) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			slog.Info("starting resp server", "address", respSRV.Addr)
			if err := respSRV.ListenAndServe(); err != nil && !errors.Is(err, resp.ErrServerClosed) {
				slog.Error("failed to start resp server", "error", err)
				os.Exit(1)
			}
		}()
	}

	if len(cfg.GRPCServerAddr) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			slog.Info("starting http server", "address", httpSRV.Addr)
			if err := httpSRV.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start http server", "error", err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}
