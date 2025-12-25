package pool

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/shv-ng/relay/internal/algo"
	"github.com/shv-ng/relay/internal/backend"
)

type Pool struct {
	ctx context.Context

	backends []*backend.Backend
	selector algo.Picker
}

func New(ctx context.Context, selector algo.Picker, urls []string) *Pool {
	var backends []*backend.Backend
	for _, urlStr := range urls {
		parsedUrl, err := url.Parse(urlStr)
		if err != nil {
			slog.Error("URL parsing", "url", parsedUrl, "err", err)
			continue
		}
		b := backend.New(parsedUrl)
		go b.HealthCheck()
		backends = append(backends, b)
	}

	selector.Init(backends)

	return &Pool{
		ctx:      ctx,
		backends: backends,
		selector: selector,
	}
}

func (p *Pool) GetNext() *backend.Backend {
	return p.selector.Next()
}

func (p *Pool) HealthCheckStart() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, b := range p.backends {
				go b.HealthCheck()
			}
		case <-p.ctx.Done():
			ticker.Stop()
			return
		}
	}
}
