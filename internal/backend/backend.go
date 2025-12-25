package backend

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy httputil.ReverseProxy

	alive atomic.Bool

	Weight          int64
	EffectiveWeight atomic.Int64

	ActiveConnection atomic.Int64
}

func New(url *url.URL, weight int64) *Backend {
	if weight <= 0 {
		weight = 1
	}
	return &Backend{
		URL:          url,
		ReverseProxy: *httputil.NewSingleHostReverseProxy(url),
		Weight:       weight,
	}
}

// HealthCheck checks and update health of each Backend with each 10sec interval
func (b *Backend) HealthCheck(timeout time.Duration) {
	client := http.Client{Timeout: timeout}

	healthUrl := b.URL.JoinPath("/health")
	res, err := client.Get(healthUrl.String())

	if err != nil {
		slog.Warn("Health check failed", "url", b.URL.String(), "err", err)

		b.SetAlive(false)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		slog.Info("Health check passed", "url", b.URL.String())

		b.SetAlive(true)
		return
	}

	slog.Warn("Health check failed", "url", b.URL.String(), "err", err)
	b.SetAlive(false)
}

func (b *Backend) Alive() bool {
	return b.alive.Load()
}

func (b *Backend) SetAlive(alive bool) {
	b.alive.Store(alive)
}
