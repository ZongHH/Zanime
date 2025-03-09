package service

import (
	"context"
	"managerService/internal/interfaces/dto/request"
	"managerService/internal/interfaces/dto/response"
)

type UserActionLogsService interface {
	// GetUserActionLogs 获取用户行为日志
	// ctx 上下文
	// request 请求参数
	// 返回用户行为日志和可能的错误
	GetUserActionLogs(ctx context.Context, request *request.UserActionLogsRequest) (*response.UserActionLogsResponse, error)
}
