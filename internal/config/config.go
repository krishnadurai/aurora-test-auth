package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/cache/redis"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cache cache.Cache
}

func LoadConfig(ctx context.Context, path string) (Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}

	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return Config{}, err
	}

	raw := struct {
		Cache yaml.Node `yaml:"cache"`
	}{}
	err = yaml.Unmarshal(b, &raw)
	if err != nil {
		return Config{}, err
	}

	cfgTypes := struct {
		Cache struct {
			Type string `yaml:"type"`
		}
	}{}
	err = yaml.Unmarshal(b, &cfgTypes)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	cfg.Cache, err = DecodeCache(ctx, cfgTypes.Cache.Type, raw.Cache)
	if err != nil {
		return Config{}, nil
	}

	return cfg, nil
}

func DecodeCache(ctx context.Context, typeName string, node yaml.Node) (cache.Cache, error) {
	switch typeName {
	case "redis":
		return redis.New(ctx, node)
	}

	return nil, fmt.Errorf("unsupported cache type: %s", typeName)
}
