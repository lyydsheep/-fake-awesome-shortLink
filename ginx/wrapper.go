package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapBody[Req any, T any](fn func(ctx *gin.Context, req Req) (Result[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			fmt.Println(fmt.Errorf("wrapBody failed to bind, %v", err))
			ctx.JSON(http.StatusOK, Result[string]{
				Msg: "系统错误",
			})
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			fmt.Println(fmt.Errorf("wrapBody failed to fn, %v", err))
			ctx.JSON(http.StatusOK, Result[string]{
				Msg: "系统错误",
			})
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
