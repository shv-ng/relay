package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shv-ng/relay/internal/algo"
	"github.com/shv-ng/relay/internal/config"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New("config.yml")
	if err != nil {
		return err
	}
	if err := setupLogging(cfg.LogFile); err != nil {
		return err
	}

	sp, err := setupPool(ctx, cfg)
	if err != nil {
		return err
	}

	return startServer(ctx, cfg.Port, sp)
}

func startServer(ctx context.Context, port int, sp *pool.Pool) error {
	srvErr := make(chan error, 1)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, sp)
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		slog.Info("Starting Relay", "port", port)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}

	case <-ctx.Done():
		slog.Info("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}
	return nil
}

func proxyRequest(w http.ResponseWriter, r *http.Request, sp *pool.Pool) {
	peerCtx := context.WithValue(r.Context(), algo.ClientIPKey, getClientIP(r))
	peer := sp.GetNext(peerCtx)

	if peer == nil {
		slog.Warn("No backends available", "path", r.URL.Path, "method", r.Method)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	peer.ActiveConnection.Add(1)
	defer peer.ActiveConnection.Add(-1)

	peer.ReverseProxy.ServeHTTP(w, r)
}

func setupPool(ctx context.Context, cfg *config.Config) (*pool.Pool, error) {

	algorithm := getAlgo(cfg.Algorithm)
	if algorithm == nil {
		return nil, fmt.Errorf("invalid algorithm %q: must be one of: round-robin, weighted-round-robin, random, least-conn, ip-hash", cfg.Algorithm)
	}

	sp := pool.New(ctx, algorithm, cfg.HealthCheckInterval, cfg.HealthCheckTimeout)
	sp.AddBackends(cfg.Backends)
	go sp.HealthCheckStart()

	return sp, nil
}

func setupLogging(logFile string) error {

	// logs will store here
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), nil))
	slog.SetDefault(logger)
	return nil
}

func getClientIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	return r.RemoteAddr
}

func getAlgo(algoStr string) algo.Picker {
	switch strings.ToLower(algoStr) {
	case "round-robin":
		return algo.NewRoundRobin()
	case "weighted-round-robin":
		return algo.NewWeightedRoundRobin()
	case "random":
		return algo.NewRandom()
	case "least-conn":
		return algo.NewLeastConnection()
	case "ip-hash":
		return algo.NewIPHash()
	}
	return nil
}
