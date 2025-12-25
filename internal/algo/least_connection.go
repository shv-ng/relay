package algo

import (
	"context"

	"github.com/shv-ng/relay/internal/backend"
)

type leastConnection struct {
	backends []*backend.Backend
}

func NewLeastConnection() Picker {
	return &leastConnection{}
}

func (i *leastConnection) Init(backends []*backend.Backend) {
	i.backends = backends
}

func (l *leastConnection) Next(ctx context.Context) *backend.Backend {
	if len(l.backends) == 0 {
		return nil
	}
	idx := -1
	var minConnection int64
	for i, b := range l.backends {
		if !b.Alive() {
			continue
		}
		conn := b.ActiveConnection.Load()
		if idx == -1 || conn < minConnection {
			minConnection = conn
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}
	return l.backends[idx]
}
