//go:build wireinject

package startup

import (
	"awesome-shortLink/internal/repository"
	"awesome-shortLink/internal/repository/cache"
	"awesome-shortLink/internal/repository/dao"
	"awesome-shortLink/internal/repository/filter"
	"awesome-shortLink/internal/service"
	"awesome-shortLink/internal/web"
	"awesome-shortLink/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitShortLinkHdl() *web.ShortLinkHandler {
	wire.Build(web.NewShortLinkHandler, service.NewShortLinkServiceBasic,
		repository.NewShortLinkRepositoryV2, filter.NewBloomFilterV1, cache.NewShortLinkRedisV1,
		dao.NewShortLinkDAOV1, ioc.InitDB, ioc.InitZap, ioc.InitRedis)
	return new(web.ShortLinkHandler)
}

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitGinEngine, InitShortLinkHdl, ioc.InitMiddleWare)
	return new(gin.Engine)
}
