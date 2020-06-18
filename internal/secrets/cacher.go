package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/google/exposure-notifications-server/pkg/cache"
)

// Compile-time check to verify implements interface.
var _ SecretManager = (*Cacher)(nil)

// Cacher is a secret manager implementation that wraps another secret manager
// and caches secret values.
type Cacher struct {
	sm    SecretManager
	cache *cache.Cache
}

// NewCacher creates a new secret manager that caches results for the given ttl.
func NewCacher(ctx context.Context, f SecretManagerFunc, ttl time.Duration) (SecretManager, error) {
	sm, err := f(ctx)
	if err != nil {
		return nil, fmt.Errorf("cacher: %w", err)
	}

	return WrapCacher(ctx, sm, ttl)
}

// WrapCacher wraps an existing SecretManager with caching.
func WrapCacher(ctx context.Context, sm SecretManager, ttl time.Duration) (SecretManager, error) {
	cache, err := cache.New(ttl)
	if err != nil {
		return nil, err
	}
	return &Cacher{
		sm:    sm,
		cache: cache,
	}, nil
}

// GetSecretValue implements the SecretManager interface, but caches values and
// retrieves them from the cache.
func (sm *Cacher) GetSecretValue(ctx context.Context, name string) (string, error) {
	lookup := func() (interface{}, error) {
		// Delegate lookup to parent sm.
		plaintext, err := sm.sm.GetSecretValue(ctx, name)
		if err != nil {
			return "", err
		}
		return plaintext, nil
	}

	cacheVal, err := sm.cache.WriteThruLookup(name, lookup)
	if err != nil {
		return "", nil
	}

	plaintext := cacheVal.(string)
	return plaintext, err
}
