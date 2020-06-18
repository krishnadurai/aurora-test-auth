package main

import (
	"context"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/config"
	"github.com/krishnadurai/aurora-test-auth/internal/interrupt"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
)

var cfg config.Config

func main() {
	ctx, cancel := interrupt.Context()
	defer cancel()

	logger := logging.FromContext(ctx)

	var err error

	cfg, err = config.LoadConfig(ctx, "config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	err = realMain(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}

func testRedis(ctx context.Context, cacheDB cache.Cache) {
	logger := logging.FromContext(ctx)
	setResult, err := cacheDB.Set(ctx, "test", "test", 20000000)
	if err != nil {
		logger.Error("error")
	}
	logger.Info(setResult)

	getResult, err := cacheDB.Get(ctx, "test")
	if err != nil {
		logger.Error(err)
	}
	logger.Info(getResult)
}

func realMain(ctx context.Context) error {
	testRedis(ctx, cfg.Cache)
	return nil
}
