package algo

import (
	"context"

	"github.com/shv-ng/relay/internal/backend"
)

type ContextKey string

const ClientIPKey ContextKey = "clientIP"

type Picker interface {
	Init([]*backend.Backend)
	Next(context.Context) *backend.Backend
}
