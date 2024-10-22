package example

import (
	"awesome-shortLink/dao"
	"awesome-shortLink/ioc"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func init() {
	ioc.InitViper()
}

func TestDB(t *testing.T) {
	dsnVal := viper.Get("db.dsn")
	dsn := dsnVal.(string)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&dao.ShortLink{})
	if err != nil {
		fmt.Println(err)
	}
}
