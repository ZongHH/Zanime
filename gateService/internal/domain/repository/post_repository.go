package repository

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
)

// PostRepository 定义帖子相关的数据库操作接口
type PostRepository interface {
	// BeginTx 开启事务
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - *sql.Tx: 事务对象
	// - error: 开启事务过程中的错误信息
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CreatePost 创建新帖子
	// 参数:
	// - ctx: 上下文
	// - post: 帖子实体
	// 返回:
	// - int64: 创建的帖子ID
	// - error: 错误信息
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)

	// UpdatePost 更新帖子信息
	// 参数:
	// - ctx: 上下文
	// - post: 需要更新的帖子实体
	UpdatePost(ctx context.Context, post *entity.Post) error

	// DeletePost 删除帖子
	// 参数:
	// - ctx: 上下文
	// - postID: 要删除的帖子ID
	DeletePost(ctx context.Context, postID int64) error

	// CreatePostTx 在事务中创建新帖子
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - post: 帖子实体
	// 返回:
	// - int64: 创建的帖子ID
	// - error: 错误信息
	CreatePostTx(ctx context.Context, tx *sql.Tx, post *entity.Post) (int64, error)

	// UpdatePostTx 在事务中更新帖子信息
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - post: 需要更新的帖子实体
	UpdatePostTx(ctx context.Context, tx *sql.Tx, post *entity.Post) error

	// DeletePostTx 在事务中删除帖子
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postID: 要删除的帖子ID
	DeletePostTx(ctx context.Context, tx *sql.Tx, postID int64) error

	// GetPostByID 根据ID获取帖子
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// 返回:
	// - *entity.Post: 帖子实体
	GetPostByID(ctx context.Context, postID int64) (*entity.Post, error)

	// GetPostsByUserID 获取用户的所有帖子
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// - page: 页码
	// - pageSize: 每页数量
	// 返回:
	// - []*entity.Post: 帖子列表
	// - error: 错误信息
	GetPostsByUserID(ctx context.Context, userID, page, pageSize int) ([]*entity.Post, error)

	// GetPostsByCategoryID 获取指定分类的帖子列表
	// 参数:
	// - ctx: 上下文
	// - categoryID: 分类ID
	// - page: 页码
	// - pageSize: 每页数量
	// 返回:
	// - []*entity.Post: 帖子列表
	GetPostsByCategoryID(ctx context.Context, categoryID int, page int, pageSize int) ([]*entity.Post, error)

	// GetPostByPostID 根据PostID获取帖子
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// 返回:
	// - *entity.Post: 帖子实体
	GetPostByPostID(ctx context.Context, postID int64) (*entity.Post, error)

	// GetPostCategoryList 获取帖子分类列表
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - []*entity.PostCategory: 帖子分类列表
	// - error: 错误信息
	GetPostCategoryList(ctx context.Context) ([]*entity.PostCategory, error)

	// UpdateViewCount 更新帖子浏览量
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// - increment: 增加的浏览量,可以为负数表示减少
	UpdateViewCount(ctx context.Context, postID int64, increment int) error

	// UpdateCommentCountTx 在事务中更新帖子评论量
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postID: 帖子ID
	// - increment: 增加的评论量,可以为负数表示减少
	UpdateCommentCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error

	// UpdateLikeCountTx 在事务中更新帖子点赞量
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postID: 帖子ID
	// - increment: 增加的点赞量,可以为负数表示减少
	UpdateLikeCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error

	// InsertPostLikeTx 在事务中插入帖子点赞记录
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postLike: 帖子点赞实体
	InsertPostLikeTx(ctx context.Context, tx *sql.Tx, postLike *entity.PostLike) error

	// UpdateCategoryCountTx 在事务中更新帖子分类数量
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - categoryID: 分类ID
	// - increment: 增加的数量,可以为负数表示减少
	UpdateCategoryCountTx(ctx context.Context, tx *sql.Tx, categoryID int, increment int) error

	// GetPostLikesByUserID 获取用户点赞的帖子
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - *map[int64]int8: 帖子ID与点赞状态的映射(1:点赞 -1:取消点赞 0:无状态)
	// - error: 错误信息
	GetPostLikesByUserID(ctx context.Context, userID int) (*map[int64]int8, error)

	// GetPostFavoritesByUserID 获取用户收藏的帖子
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - *map[int64]int8: 帖子ID与收藏状态的映射(1:收藏 0:未收藏)
	// - error: 错误信息
	GetPostFavoritesByUserID(ctx context.Context, userID int) (*map[int64]int8, error)

	// InsertPostFavoriteTx 在事务中插入帖子收藏记录
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postFavorite: 帖子收藏实体
	InsertPostFavoriteTx(ctx context.Context, tx *sql.Tx, postFavorite *entity.PostFavorite) error

	// UpdateFavoriteCountTx 在事务中更新帖子收藏数
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postID: 帖子ID
	// - increment: 增加的数量,可以为负数表示减少
	UpdateFavoriteCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error

	// GetUserPostCount 获取用户发布的帖子数量
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - int: 用户发布的帖子数量
	// - error: 获取过程中的错误信息
	GetUserPostCount(ctx context.Context, userID int) (int, error)

	// GetUserFavoritePostCount 获取用户收藏的帖子数量
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - int: 用户收藏的帖子数量
	// - error: 获取过程中的错误信息
	GetUserFavoritePostCount(ctx context.Context, userID int) (int, error)

	// CreatePostImagesTx 在事务中批量创建帖子图片
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postImages: 帖子图片列表
	CreatePostImagesTx(ctx context.Context, tx *sql.Tx, postImages []*entity.PostImage) error
}

