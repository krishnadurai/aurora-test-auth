package authcodes

import (
	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/secrets"
	"github.com/krishnadurai/aurora-test-auth/internal/setup"
)

// Compile-time check to assert this config matches requirements.
var _ setup.CacheConfigProvider = (*Config)(nil)
var _ setup.SecretManagerConfigProvider = (*Config)(nil)

// Config represents the configuration and associated environment variables for
// the authcodes components.
type Config struct {
	Cache         cache.Config
	SecretManager secrets.Config

	Port string `env:"PORT, default=8080"`
}

func (c *Config) CacheConfig() *cache.Config {
	return &c.Cache
}

func (c *Config) SecretManagerConfig() *secrets.Config {
	return &c.SecretManager
}
