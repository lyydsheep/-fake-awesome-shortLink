package example

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	res, err := rdb.BFReserve(ctx, "hello:Hello", 0.001, 1000000).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("res", res)
	ok, err := rdb.BFExists(ctx, "bike:models", "Mountain").Result()
	if !ok {
		fmt.Println("not exist")
	}
	rdb.BFMAdd(ctx, "hello:Hello", "a", "b", "c")
	if rdb.BFExists(ctx, "hello:Hello", "a").Val() {
		fmt.Println("a!")
	}
}
