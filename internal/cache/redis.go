package cache

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
)

// Compile-time check to verify implements interface.
var _ Cache = (*RedisCache)(nil)

// RedisCache implements Cache.
type RedisCache struct {
	Client *redis.Client
}

// NewRedisCache creates a cache client for Redis
func NewRedisCache(ctx context.Context, config *Config) (Cache, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("Creating cache connection pool.")

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	cache := &RedisCache{
		Client: client,
	}

	return cache, nil
}

// Set key - value pairs in Redis Cache
func (cache *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	statusCmd := cache.Client.Set(ctx, key, value, expiration)
	if statusCmd.Err() != nil {
		return "", statusCmd.Err()
	}
	return statusCmd.String(), nil
}

// Get value for key in Redis Cache
func (cache *RedisCache) Get(ctx context.Context, key string) (string, error) {
	statusCmd := cache.Client.Get(ctx, key)
	if statusCmd.Err() != nil {
		return "", statusCmd.Err()
	}
	return statusCmd.String(), nil
}

// Close releases cache connections.
func (cache *RedisCache) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Infof("Closing connection pool.")
	cache.Client.Close()
}
