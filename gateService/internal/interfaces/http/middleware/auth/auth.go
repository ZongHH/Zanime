package auth

import (
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/pkg/errors"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(jwtManager *auth.JWTManager, cookieManager *auth.CookieManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := cookieManager.GetTokenCookie(c)
		if err != nil {
			c.Error(errors.NewAppError(errors.ErrUnauthorized.Code, err.Error(), err))
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			c.Error(errors.NewAppError(errors.ErrTokenInvalid.Code, err.Error(), err))
			c.Abort()
			return
		}

		c.Set("UserInfo", claims)
		c.Next()
	}
}
