package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	config := Config{
		Addr: "localhost:6379",
	}
	err := viper.UnmarshalKey("redis", &config)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})
}
