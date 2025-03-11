package security

import (
	"gateService/pkg/errors"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(tokenBucket *TokenBucket) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tokenBucket.GetToken(ctx)
		if err != nil {
			ctx.Error(errors.NewAppError(errors.ErrForbidden.Code, err.Error(), err))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
