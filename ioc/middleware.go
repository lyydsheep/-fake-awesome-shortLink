package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitMiddleWare() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORS(),
	}
}

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http:localhost") {
				return true
			}
			return false
		},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		MaxAge:           time.Hour * 12,
		ExposeHeaders:    []string{"Location"},
	})
}
