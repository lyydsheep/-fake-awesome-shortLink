package web

import (
	"awesome-shortLink/ginx"
	"awesome-shortLink/service"
	"github.com/gin-gonic/gin"
)

type ShortLinkHandler struct {
	svc service.ShortLinkService
}

func NewShortLinkHandler() *ShortLinkHandler {
	return &ShortLinkHandler{}
}

func (hdl *ShortLinkHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/shorten", ginx.WrapBody[Req, string](hdl.Shorten))
}

func (hdl *ShortLinkHandler) Shorten(req Req) (ginx.Result[string], error) {
	res, err := hdl.svc.ShortenURL(req.URL)
	if err != nil {
		return ginx.Result[string]{}, err
	}
	return ginx.Result[string]{
		Data: res.Short,
		Msg:  "OK",
	}, nil
}

type Req struct {
	URL string `json:"url"`
}
