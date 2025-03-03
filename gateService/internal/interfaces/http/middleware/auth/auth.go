package auth

import (
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
// 参数:
// - jwtManager: JWT令牌管理器,用于解析和验证令牌
// - cookieManager: Cookie管理器,用于获取Cookie中的令牌
// - userRepository: 用户仓储接口,用于检查用户状态
// 返回:
// - gin.HandlerFunc: Gin中间件处理函数
func JWTAuthMiddleware(jwtManager *auth.JWTManager, cookieManager *auth.CookieManager, userRepository repository.UserRepository) gin.HandlerFunc {
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

		// 检查体验用户账号是否已失效
		exist, _, err := userRepository.CheckInRedis(c, "test_account:deleted:"+strconv.Itoa(claims.UserInfo.UserID))
		if err != nil {
			// 检查失败,返回未授权错误
			c.Error(errors.NewAppError(errors.ErrUnauthorized.Code, err.Error(), err))
			c.Abort()
			return
		}
		if exist {
			// 体验用户账号已失效,返回未授权错误
			c.Error(errors.NewAppError(errors.ErrUnauthorized.Code, "体验用户账号已失效", nil))
			c.Abort()
			return
		}

		// 认证通过,将用户信息存入上下文
		c.Set("UserInfo", claims)
		c.Next()
	}
}
