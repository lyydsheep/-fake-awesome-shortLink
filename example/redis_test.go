package example

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	res, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("yes")
		}
	} else {
		fmt.Println(res)
	}
}
