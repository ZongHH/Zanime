package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

// ProgressRepository 定义了观看进度数据访问层的接口
type ProgressRepository interface {
	// CreateProgress 创建新的观看进度记录
	// 参数:
	//   - ctx: 上下文信息
	//   - progress: 观看进度信息
	// 返回:
	//   - error: 可能的错误信息
	CreateProgress(ctx context.Context, progress *entity.Progress) error

	// UpdateProgress 更新观看进度记录
	// 参数:
	//   - ctx: 上下文信息
	//   - progress: 更新后的观看进度信息
	// 返回:
	//   - error: 可能的错误信息
	UpdateProgress(ctx context.Context, progress *entity.Progress) error

	// GetProgress 获取用户的所有观看进度记录
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - page: 页码
	//   - pageSize: 每页数量
	// 返回:
	//   - []*entity.Progress: 观看进度记录列表
	//   - error: 可能的错误信息
	GetProgress(ctx context.Context, userID int, page int, pageSize int) ([]*entity.Progress, error)

	// GetUserWatchProgress 获取用户特定视频的观看进度
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - videoID: 视频ID
	// 返回:
	//   - *entity.Progress: 观看进度记录
	//   - error: 可能的错误信息
	GetUserWatchProgress(ctx context.Context, userID, videoID int) (*entity.Progress, error)
}
