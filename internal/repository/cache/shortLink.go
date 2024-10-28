package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed check.lua
var check string
var ErrNotExistsInFilter = errors.New("NOT EXISTS")

type ShortLinkCache interface {
	Set(ctx context.Context, key string, val any) error
	Get(ctx context.Context, key string) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Check(ctx context.Context, filterKey, val, shortKey string) (string, error)
}

type ShortLinkRedis struct {
	rdb    redis.Cmdable
	expire time.Duration
}

func (cache *ShortLinkRedis) Check(ctx context.Context, filterKey, val, shortKey string) (string, error) {
	resAny, err := cache.rdb.Eval(ctx, check, []string{filterKey, shortKey}, []string{val}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", ErrNotExistsInFilter
	}
	return resAny.(string), nil
}

func (cache *ShortLinkRedis) Incr(ctx context.Context, key string) (int64, error) {
	return cache.rdb.Incr(ctx, key).Result()
}

func (cache *ShortLinkRedis) Set(ctx context.Context, key string, val any) error {
	return cache.rdb.Set(ctx, key, val, cache.expire).Err()
}

func (cache *ShortLinkRedis) Get(ctx context.Context, key string) (string, error) {
	return cache.rdb.Get(ctx, key).Result()
}

func NewShortLinkRedisV1(rdb redis.Cmdable) ShortLinkCache {
	return &ShortLinkRedis{
		rdb:    rdb,
		expire: time.Minute * 10,
	}
}
