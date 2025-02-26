package security

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MiddleHandler struct {
	rdb *redis.Client
}

func NewMiddleHandler(rdb *redis.Client) *MiddleHandler {
	return &MiddleHandler{rdb: rdb}
}

// XSS防护中间件
func (m *MiddleHandler) XssProtection() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		// 获取当前的查询参数
		query := ctx.Request.URL.Query()
		// 创建新的查询参数map
		for k, v := range query {
			if len(v) > 0 {
				query.Set(k, htmlEscape(v[0]))
			}
		}
		// 更新URL的RawQuery
		ctx.Request.URL.RawQuery = query.Encode()
		ctx.Next()
	}
}

// CSRF防护中间件
func (m *MiddleHandler) CsrfProtection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			token := ctx.GetHeader("X-CSRF-Token")

			cookieToken, _ := ctx.Cookie("csrf_token")

			if token == "" || token != cookieToken {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "Invalid CSRF token",
				})
				return
			}
		}

		// 生成新的CSRF token
		token := generateToken()
		ctx.SetCookie("csrf_token", token, 3600, "/", "", false, true)
		ctx.Next()
	}
}

// 请求大小限制中间件
func (m *MiddleHandler) SizeLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10<<20) // 10MB
		ctx.Next()
	}
}

// 超时控制中间件
func (m *MiddleHandler) TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)

		defer cancel()

		ctx.Request = ctx.Request.WithContext(timeoutCtx)

		done := make(chan bool)
		go func() {
			ctx.Next()
			done <- true
		}()

		select {
		case <-done:
			return
		case <-timeoutCtx.Done():
			ctx.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"message": "Request timeout",
			})
			return
		}
	}
}

// 敏感信息脱敏中间件
func (m *MiddleHandler) SensitiveFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// 响应体脱敏
		if data, exists := ctx.Get("response"); exists {
			if str, ok := data.(string); ok {
				// 手机号脱敏
				str = regexp.MustCompile(`1[3-9]\d{9}`).ReplaceAllString(str, "1****")
				// 邮箱脱敏
				str = regexp.MustCompile(`[\w\.-]+@[\w\.-]+\.\w+`).ReplaceAllStringFunc(str, maskEmail)
				ctx.Set("response", str)
			}
		}
	}
}

func (m *MiddleHandler) IpLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		rateKey := fmt.Sprintf("rate_limit:%s", ip)
		banKey := fmt.Sprintf("ban_ip:%s", ip)
		banCountKey := fmt.Sprintf("ban_count:%s", ip)

		// 检查 IP 是否被封禁
		banned, err := m.rdb.Get(context.Background(), banKey).Int()
		if err == nil && banned == 1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Your IP has been banned due to too many requests",
			})
			return
		}

		// 使用 Redis MULTI/EXEC 保证原子性
		pipe := m.rdb.Pipeline()
		incr := pipe.Incr(context.Background(), rateKey)
		pipe.Expire(context.Background(), rateKey, time.Second)
		_, err = pipe.Exec(context.Background())

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Rate limit error",
			})
			return
		}

		// 获取当前计数
		count := incr.Val()

		// 如果请求次数超过100次
		if count > 100 {
			// 增加封禁计数
			banCount, _ := m.rdb.Incr(context.Background(), banCountKey).Result()
			m.rdb.Expire(context.Background(), banCountKey, time.Hour*24)

			// 如果在24小时内触发限制3次以上，封禁IP 12小时
			if banCount >= 3 {
				m.rdb.Set(context.Background(), banKey, 1, time.Hour*12)
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "Your IP has been banned for 12 hours due to repeated violations",
				})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		ctx.Next()
	}
}

// 辅助函数
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "/", "&#47;")
	return s
}

func generateToken() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(hash.Sum(nil))
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	username := parts[0]
	if len(username) > 2 {
		username = username[:2] + "****"
	}
	return username + "@" + parts[1]
}
