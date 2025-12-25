package algo

import (
	"sync/atomic"

	"github.com/shv-ng/relay/internal/backend"
)

type roundRobin struct {
	backends []*backend.Backend

	current atomic.Uint64
	total   uint64
}

func NewRoundRobin() Picker {
	return &roundRobin{}
}

func (r *roundRobin) Init(backends []*backend.Backend) {
	r.backends = backends
	r.total = uint64(len(backends))
}

func (r *roundRobin) Next() *backend.Backend {
	if r.total == 0 {
		return nil
	}
	val := r.current.Add(1)
	for i := range r.total {
		idx := (i + val) % r.total
		if r.backends[idx].Alive() {
			return r.backends[idx]
		}
	}
	return nil
}
