package example

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	node.Generate()
	id := node.Generate()
	fmt.Printf("Base2 ID:%s\n", id.Base2())
	fmt.Println(id)
}
