package ioc

import (
	"awesome-shortLink/internal/repository/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	config := Config{
		DSN: "root:root@tcp(127.0.0.1:13366)/shortlink",
	}
	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		fmt.Println(err)
	}
	db, err := gorm.Open(mysql.Open(config.DSN))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&dao.ShortLink{})
	if err != nil {
		panic(err)
	}
	return db
}
