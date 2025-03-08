package repository

import "context"

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
}
