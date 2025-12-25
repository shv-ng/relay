package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shv-ng/relay/internal/algo"
	"github.com/shv-ng/relay/internal/pool"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Relay failed", "err", err)
		os.Exit(1)
	}
	slog.Info("Relay stopped cleanly")
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Received shutdown signal")
		cancel()
	}()

	urls := []string{"http://localhost:8000", "http://localhost:8001"}

	rr := algo.NewRoundRobin()
	sp := pool.New(ctx, rr, urls)
	go sp.HealthCheckStart()

	port := 8080

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			peer := sp.GetNext()
			if peer == nil {
				slog.Warn("No backends available", "path", r.URL.Path, "method", r.Method)
				http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
				return
			}
			peer.ReverseProxy.ServeHTTP(w, r)
		}),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "err", err)
		}
	}()

	slog.Info("Starting Relay", "port", port)

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
