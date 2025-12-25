package algo

import (
	"github.com/shv-ng/relay/internal/backend"
)

type Picker interface {
	Init([]*backend.Backend)
	Next() *backend.Backend
}
