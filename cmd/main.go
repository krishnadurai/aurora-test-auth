package main

import (
	"context"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/interrupt"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
)

func main() {
	ctx, done := interrupt.Context()
	defer done()

	if err := realMain(ctx); err != nil {
		logger := logging.FromContext(ctx)
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
	var config cache.Config
	cacheDB, err := cache.NewRedisCache(ctx, config.Cache())
	if err != nil {
		return err
	}
	testRedis(ctx, cacheDB)
	return nil
}
