package main

import (
	"context"
	"fmt"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
)

func testRedis(ctx context.Context, cacheDB cache.Cache) {
	setResult, err := cacheDB.Set(ctx, "test", "test", 1000000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(setResult)
	cacheDB.Get(ctx, "test")
	fmt.Println("Redis")
}

func main() {
	ctx := context.Background()
	var config cache.Config
	cacheDB, _ := cache.NewRedisCache(ctx, config.Cache())
	testRedis(ctx, cacheDB)
}