// PostImageRepository 定义帖子图片相关的数据库操作接口
type PostImageRepository interface {
	// CreatePostImage 创建帖子图片
	// 参数:
	// - ctx: 上下文
	// - postImage: 帖子图片实体
	CreatePostImage(ctx context.Context, postImage *entity.PostImage) error

	// UpdatePostImage 更新帖子图片信息
	// 参数:
	// - ctx: 上下文
	// - postImage: 需要更新的帖子图片实体
	UpdatePostImage(ctx context.Context, postImage *entity.PostImage) error

	// DeletePostImage 删除帖子图片
	// 参数:
	// - ctx: 上下文
	// - postImageID: 要删除的图片ID
	DeletePostImage(ctx context.Context, postImageID int64) error

	// GetPostImageByID 根据ID获取帖子图片
	// 参数:
	// - ctx: 上下文
	// - postImageID: 图片ID
	// 返回:
	// - *entity.PostImage: 帖子图片实体
	GetPostImageByID(ctx context.Context, postImageID int64) (*entity.PostImage, error)
}

// PostTagRelationRepository 定义帖子标签关联的数据库操作接口
type PostTagRelationRepository interface {
	// BeginTx 开启事务
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - *sql.Tx: 事务对象
	// - error: 开启事务过程中的错误信息
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CreatePostTagRelation 创建帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - postTagRelation: 帖子标签关联实体数组
	CreatePostTagRelation(ctx context.Context, postTagRelation []*entity.PostTagRelation) error

	// UpdatePostTagRelation 更新帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - postTagRelation: 需要更新的帖子标签关联实体
	UpdatePostTagRelation(ctx context.Context, postTagRelation *entity.PostTagRelation) error

	// DeletePostTagRelation 删除帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - postTagRelation: 要删除的帖子标签关联实体
	DeletePostTagRelation(ctx context.Context, postTagRelation *entity.PostTagRelation) error

	// CreatePostTagRelationTx 在事务中创建帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTagRelation: 帖子标签关联实体数组
	CreatePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelation []*entity.PostTagRelation) error

	// UpdatePostTagRelationTx 在事务中更新帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTagRelation: 需要更新的帖子标签关联实体
	UpdatePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelation *entity.PostTagRelation) error

	// DeletePostTagRelationTx 在事务中删除帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTagRelation: 要删除的帖子标签关联实体
	DeletePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelation *entity.PostTagRelation) error

	// GetPostTagRelationByID 根据ID获取帖子标签关联
	// 参数:
	// - ctx: 上下文
	// - postTagRelationID: 关联ID
	// 返回:
	// - *entity.PostTagRelation: 帖子标签关联实体
	GetPostTagRelationByID(ctx context.Context, postTagRelationID int64) (*entity.PostTagRelation, error)
}

