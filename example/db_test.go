package example

import (
	"gorm.io/gorm"
	"testing"
)

func TestDB(t *testing.T) {
	db, _ := gorm.Open()
}
