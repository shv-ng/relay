package pool

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/shv-ng/relay/internal/algo"
	"github.com/shv-ng/relay/internal/backend"
	"github.com/shv-ng/relay/internal/config"
)

type Pool struct {
	ctx context.Context

	backends []*backend.Backend
	selector algo.Picker

	interval time.Duration
	timeout  time.Duration
}

func New(ctx context.Context, selector algo.Picker, interval, timeout int) *Pool {
	return &Pool{
		ctx:      ctx,
		selector: selector,
		interval: time.Second * time.Duration(interval),
		timeout:  time.Second * time.Duration(timeout),
	}
}

func (p *Pool) AddBackends(backendCfg []config.BackendConf) {
	var backends []*backend.Backend
	for _, cfg := range backendCfg {
		parsedUrl, err := url.Parse(cfg.URL)
		if err != nil {
			slog.Error("URL parsing", "url", parsedUrl, "err", err)
			continue
		}
		b := backend.New(parsedUrl, int64(cfg.Weight))
		go b.HealthCheck(p.timeout)
		backends = append(backends, b)
	}

	p.backends = backends
	p.selector.Init(backends)
}

func (p *Pool) GetNext(ctx context.Context) *backend.Backend {
	return p.selector.Next(ctx)
}

func (p *Pool) HealthCheckStart() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, b := range p.backends {
				go b.HealthCheck(p.timeout)
			}
		case <-p.ctx.Done():
			return
		}
	}
}
