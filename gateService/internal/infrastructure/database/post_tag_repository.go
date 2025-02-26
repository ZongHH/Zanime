package database

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
	"strings"
)

type PostTagRepositoryImpl struct {
	db *sql.DB
}

func NewPostTagRepositoryImpl(db *sql.DB) *PostTagRepositoryImpl {
	return &PostTagRepositoryImpl{db: db}
}

func (r *PostTagRepositoryImpl) CreatePostTag(ctx context.Context, postTags []*entity.PostTag) ([]int, error) {
	if len(postTags) == 0 {
		return nil, nil
	}

	// 构建批量插入的SQL语句
	query := `
		INSERT INTO post_tags (name, post_count) 
		VALUES `
	vals := []interface{}{}
	placeholders := make([]string, 0, len(postTags))

	for _, tag := range postTags {
		placeholders = append(placeholders, "(?, 1)")
		vals = append(vals, tag.Name)
	}

	query += strings.Join(placeholders, ",") + " ON DUPLICATE KEY UPDATE post_count = post_count + 1"

	result, err := r.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	// 获取插入的最后一个ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 构建返回的ID数组
	ids := make([]int, len(postTags))
	for i := range postTags {
		ids[i] = int(lastID) + i
	}

	return ids, nil
}

func (r *PostTagRepositoryImpl) UpdatePostTag(ctx context.Context, postTag *entity.PostTag) error {
	query := "UPDATE post_tags SET name = ? WHERE tag_id = ?"
	_, err := r.db.ExecContext(ctx, query, postTag.Name, postTag.TagID)
	return err
}

func (r *PostTagRepositoryImpl) DeletePostTag(ctx context.Context, postTagID int64) error {
	query := "DELETE FROM post_tags WHERE tag_id = ?"
	_, err := r.db.ExecContext(ctx, query, postTagID)
	return err
}

func (r *PostTagRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *PostTagRepositoryImpl) CreatePostTagTx(ctx context.Context, tx *sql.Tx, postTags []*entity.PostTag) ([]int, error) {
	if len(postTags) == 0 {
		return nil, nil
	}

	// 构建批量插入的SQL语句
	query := `
		INSERT INTO post_tags (name, post_count) 
		VALUES `
	vals := []interface{}{}
	placeholders := make([]string, 0, len(postTags))

	for _, tag := range postTags {
		placeholders = append(placeholders, "(?, 1)")
		vals = append(vals, tag.Name)
	}

	query += strings.Join(placeholders, ",") + " ON DUPLICATE KEY UPDATE post_count = post_count + 1"

	_, err := tx.ExecContext(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	// 构建查询条件
	nameParams := make([]string, 0, len(postTags))
	queryVals := make([]interface{}, 0, len(postTags))
	for _, tag := range postTags {
		nameParams = append(nameParams, "?")
		queryVals = append(queryVals, tag.Name)
	}

	// 查询插入或更新的标签ID
	selectQuery := "SELECT tag_id, name FROM post_tags WHERE name IN (" + strings.Join(nameParams, ",") + ")"
	rows, err := tx.QueryContext(ctx, selectQuery, queryVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建name到tag的映射
	tagMap := make(map[string]*entity.PostTag)
	for _, tag := range postTags {
		tagMap[tag.Name] = tag
	}

	// 填充TagID
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		if tag, ok := tagMap[name]; ok {
			tag.TagID = id
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 构建返回的ID数组,保持原有顺序
	ids := make([]int, len(postTags))
	for i, tag := range postTags {
		ids[i] = tag.TagID
	}

	return ids, nil
}

func (r *PostTagRepositoryImpl) UpdatePostTagTx(ctx context.Context, tx *sql.Tx, postTag *entity.PostTag) error {
	query := "UPDATE post_tags SET name = ? WHERE tag_id = ?"
	_, err := tx.ExecContext(ctx, query, postTag.Name, postTag.TagID)
	return err
}

func (r *PostTagRepositoryImpl) DeletePostTagTx(ctx context.Context, tx *sql.Tx, postTagID int64) error {
	query := "DELETE FROM post_tags WHERE tag_id = ?"
	_, err := tx.ExecContext(ctx, query, postTagID)
	return err
}

func (r *PostTagRepositoryImpl) GetPostTagByPostID(ctx context.Context, postID int64) ([]*entity.PostTag, error) {
	// 先从关联表获取tag_ids
	relationQuery := "SELECT tag_id FROM post_tag_relations WHERE post_id = ?"
	rows, err := r.db.QueryContext(ctx, relationQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 收集所有tag_id
	var tagIDs []int
	for rows.Next() {
		var tagID int
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, tagID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tagIDs) == 0 {
		return []*entity.PostTag{}, nil
	}

	// 构建IN查询的占位符
	placeholders := make([]string, len(tagIDs))
	args := make([]interface{}, len(tagIDs))
	for i, id := range tagIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// 根据tag_ids获取标签详情
	tagQuery := "SELECT tag_id, name, post_count, created_at FROM post_tags WHERE tag_id IN (" + strings.Join(placeholders, ",") + ")"
	tagRows, err := r.db.QueryContext(ctx, tagQuery, args...)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	var postTags []*entity.PostTag
	for tagRows.Next() {
		var postTag entity.PostTag
		if err := tagRows.Scan(&postTag.TagID, &postTag.Name, &postTag.PostCount, &postTag.CreatedAt); err != nil {
			return nil, err
		}
		postTags = append(postTags, &postTag)
	}
	if err = tagRows.Err(); err != nil {
		return nil, err
	}

	return postTags, nil
}
