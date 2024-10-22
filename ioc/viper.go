package ioc

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}
