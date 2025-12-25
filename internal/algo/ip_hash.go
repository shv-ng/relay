package algo

import (
	"context"
	"hash/fnv"

	"github.com/shv-ng/relay/internal/backend"
)

type ipHash struct {
	backends []*backend.Backend
	total    uint64
}

func NewIPHash() Picker {
	return &ipHash{}
}

func (i *ipHash) Init(backends []*backend.Backend) {
	i.backends = backends
	i.total = uint64(len(backends))
}

func (i *ipHash) Next(ctx context.Context) *backend.Backend {
	if len(i.backends) == 0 {
		return nil
	}

	clientIP, ok := ctx.Value(ClientIPKey).(string)
	// no ip found, return whoever is alive
	if !ok || clientIP == "" {
		for _, b := range i.backends {
			if b.Alive() {
				return b
			}
		}
		// all dead
		return nil
	}

	hash := getHash(clientIP)
	startIdx := hash % i.total

	// find alive from and after idx
	for j := range i.total {
		idx := (startIdx + j) % i.total

		if i.backends[idx].Alive() {
			return i.backends[idx]
		}
	}

	return nil
}

func getHash(ip string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(ip))
	return hash.Sum64()
}
