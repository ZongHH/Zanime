package repository

import (
	"context"
	"managerService/internal/domain/entity"
)

// StatisticsRepository 定义了统计数据仓储层的接口
// 提供了获取各类统计数据的方法
type StatisticsRepository interface {
	// GetUserCount 获取系统中的总用户数量
	// ctx 用于传递上下文信息
	// 返回用户总数和可能的错误
	GetUserCount(ctx context.Context) (int, error)

	// GetAnimeCount 获取系统中的动画总数量
	// ctx 用于传递上下文信息
	// 返回动画总数和可能的错误
	GetAnimeCount(ctx context.Context) (int, error)

	// GetTodayPlayCount 获取今日的播放总次数
	// ctx 用于传递上下文信息
	// 返回今日播放次数和可能的错误
	GetTodayPlayCount(ctx context.Context) (int, error)

	// GetActiveUserCount 获取活跃用户数量
	// ctx 用于传递上下文信息
	// 返回活跃用户数和可能的错误
	GetActiveUserCount(ctx context.Context) (int, error)

	// GetNewAnime 获取最新上线的动漫
	// ctx 用于传递上下文信息
	// page 页码
	// pageSize 每页数量
	// 返回最新上线的动漫和可能的错误
	GetNewAnime(ctx context.Context, page int, pageSize int) ([]*entity.Anime, error)
}

// UserActionLogsRepository 定义了用户行为日志仓储层的接口
// 提供了获取用户行为日志的方法
type UserActionLogsRepository interface {
	// GetUserActionLogs 获取用户行为日志
	// ctx 用于传递上下文信息
	// page 页码
	// pageSize 每页数量
	// 返回用户行为日志和可能的错误
	GetUserActionLogs(ctx context.Context, page int, pageSize int) ([]*entity.UserActionLog, error)
}
