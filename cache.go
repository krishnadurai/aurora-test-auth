package aurora_test_auth

import (
	"context"
	"time"
)

// Cache is an interface for an in-memory DB
type Cache interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Close(ctx context.Context) error
}

