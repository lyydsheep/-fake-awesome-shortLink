package tools

import (
	"fmt"
	"testing"
)

func TestBase62(t *testing.T) {
	oldData := int64(123456)
	fmt.Println(Encode(oldData))
}
