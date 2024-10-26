package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type ShortLinkCache interface {
	BFAdd(ctx context.Context, key string, val any) error
	BFExists(ctx context.Context, key string, val any) bool
	Set(ctx context.Context, key string, val any) error
	Get(ctx context.Context, key string) (string, error)
}

type ShortLinkRedis struct {
	rdb    redis.Cmdable
	expire time.Duration
}

func (cache *ShortLinkRedis) BFAdd(ctx context.Context, key string, val any) error {
	return cache.rdb.BFAdd(ctx, key, val).Err()
}

func (cache *ShortLinkRedis) BFExists(ctx context.Context, key string, val any) bool {
	return cache.rdb.BFExists(ctx, key, val).Val()
}

func (cache *ShortLinkRedis) Set(ctx context.Context, key string, val any) error {
	return cache.rdb.Set(ctx, key, val, cache.expire).Err()
}

func (cache *ShortLinkRedis) Get(ctx context.Context, key string) (string, error) {
	return cache.rdb.Get(ctx, key).Result()
}

func NewShortLinkRedis(rdb redis.Cmdable) ShortLinkCache {
	return &ShortLinkRedis{
		rdb: rdb,
	}
}
