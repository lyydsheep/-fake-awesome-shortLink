package ioc

import (
	"awesome-shortLink/internal/web"
	"github.com/gin-gonic/gin"
)

func InitGinEngine(shortLinKHdl *web.ShortLinkHandler, hdls ...gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(hdls...)
	shortLinKHdl.RegisterRoutes(server)
	return server
}