// PostTagRepository 定义帖子标签的数据库操作接口
type PostTagRepository interface {
	// BeginTx 开启事务
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - *sql.Tx: 事务对象
	// - error: 开启事务过程中的错误信息
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CreatePostTag 创建帖子标签
	// 参数:
	// - ctx: 上下文
	// - postTags: 帖子标签实体数组
	// 返回:
	// - []int: 创建的标签ID列表
	CreatePostTag(ctx context.Context, postTags []*entity.PostTag) ([]int, error)

	// UpdatePostTag 更新帖子标签
	// 参数:
	// - ctx: 上下文
	// - postTag: 需要更新的帖子标签实体
	UpdatePostTag(ctx context.Context, postTag *entity.PostTag) error

	// DeletePostTag 删除帖子标签
	// 参数:
	// - ctx: 上下文
	// - postTagID: 要删除的标签ID
	DeletePostTag(ctx context.Context, postTagID int64) error

	// CreatePostTagTx 在事务中创建帖子标签
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTags: 帖子标签实体数组
	// 返回:
	// - []int: 创建的标签ID列表
	CreatePostTagTx(ctx context.Context, tx *sql.Tx, postTags []*entity.PostTag) ([]int, error)

	// UpdatePostTagTx 在事务中更新帖子标签
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTag: 需要更新的帖子标签实体
	UpdatePostTagTx(ctx context.Context, tx *sql.Tx, postTag *entity.PostTag) error

	// DeletePostTagTx 在事务中删除帖子标签
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postTagID: 要删除的标签ID
	DeletePostTagTx(ctx context.Context, tx *sql.Tx, postTagID int64) error

	// GetPostTagByPostID 根据帖子ID获取帖子标签
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// 返回:
	// - []*entity.PostTag: 帖子标签列表
	// - error: 错误信息
	GetPostTagByPostID(ctx context.Context, postID int64) ([]*entity.PostTag, error)
}

