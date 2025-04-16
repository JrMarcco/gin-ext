package ginext

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func W(fn func(ctx *gin.Context) (Res, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)

		// TODO: error log, need to consider the error type
		if err != nil {
			ctx.PureJSON(http.StatusInternalServerError, res)
			return
		}

		ctx.PureJSON(http.StatusOK, res)
	}
}

func B[T any](fn func(ctx *gin.Context, req T) (Res, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			// TODO: error log
			return
		}

		res, err := fn(ctx, req)
		// TODO: error log, need to consider the error type
		if err != nil {
			ctx.PureJSON(http.StatusInternalServerError, res)
			return
		}

		ctx.PureJSON(http.StatusOK, res)
	}
}
