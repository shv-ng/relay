package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy httputil.ReverseProxy

	mu    sync.RWMutex
	alive bool
}

func NewBackend(u *url.URL) *Backend {
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

func (b *Backend) CheckHealth() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	u := b.URL.JoinPath("/health")
	res, err := client.Get(u.String())
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

type ServerPool struct {
	Backend []*Backend
	current uint64
	total   uint64
}

func (s *ServerPool) NextIndex() uint64 {
	return atomic.AddUint64(&s.current, uint64(1)) % s.total
}

func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()

	for i := next; i < s.total+next; i++ {
		idx := i % s.total

		if s.Backend[idx].Alive() {
			atomic.StoreUint64(&s.current, idx)
			return s.Backend[idx]
		}
	}
	return nil
}

func NewServerPool() *ServerPool {
	return &ServerPool{}
}

func (s *ServerPool) AddBackend(v ...string) {
	for _, u := range v {
		uu, err := url.Parse(u)
		if err != nil {
			slog.Error("URL parsing", "url", uu, "err", err)
			continue
		}
		s.total++

		b := NewBackend(uu)
		s.Backend = append(s.Backend, b)
		go b.CheckHealth()
	}
}
func (s *ServerPool) CheckHealth(done chan bool) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, b := range s.Backend {
					go b.CheckHealth()
				}
			case <-done:
				return
			}
		}
	}()

}

func lb(sp *ServerPool, w http.ResponseWriter, r *http.Request) {
	peer := sp.GetNextPeer()

	if peer == nil {
		slog.Warn("No backends available", "path", r.URL.Path, "method", r.Method)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}
	slog.Info("Routing request", "backend", peer.URL.String(), "path", r.URL.Path, "method", r.Method)
	peer.ReverseProxy.ServeHTTP(w, r)
}

func run() error {
	done := make(chan bool)
	defer func() { done <- true }()

	urls := []string{"http://localhost:8000"}

	sp := NewServerPool()
	slog.Info("Added backends", "count", len(urls), "backends", urls)
	sp.AddBackend(urls...)
	go sp.CheckHealth(done)

	port := 8080
	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lb(sp, w, r)
		}),
	}
	slog.Info("Starting load balancer", "port", port)
	return s.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start", "err", err)
		os.Exit(1)
	}
}
