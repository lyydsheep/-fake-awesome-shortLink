package example

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

//go:embed check.lua
var check string

func TestLua(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	res, err := rdb.Eval(context.Background(), check, []string{"filter", "a"}, "a").Result()
	fmt.Println(res, err)
}
