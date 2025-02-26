package repository

import (
	"context"
	"crawler/internal/domain/entity"
)

// VideoRepository 接口定义了与视频相关的数据库操作方法
type VideoRepository interface {
	// CreateVideo 方法用于在数据库中创建一个新的视频记录
	// 参数 ctx 是上下文，用于控制请求的生命周期
	// 参数 anime 是要创建的视频实体
	// 返回值是一个 error，表示操作是否成功
	CreateVideo(ctx context.Context, anime *entity.Video) error

	// CreateGenre 方法用于在数据库中创建多个视频类型记录
	// 参数 ctx 是上下文，用于控制请求的生命周期
	// 参数 genres 是要创建的动漫类型的切片
	// 返回值是一个 error，表示操作是否成功
	CreateGenre(ctx context.Context, genres *[]entity.AnimeGenre) error

	// CreateVideoUrl 方法用于在数据库中创建多个视频链接记录
	// 参数 ctx 是上下文，用于控制请求的生命周期
	// 参数 videoUrls 是要创建的视频链接的切片
	// 返回值是一个 error，表示操作是否成功
	CreateVideoUrl(ctx context.Context, videoUrls *[]entity.VideoUrl) error

	// GetVideoIDByVideoName 方法用于根据视频名称获取视频的 ID
	// 参数 ctx 是上下文，用于控制请求的生命周期
	// 参数 videoName 是要查询的视频名称
	// 返回值是视频的 ID 和可能发生的错误
	GetVideoIDByVideoName(ctx context.Context, videoName string) (int, error)

	// GetVideoURLByVideoIDANDEpisode 方法用于根据视频 ID 和集数获取视频的链接
	// 参数 ctx 是上下文，用于控制请求的生命周期
	// 参数 videoID 是要查询的视频 ID
	// 参数 Episode 是要查询的集数
	// 返回值是视频的链接和可能发生的错误
	GetVideoURLByVideoIDANDEpisode(ctx context.Context, videoID int, Episode string) (string, error)
}
