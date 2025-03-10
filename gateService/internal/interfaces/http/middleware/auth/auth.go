package auth

import (
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/pkg/errors"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
// 参数:
// - jwtManager: JWT令牌管理器,用于解析和验证令牌
// - cookieManager: Cookie管理器,用于获取Cookie中的令牌
// 返回:
// - gin.HandlerFunc: Gin中间件处理函数
func JWTAuthMiddleware(jwtManager *auth.JWTManager, cookieManager *auth.CookieManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie中获取令牌
		token, err := cookieManager.GetTokenCookie(c)
		if err != nil {
			// 获取令牌失败,返回未授权错误
			c.Error(errors.NewAppError(errors.ErrUnauthorized.Code, err.Error(), err))
			c.Abort()
			return
		}

		// 解析令牌获取声明信息
		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			// 令牌解析失败,返回令牌无效错误
			c.Error(errors.NewAppError(errors.ErrTokenInvalid.Code, err.Error(), err))
			c.Abort()
			return
		}

		// 认证通过,将用户信息存入上下文
		c.Set("UserInfo", claims)
		c.Next()
	}
}
