package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"gateService/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AppError 自定义错误结构
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // 原始错误，不返回给前端
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
}

// NewAppError 创建自定义错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ErrorResponse 统一的错误响应结构
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 增加更多错误码
const (
	ErrCodeSuccess        = 0
	ErrCodeParamInvalid   = 400
	ErrCodeUnauthorized   = 401
	ErrCodeTokenInvalid   = 402
	ErrCodeForbidden      = 403
	ErrCodeNotFound       = 404
	ErrCodeInternalError  = 500
	ErrCodeBusinessError  = 1001
	ErrCodeDatabaseError  = 1002
	ErrCodeValidationFail = 1003

	// 业务错误码 (2000-2999)
	ErrCodeUserNotFound  = 2001
	ErrCodePasswordWrong = 2002
	ErrCodeTokenExpired  = 2003

	// 第三方服务错误 (3000-3999)
	ErrCodeThirdPartyError = 3001
	ErrCodeRPCError        = 3002

	// 系统错误 (5000-5999)
	ErrCodeConfigError  = 5001
	ErrCodeNetworkError = 5002
)

// ErrorConfig 错误处理配置
type ErrorConfig struct {
	Env            string        // 环境：development/production
	AlertThreshold int           // 错误告警阈值
	AlertInterval  time.Duration // 告警间隔
}

var (
	errorConfig = ErrorConfig{
		Env:            "development",
		AlertThreshold: 10,
		AlertInterval:  time.Minute * 5,
	}
	errorCount = make(map[int]int) // 记录错误次数
	lastAlert  = time.Now()        // 上次告警时间
)

// SetErrorConfig 设置错误处理配置
func SetErrorConfig(config ErrorConfig) {
	errorConfig = config
}

// logError 记录错误日志
func logError(err error, stack string, c *gin.Context) {
	fields := []zap.Field{
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
	}

	if appErr, ok := err.(*AppError); ok {
		fields = append(fields,
			zap.Int("error_code", appErr.Code),
			zap.Error(appErr.Err),
		)

		// 记录错误次数
		errorCount[appErr.Code]++

		// 检查是否需要告警
		if errorCount[appErr.Code] >= errorConfig.AlertThreshold &&
			time.Since(lastAlert) > errorConfig.AlertInterval {
			sendAlert(appErr)
			lastAlert = time.Now()
		}
	}

	if stack != "" {
		fields = append(fields, zap.String("stack", stack))
	}

	logger.Log.Error("request error", fields...)
}

// sendAlert 发送告警
func sendAlert(err *AppError) {
	// TODO: 实现告警逻辑，可以是发送邮件、钉钉、短信等
	logger.Log.Warn("error alert",
		zap.Int("error_code", err.Code),
		zap.String("message", err.Message),
		zap.Int("count", errorCount[err.Code]),
	)
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer确保panic后的错误处理逻辑一定会执行
		defer func() {
			// recover()捕获panic抛出的错误
			if err := recover(); err != nil {
				// 获取完整的堆栈信息用于调试
				stack := string(debug.Stack())
				// 记录错误和堆栈信息到日志系统
				logError(fmt.Errorf("%v", err), stack, c)

				// 构造返回给客户端的错误响应
				var response ErrorResponse

				// 根据环境配置返回不同详细程度的错误信息
				if errorConfig.Env == "development" {
					// 开发环境下返回详细错误信息和堆栈信息
					response = ErrorResponse{
						Code:    500,                    // 服务器内部错误状态码
						Message: fmt.Sprintf("%v", err), // 原始错误信息
						Data:    stack,                  // 堆栈信息用于调试
					}
				} else {
					// 生产环境下只返回简单的错误提示
					response = ErrorResponse{
						Code:    500,
						Message: "服务器内部错误", // 对用户友好的错误提示
					}
				}

				// 返回500状态码和错误响应
				c.JSON(http.StatusInternalServerError, response)
				// 终止后续中间件的执行
				c.Abort()
			}
		}()

		c.Next()

		// 处理请求过程中的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logError(err, "", c)

			switch e := err.(type) {
			case *AppError:
				response := ErrorResponse{
					Code:    e.Code,
					Message: e.Message,
				}

				// 在开发环境下，返回原始错误信息
				if errorConfig.Env == "development" && e.Err != nil {
					response.Data = e.Err.Error()
				}

				c.JSON(http.StatusOK, response)
			default:
				response := ErrorResponse{
					Code:    500,
					Message: "服务器内部错误",
				}

				if errorConfig.Env == "development" {
					response.Message = err.Error()
				}

				c.JSON(http.StatusInternalServerError, response)
			}
		}
	}
}

// 预定义一些常用错误
var (
	ErrParamInvalid    = NewAppError(ErrCodeParamInvalid, "参数无效", nil)
	ErrUnauthorized    = NewAppError(ErrCodeUnauthorized, "未授权访问", nil)
	ErrTokenInvalid    = NewAppError(ErrCodeTokenInvalid, "令牌无效", nil)
	ErrForbidden       = NewAppError(ErrCodeForbidden, "禁止访问", nil)
	ErrNotFound        = NewAppError(ErrCodeNotFound, "资源不存在", nil)
	ErrInternalError   = NewAppError(ErrCodeInternalError, "服务器内部错误", nil)
	ErrDatabaseError   = NewAppError(ErrCodeDatabaseError, "数据库操作失败", nil)
	ErrValidationFail  = NewAppError(ErrCodeValidationFail, "数据验证失败", nil)
	ErrThirdPartyError = NewAppError(ErrCodeThirdPartyError, "第三方服务错误", nil)
	ErrRPCError        = NewAppError(ErrCodeRPCError, "RPC服务错误", nil)
	ErrConfigError     = NewAppError(ErrCodeConfigError, "配置错误", nil)
	ErrNetworkError    = NewAppError(ErrCodeNetworkError, "网络错误", nil)
)
