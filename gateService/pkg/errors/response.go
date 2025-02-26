package errors

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, err error) {
	if e, ok := err.(*AppError); ok {
		c.JSON(200, Response{
			Code:    e.Code,
			Message: e.Message,
		})
		return
	}

	c.JSON(200, Response{
		Code:    ErrCodeInternalError,
		Message: err.Error(),
	})
}

// FailWithCode 指定错误码的失败响应
func FailWithCode(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}
