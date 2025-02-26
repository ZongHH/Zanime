package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gateService/internal/domain/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type CommentRepositoryImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewCommentRepositoryImpl(db *sql.DB, rdb *redis.Client) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		db:  db,
		rdb: rdb,
	}
}

func (c *CommentRepositoryImpl) CreateFirstLevelComment(ctx context.Context, comment *entity.Comment) error {
	// 一级评论
	query := "INSERT INTO comments (video_id, user_id, content) VALUES (?, ?, ?)"
	result, err := c.db.ExecContext(ctx, query, comment.VideoID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.CommentID = int(commentID)
	return nil
}

func (c *CommentRepositoryImpl) CreateSecondLevelComment(ctx context.Context, comment *entity.Comment) error {
	// 二级评论
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入二级评论
	query := "INSERT INTO comments (video_id, user_id, to_user_id, content, root_id, parent_id, level) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, comment.VideoID, comment.UserID, comment.ToUserID, comment.Content, comment.RootID, comment.ParentID, comment.Level)
	if err != nil {
		return err
	}

	// 更新父评论的回复数
	updateQuery := "UPDATE comments SET reply_num = reply_num + 1 WHERE comment_id = ?"
	_, err = tx.ExecContext(ctx, updateQuery, comment.RootID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.CommentID = int(commentID)
	return nil
}

func (c *CommentRepositoryImpl) GetCommentsByVideoID(ctx context.Context, videoID int, page int, pageSize int) ([]*entity.Comment, error) {
	query := `
			SELECT c.*, u.username, u.avatar_url 
			FROM comments c 
			JOIN user_infos u 
			ON c.user_id = u.user_id 
			WHERE c.video_id = ? AND c.level = 1 AND c.status = 1 LIMIT ? OFFSET ?`
	rows, err := c.db.QueryContext(ctx, query, videoID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(&comment.CommentID, &comment.VideoID, &comment.UserID, &comment.ToUserID,
			&comment.Content, &comment.RootID, &comment.ParentID, &comment.ReplyNum, &comment.Level,
			&comment.Status, &comment.CreatedAt, &comment.UserName, &comment.AvatarURL)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (c *CommentRepositoryImpl) GetCommentsByRootID(ctx context.Context, rootID int, page int, pageSize int) ([]*entity.Comment, error) {
	query := `
			SELECT c.*, u.username, u.avatar_url 
			FROM comments c 
			JOIN user_infos u 
			ON c.user_id = u.user_id 
			WHERE c.root_id = ? AND c.level = 2 AND c.status = 1 LIMIT ? OFFSET ?`
	rows, err := c.db.QueryContext(ctx, query, rootID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(&comment.CommentID, &comment.VideoID, &comment.UserID, &comment.ToUserID,
			&comment.Content, &comment.RootID, &comment.ParentID, &comment.ReplyNum, &comment.Level,
			&comment.Status, &comment.CreatedAt, &comment.UserName, &comment.AvatarURL)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (c *CommentRepositoryImpl) GetTotalCount(ctx context.Context, videoID int) (int, error) {
	// 查询总评论数
	var total int
	query := "SELECT COUNT(*) FROM comments WHERE video_id = ? AND level = 1 AND status = 1"
	err := c.db.QueryRowContext(ctx, query, videoID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (c *CommentRepositoryImpl) SetCommentsCacheByVideoID(ctx context.Context, videoID int, comments []*entity.Comment) error {
	// 设置缓存键
	key := fmt.Sprintf("video%d:comments", videoID)

	// 创建一个新的 Redis 事务
	pipe := c.rdb.Pipeline()

	// 将每个评论添加到有序集合中
	for _, comment := range comments {
		// 序列化评论数据
		jsonData, err := json.Marshal(comment)
		if err != nil {
			return err
		}

		// 使用回复数作为分数，将评论添加到有序集合
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(comment.ReplyNum),
			Member: string(jsonData),
		})
	}

	// 设置过期时间为1小时
	pipe.Expire(ctx, key, time.Hour)

	// 执行事务
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentRepositoryImpl) SetCommentsCacheByParentID(ctx context.Context, videoID int, parentID int, comments []*entity.Comment) error {
	// 设置缓存键
	key := fmt.Sprintf("video%d:comment%d:replies", videoID, parentID)

	// 创建一个新的 Redis 事务
	pipe := c.rdb.Pipeline()

	// 将每个评论添加到有序集合中
	for _, comment := range comments {
		// 序列化评论数据
		jsonData, err := json.Marshal(comment)
		if err != nil {
			return err
		}

		// 使用点赞数作为分数，将评论添加到有序集合
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(comment.ReplyNum),
			Member: string(jsonData),
		})
	}

	// 设置过期时间为1小时
	pipe.Expire(ctx, key, time.Hour)

	// 执行事务
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentRepositoryImpl) GetCommentsCacheByVideoID(ctx context.Context, videoID int, page int) ([]*entity.Comment, error) {
	// 获取缓存键
	key := fmt.Sprintf("video%d:comments", videoID)

	// 从有序集合中获取数据,按点赞数从高到低排序
	start := int64((page - 1) * 20)
	stop := int64(page*20 - 1) // 取20条数据

	members, err := c.rdb.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	// 如果members为空,说明缓存不存在
	if len(members) == 0 {
		return nil, nil
	}

	comments := make([]*entity.Comment, 0, len(members))
	for _, member := range members {
		var comment entity.Comment
		if err := json.Unmarshal([]byte(member), &comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (c *CommentRepositoryImpl) GetCommentsCacheByParentID(ctx context.Context, videoID int, parentID int, page int) ([]*entity.Comment, error) {
	// 获取缓存键
	key := fmt.Sprintf("video%d:comment%d:replies", videoID, parentID)

	// 从有序集合中获取数据,按点赞数从高到低排序
	start := int64((page - 1) * 8)
	stop := int64(page*8 - 1) // 取8条数据

	members, err := c.rdb.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	// 如果members为空,说明缓存不存在
	if len(members) == 0 {
		return nil, nil
	}

	comments := make([]*entity.Comment, 0, len(members))
	for _, member := range members {
		var comment entity.Comment
		if err := json.Unmarshal([]byte(member), &comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