// PostCommentRepository 定义帖子评论相关的数据库操作接口
// 提供帖子评论的增删改查等基本操作方法
type PostCommentRepository interface {
	// BeginTx 开启事务
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - *sql.Tx: 事务对象
	// - error: 开启事务过程中的错误信息
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CreatePostComment 创建帖子评论
	// 参数:
	// - ctx: 上下文
	// - postComment: 需要创建的帖子评论实体
	// 返回:
	// - error: 创建过程中的错误信息
	CreatePostComment(ctx context.Context, postComment *entity.PostComment) error

	// CreatePostCommentTx 在事务中创建新的帖子评论
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - postComment: 需要创建的帖子评论实体
	// 返回:
	// - error: 创建过程中的错误信息
	CreatePostCommentTx(ctx context.Context, tx *sql.Tx, postComment *entity.PostComment) error

	// DeletePostComment 删除帖子评论
	// 参数:
	// - ctx: 上下文
	// - postCommentID: 要删除的评论ID
	// 返回:
	// - error: 删除过程中的错误信息
	DeletePostComment(ctx context.Context, postCommentID int64) error

	// GetPostCommentsByPostID 根据帖子ID获取所有评论
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// - page: 分页页码,从1开始
	// - pageSize: 每页评论数量
	// 返回:
	// - []*entity.PostComment: 评论列表
	// - error: 获取过程中的错误信息
	GetPostCommentsByPostID(ctx context.Context, postID int64, page int, pageSize int) ([]*entity.PostComment, error)

	// GetPostCommentsByRootID 获取指定根评论下的所有子评论
	// 参数:
	// - ctx: 上下文
	// - rootID: 根评论ID
	// - page: 分页页码,从1开始
	// - pageSize: 每页评论数量
	// 返回:
	// - []*entity.PostComment: 子评论列表
	// - error: 获取过程中的错误信息
	GetPostCommentsByRootID(ctx context.Context, rootID int64, page int, pageSize int) ([]*entity.PostComment, error)

	// GetPostCommentLikesByUserID 获取用户对评论的点赞状态
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - *map[int64]int8: 评论ID与点赞状态的映射(1:点赞 -1:踩 0:无状态)
	// - error: 获取过程中的错误信息
	GetPostCommentLikesByUserID(ctx context.Context, userID int) (*map[int64]int8, error)

	// GetCommentTotalPage 获取评论总页数
	// 参数:
	// - ctx: 上下文
	// - postID: 帖子ID
	// - pageSize: 每页评论数量
	// 返回:
	// - int: 评论总页数
	// - error: 获取过程中的错误信息
	GetCommentTotalPage(ctx context.Context, postID int64, pageSize int) (int, error)

	// UpdateReplyCount 更新评论的回复数量
	// 参数:
	// - ctx: 上下文
	// - commentID: 评论ID
	// - increment: 增加的回复数量,可以为负数表示减少
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateReplyCount(ctx context.Context, commentID int64, increment int) error

	// UpdateReplyCountTx 在事务中更新评论的回复数量
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - commentID: 评论ID
	// - increment: 增加的回复数量,可以为负数表示减少
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateReplyCountTx(ctx context.Context, tx *sql.Tx, commentID int64, increment int) error

	// UpdateLikeCount 更新评论的点赞数量
	// 参数:
	// - ctx: 上下文
	// - commentID: 评论ID
	// - increment: 增加的点赞数量,可以为负数表示减少
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateLikeCount(ctx context.Context, commentID int64, increment int) error

	// InsertCommentLike 插入评论点赞记录
	// 参数:
	// - ctx: 上下文
	// - commentLike: 评论点赞实体,包含用户ID、评论ID和点赞状态等信息
	// 返回:
	// - error: 插入过程中的错误信息
	InsertCommentLike(ctx context.Context, commentLike *entity.PostCommentLike) error

	// UpdateLikeCountTx 在事务中更新评论的点赞数量
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - commentID: 评论ID
	// - increment: 增加的点赞数量,可以为负数表示减少
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateLikeCountTx(ctx context.Context, tx *sql.Tx, commentID int64, increment int) error

	// InsertCommentLikeTx 在事务中插入评论点赞记录
	// 参数:
	// - ctx: 上下文
	// - tx: 事务对象
	// - commentLike: 评论点赞实体,包含用户ID、评论ID和点赞状态等信息
	// 返回:
	// - error: 插入过程中的错误信息
	InsertCommentLikeTx(ctx context.Context, tx *sql.Tx, commentLike *entity.PostCommentLike) error

	// SetCommentVirtualID 在Redis中设置评论的虚拟ID
	// 参数:
	// - ctx: 上下文,用于传递请求上下文
	// - userID: 用户ID
	// - commentID: 评论的实际ID
	// - virtualID: 要设置的虚拟ID
	// 返回:
	// - error: 设置过程中的错误信息,如Redis连接失败等
	SetCommentVirtualID(ctx context.Context, userID int, commentID int64, virtualID int64) error

	// GetCommentVirtualID 从Redis中获取评论的虚拟ID
	// 参数:
	// - ctx: 上下文,用于传递请求上下文
	// - userID: 用户ID
	// - commentID: 评论的实际ID
	// 返回:
	// - int64: 评论对应的虚拟ID,如果不存在则返回0
	// - error: 获取过程中的错误信息,如Redis连接失败等
	GetCommentVirtualID(ctx context.Context, userID int, commentID int64) (int64, error)

	// GetUserCommentCount 获取用户发布的评论数量
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - int: 用户发布的评论数量
	// - error: 获取过程中的错误信息
	GetUserCommentCount(ctx context.Context, userID int) (int, error)

	// GetCommentUserIDByCommentID 获取评论目标用户ID
	// 参数:
	// - ctx: 上下文
	// - tx: 数据库事务
	// - commentID: 评论ID
	// 返回:
	// - int: 目标用户ID
	// - error: 获取过程中的错误信息
	GetCommentUserIDByCommentID(ctx context.Context, tx *sql.Tx, commentID int64) (int, error)
}
