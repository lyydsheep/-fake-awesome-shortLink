package ioc

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
)

func InitSnowFlakeNode() *snowflake.Node {
	node, err := snowflake.NewNode(rand.Int63() % 1024)
	if err != nil {
		panic(err)
	}
	return node
}
