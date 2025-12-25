package algo

import (
	"context"

	"github.com/shv-ng/relay/internal/backend"
)

type weightedRoundRobin struct {
	backends []*backend.Backend

	totalWeight int64
}

func NewWeightedRoundRobin() Picker {
	return &weightedRoundRobin{}
}

func (r *weightedRoundRobin) Init(backends []*backend.Backend) {
	r.backends = backends
	var total int64
	for _, b := range r.backends {
		total += b.Weight
	}
	r.totalWeight = total
}

func (r *weightedRoundRobin) Next(ctx context.Context) *backend.Backend {
	if r.totalWeight == 0 {
		return nil
	}
	idx := -1
	var maxWeight int64
	for i, b := range r.backends {
		w := b.EffectiveWeight.Add(int64(b.Weight))
		if b.Alive() && w > int64(maxWeight) {
			maxWeight = w
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}
	r.backends[idx].EffectiveWeight.Add(-r.totalWeight)
	return r.backends[idx]
}
