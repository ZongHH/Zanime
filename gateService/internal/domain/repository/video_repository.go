package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

// VideoRepository 定义了视频仓储的接口规范
type VideoRepository interface {
	// 数据库相关操作

	// CreateVideo 创建新的视频记录
	// 参数:
	//   - ctx: 上下文信息
	//   - video: 要创建的视频实体
	// 返回:
	//   - error: 可能的错误信息
	CreateVideo(ctx context.Context, video *entity.Video) error

	// UpdateVideo 更新视频信息
	// 参数:
	//   - ctx: 上下文信息
	//   - video: 包含更新信息的视频实体
	// 返回:
	//   - error: 可能的错误信息
	UpdateVideo(ctx context.Context, video *entity.Video) error

	// DeleteVideo 删除指定ID的视频
	// 参数:
	//   - ctx: 上下文信息
	//   - videoID: 要删除的视频ID
	// 返回:
	//   - error: 可能的错误信息
	DeleteVideo(ctx context.Context, videoID int) error

	// GetVideoByID 根据ID获取视频信息
	// 参数:
	//   - ctx: 上下文信息
	//   - videoID: 视频ID
	// 返回:
	//   - *entity.Video: 视频信息
	//   - error: 可能的错误信息
	GetVideoByID(ctx context.Context, videoID int) (*entity.Video, error)

	// GetVideosByVideoName 根据视频名称搜索视频
	// 参数:
	//   - ctx: 上下文信息
	//   - videoName: 视频名称
	// 返回:
	//   - []*entity.Video: 视频列表
	//   - error: 可能的错误信息
	GetVideosByVideoName(ctx context.Context, videoName string) ([]*entity.Video, error)

	// GetVideosALLEpisodesByVideoName 获取指定视频名称的所有剧集
	// 参数:
	//   - ctx: 上下文信息
	//   - videoName: 视频名称
	//   - page: 分页页码
	// 返回:
	//   - []*entity.Video: 视频剧集列表
	//   - error: 可能的错误信息
	GetVideosALLEpisodesByVideoName(ctx context.Context, videoName string, page int) ([]*entity.Video, error)

	// GetVideoInfoWithEposidesByVideoID 获取视频详细信息及其剧集
	// 参数:
	//   - ctx: 上下文信息
	//   - videoID: 视频ID
	// 返回:
	//   - *entity.Video: 视频详细信息及剧集
	//   - error: 可能的错误信息
	GetVideoInfoWithEposidesByVideoID(ctx context.Context, videoID int) (*entity.Video, error)

	// GetVideosByFilters 根据过滤条件获取视频列表
	// 参数:
	//   - ctx: 上下文信息
	//   - video: 包含过滤条件的视频实体
	//   - page: 分页页码
	//   - limit: 每页数量
	// 返回:
	//   - []*entity.Video: 视频列表
	//   - int: 总数量
	//   - error: 可能的错误信息
	GetVideosByFilters(ctx context.Context, video *entity.Video, page, limit int) ([]*entity.Video, int, error)

	// GetVideoFilters 获取视频过滤条件选项
	// 参数:
	//   - ctx: 上下文信息
	// 返回:
	//   - *entity.VideoFilters: 过滤条件选项
	//   - error: 可能的错误信息
	GetVideoFilters(ctx context.Context) (*entity.VideoFilters, error)

	// GetAnimesByGenre 根据动漫类型获取动漫列表
	// 参数:
	//   - ctx: 上下文信息
	//   - genre: 动漫类型
	//   - page: 分页页码
	//   - limit: 每页数量
	// 返回:
	//   - []*entity.Video: 动漫列表
	//   - error: 可能的错误信息
	GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*entity.Video, error)

	// GetTopAnimeGenres 获取热门动漫类型
	// 参数:
	//   - ctx: 上下文信息
	// 返回:
	//   - []string: 热门动漫类型
	//   - error: 可能的错误信息
	GetTopAnimeGenres(ctx context.Context) ([]string, error)

	// Redis相关操作

	// CacheVideoURL 缓存视频URL
	// 参数:
	//   - context: 上下文信息
	//   - string: 缓存key
	//   - string: 视频URL
	// 返回:
	//   - error: 可能的错误信息
	CacheVideoURL(context.Context, string, string) error

	// GetVideoURL 获取缓存的视频URL
	// 参数:
	//   - context: 上下文信息
	//   - string: 缓存key
	// 返回:
	//   - string: 视频URL
	//   - error: 可能的错误信息
	GetVideoURL(context.Context, string) (string, error)

	// AddAnimeCollection 添加用户动漫收藏记录
	// 参数:
	//   - ctx: 上下文信息
	//   - collection: 收藏记录实体
	// 返回:
	//   - error: 可能的错误信息
	AddAnimeCollection(ctx context.Context, collection *entity.UserAnimeCollection) error

	// DeleteAnimeCollection 删除用户动漫收藏记录
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - videoID: 视频ID
	// 返回:
	//   - error: 可能的错误信息
	DeleteAnimeCollection(ctx context.Context, userID, videoID int) error

	// GetAnimeCollectionByUser 获取用户的动漫收藏列表
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - page: 页码
	//   - limit: 每页数量
	// 返回:
	//   - []*entity.Video: 收藏的动漫列表
	//   - int: 总数量
	//   - error: 可能的错误信息
	GetAnimeCollectionByUser(ctx context.Context, userID int, page, limit int) ([]*entity.Video, int, error)

	// GetAnimeCollectionByUserAndVideoIDs 批量获取用户动漫收藏状态
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - videoIDs: 视频ID列表
	// 返回:
	//   - *map[int]bool: 视频ID到收藏状态的映射
	//   - error: 可能的错误信息
	GetAnimeCollectionByUserAndVideoIDs(ctx context.Context, userID int, videoIDs []int) (*map[int]bool, error)

	// GetAnimeCollectionByUserAndVideoID 获取用户动漫收藏状态
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户ID
	//   - videoID: 视频ID
	// 返回:
	//   - bool: 是否收藏
	//   - error: 可能的错误信息
	GetAnimeCollectionByUserAndVideoID(ctx context.Context, userID int, videoID int) (bool, error)

	// GetAnimeGenres 获取动漫类型
	// 参数:
	//   - ctx: 上下文信息
	//   - videoID: 视频ID
	// 返回:
	//   - []string: 动漫类型
	//   - error: 可能的错误信息
	GetAnimeGenres(ctx context.Context, videoID int) ([]string, error)
}
