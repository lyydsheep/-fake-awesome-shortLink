package example

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestDistribute(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	time.Sleep(time.Second)
	ctx := context.Background()
	for !rdb.SetNX(ctx, "1", 1, 10*time.Second).Val() {
	}
	for range 3 {
		time.Sleep(time.Second)
		fmt.Println("this is func2")
	}
	rdb.Del(ctx, "1")
}
