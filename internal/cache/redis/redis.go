package redis

import (
	"context"
	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"gopkg.in/yaml.v3"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
)

type Config struct {
	Addr     string `yaml:"addr" default:"localhost:6379"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db" default:"0"`
}

// RedisCache implements Cache.
type Cache struct {
	Client *redis.Client
}

// Verify at compile-time that RedisCache implements interface.
var _ cache.Cache = (*Cache)(nil)

// Set key - value pairs in Redis Cache
func (c *Cache) Set(ctx context.Context, key string, value string, expiration time.Duration) (string, error) {
	statusCmd := c.Client.Set(ctx, key, value, expiration)
	if statusCmd.Err() != nil {
		return "", statusCmd.Err()
	}
	return statusCmd.String(), nil
}

// Get value for key in Redis Cache
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	statusCmd := c.Client.Get(ctx, key)
	if statusCmd.Err() != nil {
		return "", statusCmd.Err()
	}
	return statusCmd.String(), nil
}

// Close releases cache connections.
func (c *Cache) Close(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Infof("Closing connection pool.")
	return c.Client.Close()
}

// New creates a cache client for Redis
func New(ctx context.Context, node yaml.Node) (*Cache, error) {
	logger := logging.FromContext(ctx)
	logger.Infof("Creating cache connection pool.")

	var cfg Config
	err := node.Decode(&cfg)
	if err != nil {
		return &Cache{}, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	cache := &Cache{
		Client: client,
	}

	return cache, nil
}
