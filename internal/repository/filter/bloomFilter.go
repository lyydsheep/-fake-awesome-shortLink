package filter

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	bloomFilterKey = "bloomFilter:shortURL"
)

type BloomFilter interface {
	BFAdd(ctx context.Context, key string, val any) error
	BFExists(ctx context.Context, key string, val any) bool
}

type BloomFilterV1 struct {
	rdb redis.Cmdable
}

func (bf *BloomFilterV1) BFAdd(ctx context.Context, key string, val any) error {
	return bf.rdb.BFAdd(ctx, key, val).Err()
}

func (bf *BloomFilterV1) BFExists(ctx context.Context, key string, val any) bool {
	return bf.rdb.BFExists(ctx, key, val).Val()
}

func NewBloomFilterV1(rdb redis.Cmdable) BloomFilter {
	rdb.BFReserve(context.Background(), bloomFilterKey, 0.001, 100000)
	return &BloomFilterV1{
		rdb: rdb,
	}
}
