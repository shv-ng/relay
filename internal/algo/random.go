package algo

import (
	"context"
	"math/rand"

	"github.com/shv-ng/relay/internal/backend"
)

type random struct {
	backends []*backend.Backend

	total int
}

func NewRandom() Picker {
	return &random{}
}

func (r *random) Init(backends []*backend.Backend) {
	r.backends = backends
	r.total = len(backends)
}

func (r *random) Next(ctx context.Context) *backend.Backend {
	if r.total == 0 {
		return nil
	}

	randomIdx := rand.Intn(r.total)
	for i := range r.total {
		idx := (randomIdx + i) % r.total

		if r.backends[idx].Alive() {
			return r.backends[idx]
		}
	}
	return nil
}
