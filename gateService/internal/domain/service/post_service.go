// package service 提供了与帖子相关的业务逻辑服务
// 包含帖子的创建、更新、删除、查询等核心功能
// 以及帖子评论的管理功能
package service

import (
	"context"
	"gateService/internal/domain/entity"
	"gateService/internal/interfaces/dto"
)

// PostService 定义了帖子服务的接口
// 提供了帖子管理的完整功能集,包括:
// - 帖子的基本CRUD操作
// - 帖子评论的管理
// - 用户帖子的查询
type PostService interface {
	// CreatePost 创建新的帖子
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - post: 创建帖子的请求数据,包含帖子的标题、内容等信息
	// 返回:
	// - *dto.CreatePostResponse: 创建帖子的响应数据
	// - error: 创建过程中的错误信息
	CreatePost(ctx context.Context, post *dto.CreatePostRequest) (*dto.CreatePostResponse, error)

	// UpdatePost 更新已有帖子
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - post: 需要更新的帖子实体,包含更新后的帖子信息
	// 返回:
	// - error: 更新过程中的错误信息,如果更新成功则返回nil
	UpdatePost(ctx context.Context, post *entity.Post) error

	// DeletePost 删除指定帖子
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - postID: 要删除的帖子ID
	// 返回:
	// - error: 删除过程中的错误信息,如果删除成功则返回nil
	DeletePost(ctx context.Context, postID int64) error

	// GetPostByID 根据帖子ID获取帖子信息
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - postID: 帖子ID
	// 返回:
	// - *entity.Post: 帖子实体信息
	// - error: 获取过程中的错误信息
	GetPostByID(ctx context.Context, postID int64) (*entity.Post, error)

	// GetPostCommentsByPostID 获取帖子的一级评论信息
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取评论的请求参数,包含分页信息等
	// 返回:
	// - *dto.GetCommentsResponse: 评论列表响应数据
	// - error: 获取过程中的错误信息
	GetPostCommentsByPostID(ctx context.Context, request *dto.GetPostCommentsRequest) (*dto.GetPostCommentsResponse, error)

	// GetPostCommentsByRootID 获取指定根评论的子评论列表
	// 参数:
	// - ctx: 上下文,用于传递请求上下文
	// - request: 获取评论请求,包含根评论ID、分页参数和当前用户ID等信息
	// 返回:
	// - []*dto.ReplyItem: 子评论列表数据
	// - error: 获取子评论过程中的错误信息
	GetPostCommentsByRootID(ctx context.Context, request *dto.GetPostCommentsRequest) ([]*dto.ReplyItem, error)

	// SubmitComment 提交评论或回复
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 提交评论的请求参数,包含评论内容、回复目标等
	// 返回:
	// - *dto.SubmitCommentResponse: 提交评论的响应数据
	// - error: 提交过程中的错误信息
	SubmitComment(ctx context.Context, request *dto.SubmitPostCommentRequest) (*dto.SubmitPostCommentResponse, error)

	// SubmitPostReply 提交帖子回复
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 提交回复的请求参数,包含回复内容、回复目标等
	// 返回:
	// - *dto.SubmitPostReplyResponse: 提交回复的响应数据
	// - error: 提交过程中的错误信息
	SubmitPostReply(ctx context.Context, request *dto.SubmitPostReplyRequest) (*dto.SubmitPostReplyResponse, error)

	// CommentLike 处理评论点赞/取消点赞的请求
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 点赞请求参数,包含用户ID、评论ID和点赞状态
	// 返回:
	// - *dto.CommentLikeResponse: 点赞操作的响应数据
	// - error: 点赞过程中的错误信息
	CommentLike(ctx context.Context, request *dto.CommentLikeRequest) (*dto.CommentLikeResponse, error)

	// GetPostCategoryList 获取帖子分类列表
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取分类列表的请求参数
	// 返回:
	// - *dto.GetCategoryListResponse: 帖子分类列表响应数据
	// - error: 获取过程中的错误信息
	GetPostCategoryList(ctx context.Context, request *dto.GetCategoryListRequest) (*dto.GetCategoryListResponse, error)

	// GetPostsByCategoryID 获取指定分类的帖子列表
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取帖子列表的请求参数,包含分类ID、页码和每页数量
	// 返回:
	// - *dto.GetPostListResponse: 帖子列表响应数据
	// - error: 获取过程中的错误信息
	GetPostsByCategoryID(ctx context.Context, request *dto.GetPostListRequest) (*dto.GetPostListResponse, error)

	// GetPostByPostID 根据PostID获取帖子
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取帖子请求参数,包含帖子ID
	// 返回:
	// - *dto.GetPostByPostIDResponse: 帖子响应数据
	// - error: 获取过程中的错误信息
	GetPostByPostID(ctx context.Context, request *dto.GetPostByPostIDRequest) (*dto.GetPostByPostIDResponse, error)

	// PostLike 处理帖子点赞/取消点赞的请求
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 点赞请求参数,包含用户ID、帖子ID和点赞状态
	// 返回:
	// - *dto.PostLikeResponse: 点赞操作的响应数据
	// - error: 点赞过程中的错误信息
	PostLike(ctx context.Context, request *dto.PostLikeRequest) (*dto.PostLikeResponse, error)

	// PostFavorite 处理帖子收藏/取消收藏的请求
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 收藏请求参数,包含用户ID、帖子ID和收藏状态
	// 返回:
	// - *dto.PostFavoriteResponse: 收藏操作的响应数据
	// - error: 收藏过程中的错误信息
	PostFavorite(ctx context.Context, request *dto.PostFavoriteRequest) (*dto.PostFavoriteResponse, error)

	// GetRecentPosts 获取最近发布的帖子列表
	// 参数:
	// - ctx: 上下文信息,用于传递请求上下文
	// - request: 获取最近帖子的请求参数,包含页码和每页数量
	// 返回:
	// - *dto.GetRecentPostsResponse: 最近帖子列表响应数据
	// - error: 获取过程中的错误信息
	GetRecentPosts(ctx context.Context, request *dto.GetRecentPostsRequest) (*dto.GetRecentPostsResponse, error)
}
