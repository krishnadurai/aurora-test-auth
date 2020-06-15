package cache

import (
	"time"
)

//Config is a configuration spec for Cache
type Config struct {
	User               string        `envconfig:"CACHE_USER"`
	Addr               string        `envconfig:"CACHE_ADDR" default:"localhost:6379"`
	Password           string        `envconfig:"CACHE_PASSWORD"`
	DB                 int           `envconfig:"CACHE_DB"`
	PoolSize           int           `envconfig:"CACHE_POOL_SIZE"`
	MinIdleConns       int           `envconfig:"CACHE_MIN_IDLE_CONNS"`
	MaxConnAge         time.Duration `envconfig:"CACHE_MAX_CONN_AGE"`
	PoolTimeout        time.Duration `envconfig:"CACHE_POOL_TIMEOUT"`
	IdleTimeout        time.Duration `envconfig:"CACHE_IDLE_TIMEOUT"`
	IdleCheckFrequency time.Duration `envconfig:"CACHE_IDLE_CHECK_FREQ"`
}

//Cache returns a Config for Cache
func (c *Config) Cache() *Config {
	return c
}
