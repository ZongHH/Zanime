package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

// CommentRepository 定义评论相关的数据库和缓存操作接口
// 提供评论的增删改查等基本操作方法
type CommentRepository interface {
	// 数据库操作

	// CreateFirstLevelComment 创建一级评论
	// 参数:
	// - ctx: 上下文
	// - comment: 需要创建的评论实体
	// 返回:
	// - error: 创建过程中的错误信息
	CreateFirstLevelComment(ctx context.Context, comment *entity.Comment) error

	// CreateSecondLevelComment 创建二级评论(回复)
	// 参数:
	// - ctx: 上下文
	// - comment: 需要创建的评论实体
	// 返回:
	// - error: 创建过程中的错误信息
	CreateSecondLevelComment(ctx context.Context, comment *entity.Comment) error

	// GetCommentsByVideoID 根据视频ID获取评论列表
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// - page: 分页页码,从1开始
	// - pageSize: 每页显示的评论数量
	// 返回:
	// - []*entity.Comment: 评论列表
	// - error: 获取过程中的错误信息
	GetCommentsByVideoID(ctx context.Context, videoID int, page int, pageSize int) ([]*entity.Comment, error)

	// GetCommentsByRootID 获取指定根评论下的所有子评论
	// 参数:
	// - ctx: 上下文
	// - rootID: 根评论ID
	// - page: 分页页码,从1开始
	// - pageSize: 每页显示的评论数量
	// 返回:
	// - []*entity.Comment: 子评论列表
	// - error: 获取过程中的错误信息
	GetCommentsByRootID(ctx context.Context, rootID int, page int, pageSize int) ([]*entity.Comment, error)

	// GetTotalCount 获取视频评论的总数量
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// 返回:
	// - int: 总数量
	// - error: 获取过程中的错误信息
	GetTotalCount(ctx context.Context, videoID int) (int, error)

	// 缓存操作

	// SetCommentsCacheByVideoID 缓存视频的评论列表
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// - comments: 要缓存的评论列表
	// 返回:
	// - error: 缓存过程中的错误信息
	SetCommentsCacheByVideoID(ctx context.Context, videoID int, comments []*entity.Comment) error

	// SetCommentsCacheByParentID 缓存父评论下的子评论列表
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// - parentID: 父评论ID
	// - comments: 要缓存的子评论列表
	// 返回:
	// - error: 缓存过程中的错误信息
	SetCommentsCacheByParentID(ctx context.Context, videoID int, parentID int, comments []*entity.Comment) error

	// GetCommentsCacheByVideoID 从缓存获取视频的评论列表
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// - page: 分页页码,从1开始
	// 返回:
	// - []*entity.Comment: 缓存的评论列表
	// - error: 获取过程中的错误信息
	GetCommentsCacheByVideoID(ctx context.Context, videoID int, page int) ([]*entity.Comment, error)

	// GetCommentsCacheByParentID 从缓存获取父评论下的子评论列表
	// 参数:
	// - ctx: 上下文
	// - videoID: 视频ID
	// - parentID: 父评论ID
	// - page: 分页页码,从1开始
	// 返回:
	// - []*entity.Comment: 缓存的子评论列表
	// - error: 获取过程中的错误信息
	GetCommentsCacheByParentID(ctx context.Context, videoID int, parentID int, page int) ([]*entity.Comment, error)
}
