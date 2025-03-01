package database

import (
	"context"
	"database/sql"
	"fmt"
	"gateService/internal/domain/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

// PostCommentRepositoryImpl 实现了帖子评论相关的数据库操作
type PostCommentRepositoryImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewPostCommentRepositoryImpl 创建一个新的PostCommentRepositoryImpl实例
// 参数:
// - db: 数据库连接对象
// 返回:
// - *PostCommentRepositoryImpl: 帖子评论仓储实现的实例
func NewPostCommentRepositoryImpl(db *sql.DB, rdb *redis.Client) *PostCommentRepositoryImpl {
	return &PostCommentRepositoryImpl{db: db, rdb: rdb}
}

// BeginTx 开启事务
// 参数:
// - ctx: 上下文
// 返回:
// - *sql.Tx: 事务对象
// - error: 开启事务过程中的错误信息
func (p *PostCommentRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, nil)
}

// CreatePostComment 创建新的帖子评论
// 参数:
// - ctx: 上下文
// - postComment: 需要创建的帖子评论实体
// 返回:
// - error: 创建过程中的错误信息
func (p *PostCommentRepositoryImpl) CreatePostComment(ctx context.Context, postComment *entity.PostComment) error {
	query := `INSERT INTO post_comments (post_id, user_id, to_user_id, parent_id, root_id, content, level, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	var rootID sql.NullInt64
	if postComment.RootID != nil {
		rootID.Int64 = *postComment.RootID
		rootID.Valid = true
	}

	var parentID sql.NullInt64
	if postComment.ParentID != nil {
		parentID.Int64 = *postComment.ParentID
		parentID.Valid = true
	}

	var toUserID sql.NullInt32
	if postComment.ToUserID != nil {
		toUserID.Int32 = int32(*postComment.ToUserID)
		toUserID.Valid = true
	}

	result, err := p.db.ExecContext(ctx, query,
		postComment.PostID,
		postComment.UserID,
		toUserID,
		parentID,
		rootID,
		postComment.Content,
		postComment.Level,
		postComment.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	postComment.CommentID = id
	return nil
}

// CreatePostCommentTx 在事务中创建新的帖子评论
// 参数:
// - ctx: 上下文
// - tx: 事务对象
// - postComment: 需要创建的帖子评论实体
// 返回:
// - error: 创建过程中的错误信息
func (p *PostCommentRepositoryImpl) CreatePostCommentTx(ctx context.Context, tx *sql.Tx, postComment *entity.PostComment) error {
	query := `INSERT INTO post_comments (post_id, user_id, to_user_id, parent_id, root_id, content, level, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	var rootID sql.NullInt64
	if postComment.RootID != nil {
		rootID.Int64 = *postComment.RootID
		rootID.Valid = true
	}

	var parentID sql.NullInt64
	if postComment.ParentID != nil {
		parentID.Int64 = *postComment.ParentID
		parentID.Valid = true
	}

	var toUserID sql.NullInt32
	if postComment.ToUserID != nil {
		toUserID.Int32 = int32(*postComment.ToUserID)
		toUserID.Valid = true
	}

	result, err := tx.ExecContext(ctx, query,
		postComment.PostID,
		postComment.UserID,
		toUserID,
		parentID,
		rootID,
		postComment.Content,
		postComment.Level,
		postComment.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	postComment.CommentID = id
	return nil
}

// DeletePostComment 软删除帖子评论(将状态设置为0)
// 参数:
// - ctx: 上下文
// - postCommentID: 要删除的评论ID
// 返回:
// - error: 删除过程中的错误信息
func (p *PostCommentRepositoryImpl) DeletePostComment(ctx context.Context, postCommentID int64) error {
	query := "UPDATE post_comments SET status = 0 WHERE comment_id = ?"
	_, err := p.db.ExecContext(ctx, query, postCommentID)
	return err
}

// GetPostCommentsByPostID 根据帖子ID分页获取评论列表
// 参数:
// - ctx: 上下文
// - postID: 帖子ID
// - page: 页码,从1开始
// - pageSize: 每页评论数量
// 返回:
// - []*entity.PostComment: 评论列表
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetPostCommentsByPostID(ctx context.Context, postID int64, page int, pageSize int) ([]*entity.PostComment, error) {
	offset := (page - 1) * pageSize
	query := `SELECT comment_id, user_id, content, like_count, reply_count, created_at
		FROM post_comments 
		WHERE post_id = ? AND status = 1 AND level = 1
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := p.db.QueryContext(ctx, query, postID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*entity.PostComment
	for rows.Next() {
		comment := &entity.PostComment{}
		err := rows.Scan(
			&comment.CommentID,
			&comment.UserID,
			&comment.Content,
			&comment.LikeCount,
			&comment.ReplyCount,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetPostCommentsByRootID 根据根评论ID分页获取子评论列表
// 参数:
// - ctx: 上下文
// - rootID: 根评论ID
// - page: 页码,从1开始
// - pageSize: 每页评论数量
// 返回:
// - []*entity.PostComment: 子评论列表
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetPostCommentsByRootID(ctx context.Context, rootID int64, page int, pageSize int) ([]*entity.PostComment, error) {
	offset := (page - 1) * pageSize
	query := `SELECT comment_id, user_id, to_user_id, content, like_count, reply_count, created_at
		FROM post_comments 
		WHERE root_id = ? AND status = 1
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?`

	rows, err := p.db.QueryContext(ctx, query, rootID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.PostComment
	for rows.Next() {
		comment := &entity.PostComment{}
		err := rows.Scan(
			&comment.CommentID,
			&comment.UserID,
			&comment.ToUserID,
			&comment.Content,
			&comment.LikeCount,
			&comment.ReplyCount,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetPostCommentLikesByUserID 获取用户对评论的点赞状态
// 参数:
// - ctx: 上下文
// - userID: 用户ID
// 返回:
// - *map[int64]int8: 评论ID到点赞状态的映射(1:点赞 -1:踩 0:无状态)
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetPostCommentLikesByUserID(ctx context.Context, userID int) (*map[int64]int8, error) {
	query := `SELECT comment_id, status 
		FROM post_comment_likes 
		WHERE user_id = ?`

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likeMap := make(map[int64]int8)
	for rows.Next() {
		var commentID int64
		var status int8
		err := rows.Scan(&commentID, &status)
		if err != nil {
			return nil, err
		}
		likeMap[commentID] = status
	}
	return &likeMap, nil
}

// GetCommentTotalPage 获取评论总页数
// 参数:
// - ctx: 上下文
// - postID: 帖子ID
// - pageSize: 每页评论数量
// 返回:
// - int: 总页数
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetCommentTotalPage(ctx context.Context, postID int64, pageSize int) (int, error) {
	query := `SELECT COUNT(*) FROM post_comments WHERE post_id = ? AND status = 1 AND level = 1`

	var total int
	err := p.db.QueryRowContext(ctx, query, postID).Scan(&total)
	if err != nil {
		return 0, err
	}

	totalPage := (total + pageSize - 1) / pageSize

	return totalPage, nil
}

// UpdateReplyCount 更新评论的回复数量
// 参数:
// - ctx: 上下文
// - commentID: 评论ID
// - increment: 增加的回复数量,可以为负数表示减少
// 返回:
// - error: 更新过程中的错误信息
func (p *PostCommentRepositoryImpl) UpdateReplyCount(ctx context.Context, commentID int64, increment int) error {
	query := "UPDATE post_comments SET reply_count = reply_count + ? WHERE comment_id = ?"
	_, err := p.db.ExecContext(ctx, query, increment, commentID)
	return err
}

// UpdateReplyCountTx 在事务中更新评论的回复数量
// 参数:
// - ctx: 上下文
// - tx: 事务对象
// - commentID: 评论ID
// - increment: 增加的回复数量,可以为负数表示减少
// 返回:
// - error: 更新过程中的错误信息
func (p *PostCommentRepositoryImpl) UpdateReplyCountTx(ctx context.Context, tx *sql.Tx, commentID int64, increment int) error {
	query := "UPDATE post_comments SET reply_count = reply_count + ? WHERE comment_id = ?"
	_, err := tx.ExecContext(ctx, query, increment, commentID)
	return err
}

// UpdateLikeCount 更新评论的点赞数量
// 参数:
// - ctx: 上下文
// - commentID: 评论ID
// - increment: 增加的点赞数量,可以为负数表示减少
// 返回:
// - error: 更新过程中的错误信息
func (p *PostCommentRepositoryImpl) UpdateLikeCount(ctx context.Context, commentID int64, increment int) error {
	query := "UPDATE post_comments SET like_count = like_count + ? WHERE comment_id = ?"
	_, err := p.db.ExecContext(ctx, query, increment, commentID)
	return err
}

// InsertCommentLike 插入评论点赞记录
// 参数:
// - ctx: 上下文
// - commentLike: 评论点赞实体,包含用户ID、评论ID和点赞状态等信息
// 返回:
// - error: 插入过程中的错误信息
func (p *PostCommentRepositoryImpl) InsertCommentLike(ctx context.Context, commentLike *entity.PostCommentLike) error {
	query := `INSERT INTO post_comment_likes (user_id, comment_id, status) 
	VALUES (?, ?, ?) 
	ON DUPLICATE KEY UPDATE status = ?`

	_, err := p.db.ExecContext(ctx, query,
		commentLike.UserID,
		commentLike.CommentID,
		commentLike.Status,
		commentLike.Status)
	return err
}

// UpdateLikeCountTx 在事务中更新评论的点赞数量
// 参数:
// - ctx: 上下文
// - tx: 事务对象
// - commentID: 评论ID
// - increment: 增加的点赞数量,可以为负数表示减少
// 返回:
// - error: 更新过程中的错误信息
func (p *PostCommentRepositoryImpl) UpdateLikeCountTx(ctx context.Context, tx *sql.Tx, commentID int64, increment int) error {
	query := "UPDATE post_comments SET like_count = like_count + ? WHERE comment_id = ?"
	_, err := tx.ExecContext(ctx, query, increment, commentID)
	return err
}

// InsertCommentLikeTx 在事务中插入评论点赞记录
// 参数:
// - ctx: 上下文
// - tx: 事务对象
// - commentLike: 评论点赞实体,包含用户ID、评论ID和点赞状态等信息
// 返回:
// - error: 插入过程中的错误信息
func (p *PostCommentRepositoryImpl) InsertCommentLikeTx(ctx context.Context, tx *sql.Tx, commentLike *entity.PostCommentLike) error {
	query := `INSERT INTO post_comment_likes (user_id, comment_id, status) 
	VALUES (?, ?, ?) 
	ON DUPLICATE KEY UPDATE status = ?`

	_, err := tx.ExecContext(ctx, query,
		commentLike.UserID,
		commentLike.CommentID,
		commentLike.Status,
		commentLike.Status)
	return err
}

// SetCommentVirtualID 在Redis中设置评论的虚拟ID
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - userID: 用户ID
// - commentID: 评论的实际ID
// - virtualID: 要设置的虚拟ID
// 返回:
// - error: 设置过程中的错误信息,如Redis连接失败等
func (p *PostCommentRepositoryImpl) SetCommentVirtualID(ctx context.Context, userID int, commentID int64, virtualID int64) error {
	key := fmt.Sprintf("comment:virtual:user%d:%d", userID, commentID)
	err := p.rdb.Set(ctx, key, virtualID, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetCommentVirtualID 从Redis中获取评论的虚拟ID
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - userID: 用户ID
// - commentID: 评论的实际ID
// 返回:
// - int64: 评论对应的虚拟ID,如果不存在则返回0
// - error: 获取过程中的错误信息,如Redis连接失败等
func (p *PostCommentRepositoryImpl) GetCommentVirtualID(ctx context.Context, userID int, commentID int64) (int64, error) {
	key := fmt.Sprintf("comment:virtual:user%d:%d", userID, commentID)
	val, err := p.rdb.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}

// GetUserCommentCount 获取用户发布的评论数量
// 参数:
// - ctx: 上下文
// - userID: 用户ID
// 返回:
// - int: 用户发布的评论数量
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetUserCommentCount(ctx context.Context, userID int) (int, error) {
	query := "SELECT COUNT(*) FROM post_comments WHERE user_id = ? AND status = 1"
	var count int
	err := p.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentUserIDByCommentID 获取评论目标用户ID
// 参数:
// - ctx: 上下文
// - tx: 数据库事务
// - commentID: 评论ID
// 返回:
// - int: 目标用户ID
// - error: 获取过程中的错误信息
func (p *PostCommentRepositoryImpl) GetCommentUserIDByCommentID(ctx context.Context, tx *sql.Tx, commentID int64) (int, error) {
	query := "SELECT user_id FROM post_comments WHERE comment_id = ? AND status = 1"
	var userID int
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, commentID).Scan(&userID)
	} else {
		err = p.db.QueryRowContext(ctx, query, commentID).Scan(&userID)
	}

	if err != nil {
		return 0, fmt.Errorf("获取评论用户ID失败: %v", err)
	}

	return userID, nil
}
