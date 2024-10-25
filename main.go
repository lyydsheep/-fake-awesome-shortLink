package main

import (
	"awesome-shortLink/ioc"
)

func init() {
	ioc.InitViper()
}

func main() {
	server := InitWebServer()
	server.Run(":8080")
}
