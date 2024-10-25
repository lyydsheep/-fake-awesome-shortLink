package ioc

import "go.uber.org/zap"

func InitZap() *zap.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return l
}
