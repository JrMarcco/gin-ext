package ratelimit

import (
	"log"
	"net/http"
	"strings"

	"github.com/JrMarcco/gin-ext/internal/ratelimit"
	"github.com/gin-gonic/gin"
)

type Builder struct {
	limiter  ratelimit.Limiter
	keyGenFn func(ctx *gin.Context) string
	logFn    func(msg any, arg ...any)
}

func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{
		limiter: limiter,
		keyGenFn: func(ctx *gin.Context) string {
			// default key generator, use ip limit
			var sb strings.Builder
			sb.WriteString("ip-limiter:")
			sb.WriteString(ctx.ClientIP())
			return sb.String()
		},
		logFn: func(msg any, args ...any) {
			// default log function, use println
			toPrint := make([]any, 0, len(args)+1)
			toPrint = append(toPrint, msg)
			toPrint = append(toPrint, args...)
			log.Println(toPrint...)
		},
	}
}

func (b *Builder) WithKeyGenFn(keyGenFn func(ctx *gin.Context) string) *Builder {
	b.keyGenFn = keyGenFn
	return b
}

func (b *Builder) WithLogFn(logFn func(msg any, arg ...any)) *Builder {
	b.logFn = logFn
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			b.logFn(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	return b.limiter.Limit(ctx, b.keyGenFn(ctx))
}
