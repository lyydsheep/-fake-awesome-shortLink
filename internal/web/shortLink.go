package web

import (
	"awesome-shortLink/ginx"
	"awesome-shortLink/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrNotFound = service.ErrNotFound
)

type ShortLinkHandler struct {
	svc service.ShortLinkService
	l   *zap.Logger
}

func NewShortLinkHandler(svc service.ShortLinkService, l *zap.Logger) *ShortLinkHandler {
	return &ShortLinkHandler{
		svc: svc,
		l:   l,
	}
}

func (hdl *ShortLinkHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/shorten", ginx.WrapBody[Req, string](hdl.Shorten))
	server.GET("/sl/:shortURL", hdl.Obtain)
}

func (hdl *ShortLinkHandler) Shorten(ctx *gin.Context, req Req) (ginx.Result[string], error) {
	sl, err := hdl.svc.ShortenURL(ctx, req.URL)
	if err != nil {
		return ginx.Result[string]{}, err
	}
	return ginx.Result[string]{
		Data: fmt.Sprintf("http://localhost:8080/sl/%s", sl.Short),
		Msg:  "OK",
	}, nil
}

func (hdl *ShortLinkHandler) Obtain(ctx *gin.Context) {
	shortURL := ctx.Param("shortURL")
	sl, err := hdl.svc.Obtain(ctx, shortURL)
	switch {
	case errors.Is(err, ErrNotFound):
		hdl.l.Error("木有长链", zap.Error(err))
		ctx.JSON(http.StatusNotFound, ginx.Result[string]{
			Code: 4,
			Msg:  "非法短链",
		})
		return
	case err != nil:
		hdl.l.Error("获取长链失败", zap.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result[string]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.Redirect(http.StatusFound, sl.Long)
}

type Req struct {
	URL string `json:"url"`
}
