// package service 提供了与进度相关的业务逻辑服务
package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// ProgressService 定义了进度服务的接口
// 提供了获取、保存、加载进度以及发布进度消息到NSQ的功能
type ProgressService interface {
	// GetWatchHistory 获取用户观看历史记录
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取观看历史的请求参数,包含用户ID、页码和每页数量
	// 返回:
	// - *dto.WatchHistoryResponse: 观看历史响应数据,包含进度列表
	// - error: 获取过程中的错误信息
	GetWatchHistory(ctx context.Context, request *dto.WatchHistoryRequest) (*dto.WatchHistoryResponse, error)

	// SaveProgress 保存用户进度
	// 参数:
	// - ctx: 上下文信息
	// - request: 保存进度的请求数据,包含用户ID和进度信息
	// 返回:
	// - *dto.SaveProgressResponse: 保存进度的响应数据
	// - error: 保存过程中的错误信息
	SaveProgress(ctx context.Context, request *dto.SaveProgressRequest) (*dto.SaveProgressResponse, error)

	// LoadProgress 加载用户进度
	// 参数:
	// - ctx: 上下文信息
	// - request: 加载进度的请求数据,包含用户ID等信息
	// 返回:
	// - *dto.LoadProgressResponse: 加载进度的响应数据
	// - error: 加载过程中的错误信息
	LoadProgress(ctx context.Context, request *dto.LoadProgressRequest) (*dto.LoadProgressResponse, error)

	// PublishTONsq 将进度信息发布到NSQ消息队列
	// 参数:
	// - ctx: 上下文信息
	// - request: 要发布的进度数据
	// 返回:
	// - error: 发布过程中的错误信息
	PublishTONsq(context.Context, *dto.SaveProgressRequest) error
}
