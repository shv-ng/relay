package backend

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy httputil.ReverseProxy

	mu    sync.RWMutex
	alive bool
}

func New(url *url.URL) *Backend {
	return &Backend{
		URL:          url,
		ReverseProxy: *httputil.NewSingleHostReverseProxy(url),
	}
}

// HealthCheck checks and update health of each Backend with each 10sec interval
func (b *Backend) HealthCheck() {
	client := http.Client{Timeout: 3 * time.Second}

	healthUrl := b.URL.JoinPath("/health")
	res, err := client.Get(healthUrl.String())

	if err != nil {
		slog.Warn("Health check failed", "url", b.URL.String(), "err", err)

		b.SetAlive(false)
		return
	}

	if res.StatusCode == http.StatusOK {
		slog.Info("Health check passed", "url", b.URL.String())

		b.SetAlive(true)
		return
	}

	slog.Warn("Health check failed", "url", b.URL.String(), "err", err)
	b.SetAlive(false)
}

func (b *Backend) Alive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.alive = alive
}
