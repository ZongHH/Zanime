package database

import (
	"context"
	"database/sql"
	"fmt"
	"gateService/internal/domain/entity"
	"strconv"
	"strings"
)

type PostRepositoryImpl struct {
	db *sql.DB
}

func NewPostRepositoryImpl(db *sql.DB) *PostRepositoryImpl {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *PostRepositoryImpl) CreatePost(ctx context.Context, post *entity.Post) (int64, error) {
	query := "INSERT INTO posts (user_id, category_id, title, content) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, post.UserID, post.CategoryID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (r *PostRepositoryImpl) UpdatePost(ctx context.Context, post *entity.Post) error {
	query := "UPDATE posts SET title = ?, content = ?, view_count = ?, like_count = ?, comment_count = ?, favorite_count = ?, is_pinned = ?, is_featured = ?, status = ? WHERE post_id = ?"
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ViewCount, post.LikeCount, post.CommentCount, post.FavoriteCount, post.IsPinned, post.IsFeatured, post.Status, post.PostID)
	return err
}

func (r *PostRepositoryImpl) DeletePost(ctx context.Context, postID int64) error {
	query := "DELETE FROM posts WHERE post_id = ?"
	_, err := r.db.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) CreatePostTx(ctx context.Context, tx *sql.Tx, post *entity.Post) (int64, error) {
	query := "INSERT INTO posts (user_id, category_id, title, content) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, post.UserID, post.CategoryID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (r *PostRepositoryImpl) UpdatePostTx(ctx context.Context, tx *sql.Tx, post *entity.Post) error {
	query := "UPDATE posts SET title = ?, content = ?, view_count = ?, like_count = ?, comment_count = ?, favorite_count = ?, is_pinned = ?, is_featured = ?, status = ? WHERE post_id = ?"
	_, err := tx.ExecContext(ctx, query, post.Title, post.Content, post.ViewCount, post.LikeCount, post.CommentCount, post.FavoriteCount, post.IsPinned, post.IsFeatured, post.Status, post.PostID)
	return err
}

func (r *PostRepositoryImpl) DeletePostTx(ctx context.Context, tx *sql.Tx, postID int64) error {
	query := "DELETE FROM posts WHERE post_id = ?"
	_, err := tx.ExecContext(ctx, query, postID)
	return err
}

func (r *PostRepositoryImpl) GetPostByID(ctx context.Context, postID int64) (*entity.Post, error) {
	query := "SELECT * FROM posts WHERE post_id = ?"
	row := r.db.QueryRowContext(ctx, query, postID)
	var post entity.Post
	err := row.Scan(&post.PostID, &post.UserID, &post.CategoryID, &post.Title, &post.Content, &post.ViewCount, &post.LikeCount, &post.CommentCount, &post.FavoriteCount, &post.IsPinned, &post.IsFeatured, &post.Status, &post.CreatedAt, &post.UpdatedAt)
	return &post, err
}

func (r *PostRepositoryImpl) GetPostsByUserID(ctx context.Context, userID, page, pageSize int) ([]*entity.Post, error) {
	query := "SELECT * FROM posts WHERE user_id = ? ORDER BY created_at DESC LIMIT ?, ?"
	rows, err := r.db.QueryContext(ctx, query, userID, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.CategoryID, &post.Title, &post.Content, &post.ViewCount, &post.LikeCount, &post.CommentCount, &post.FavoriteCount, &post.IsPinned, &post.IsFeatured, &post.Status, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *PostRepositoryImpl) GetPostsByCategoryID(ctx context.Context, categoryID int, page int, pageSize int) ([]*entity.Post, error) {
	query := `
			SELECT p.*, COALESCE(GROUP_CONCAT(pi.image_id), '') AS image_ids, COALESCE(GROUP_CONCAT(pi.image_url), '') AS image_urls
			FROM posts p
			LEFT JOIN post_images pi ON p.post_id = pi.post_id
			WHERE p.category_id = ? AND p.status = 1
			GROUP BY p.post_id
			ORDER BY p.created_at DESC
			LIMIT ?, ?
			`
	rows, err := r.db.QueryContext(ctx, query, categoryID, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		var imageIDs string
		var imageURLs string
		err := rows.Scan(&post.PostID, &post.UserID, &post.CategoryID, &post.Title, &post.Content,
			&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.FavoriteCount,
			&post.IsPinned, &post.IsFeatured, &post.Status, &post.CreatedAt, &post.UpdatedAt, &imageIDs, &imageURLs)
		if err != nil {
			return nil, err
		}

		if imageIDs != "" && imageURLs != "" {
			IDs := strings.Split(imageIDs, ",")
			URLs := strings.Split(imageURLs, ",")
			for i, imageURL := range URLs {
				imageID, err := strconv.ParseInt(IDs[i], 10, 64)
				if err != nil {
					return nil, err
				}
				post.Images = append(post.Images, entity.PostImage{
					ImageID:  imageID,
					ImageURL: imageURL,
				})
			}
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *PostRepositoryImpl) GetPostByPostID(ctx context.Context, postID int64) (*entity.Post, error) {
	query := `
			SELECT p.*, COALESCE(GROUP_CONCAT(pi.image_id), '') AS image_ids, COALESCE(GROUP_CONCAT(pi.image_url), '') AS image_urls
			FROM posts p
			LEFT JOIN post_images pi ON p.post_id = pi.post_id
			WHERE p.post_id = ? AND p.status = 1
			GROUP BY p.post_id
			`
	row := r.db.QueryRowContext(ctx, query, postID)
	var post entity.Post
	var imageIDs string
	var imageURLs string
	err := row.Scan(&post.PostID, &post.UserID, &post.CategoryID, &post.Title, &post.Content,
		&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.FavoriteCount,
		&post.IsPinned, &post.IsFeatured, &post.Status, &post.CreatedAt, &post.UpdatedAt, &imageIDs, &imageURLs)
	if err != nil {
		return nil, err
	}

	if imageIDs != "" && imageURLs != "" {
		IDs := strings.Split(imageIDs, ",")
		URLs := strings.Split(imageURLs, ",")
		for i, imageURL := range URLs {
			imageID, err := strconv.ParseInt(IDs[i], 10, 64)
			if err != nil {
				return nil, err
			}
			post.Images = append(post.Images, entity.PostImage{
				ImageID:  imageID,
				ImageURL: imageURL,
			})
		}
	}
	return &post, nil
}

func (r *PostRepositoryImpl) GetPostCategoryList(ctx context.Context) ([]*entity.PostCategory, error) {
	query := "SELECT category_id, name, icon, post_count FROM post_categories WHERE status = 1 ORDER BY sort_order DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.PostCategory
	for rows.Next() {
		var category entity.PostCategory
		err := rows.Scan(&category.CategoryID, &category.Name, &category.Icon, &category.PostCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

func (r *PostRepositoryImpl) UpdateViewCount(ctx context.Context, postID int64, increment int) error {
	query := "UPDATE posts SET view_count = view_count + ? WHERE post_id = ?"
	result, err := r.db.ExecContext(ctx, query, increment, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("更新帖子浏览量失败,帖子不存在")
	}

	return nil
}

func (r *PostRepositoryImpl) UpdateCommentCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error {
	query := "UPDATE posts SET comment_count = comment_count + ? WHERE post_id = ?"
	result, err := tx.ExecContext(ctx, query, increment, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("更新帖子评论数失败,帖子不存在")
	}

	return nil
}

func (r *PostRepositoryImpl) UpdateLikeCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error {
	query := "UPDATE posts SET like_count = like_count + ? WHERE post_id = ?"
	result, err := tx.ExecContext(ctx, query, increment, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("更新帖子点赞数失败,帖子不存在")
	}

	return nil
}

func (r *PostRepositoryImpl) InsertPostLikeTx(ctx context.Context, tx *sql.Tx, postLike *entity.PostLike) error {
	query := `INSERT INTO post_likes (user_id, post_id, status)
			  VALUES (?, ?, ?)
			  ON DUPLICATE KEY UPDATE status = VALUES(status)`
	_, err := tx.ExecContext(ctx, query, postLike.UserID, postLike.PostID, postLike.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepositoryImpl) UpdateCategoryCountTx(ctx context.Context, tx *sql.Tx, categoryID int, increment int) error {
	query := "UPDATE post_categories SET post_count = post_count + ? WHERE category_id = ?"
	result, err := tx.ExecContext(ctx, query, increment, categoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("更新分类帖子数量失败,分类不存在")
	}

	return nil
}

func (r *PostRepositoryImpl) GetPostLikesByUserID(ctx context.Context, userID int) (*map[int64]int8, error) {
	query := "SELECT post_id, status FROM post_likes WHERE user_id = ?"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likeMap := make(map[int64]int8)
	for rows.Next() {
		var postID int64
		var status int8
		if err := rows.Scan(&postID, &status); err != nil {
			return nil, err
		}
		likeMap[postID] = status
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &likeMap, nil
}

func (r *PostRepositoryImpl) GetPostFavoritesByUserID(ctx context.Context, userID int) (*map[int64]int8, error) {
	query := "SELECT post_id, status FROM post_favorites WHERE user_id = ?"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favoriteMap := make(map[int64]int8)
	for rows.Next() {
		var postID int64
		var status int8
		if err := rows.Scan(&postID, &status); err != nil {
			return nil, err
		}
		favoriteMap[postID] = status
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &favoriteMap, nil
}

func (r *PostRepositoryImpl) InsertPostFavoriteTx(ctx context.Context, tx *sql.Tx, postFavorite *entity.PostFavorite) error {
	query := `INSERT INTO post_favorites (user_id, post_id, status) 
			  VALUES (?, ?, ?)
			  ON DUPLICATE KEY UPDATE status = ?`

	_, err := tx.ExecContext(ctx, query,
		postFavorite.UserID,
		postFavorite.PostID,
		postFavorite.Status,
		postFavorite.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepositoryImpl) UpdateFavoriteCountTx(ctx context.Context, tx *sql.Tx, postID int64, increment int) error {
	query := "UPDATE posts SET favorite_count = favorite_count + ? WHERE post_id = ?"

	_, err := tx.ExecContext(ctx, query, increment, postID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepositoryImpl) GetUserPostCount(ctx context.Context, userID int) (int, error) {
	query := "SELECT COUNT(*) FROM posts WHERE user_id = ? AND status = 1"
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepositoryImpl) GetUserFavoritePostCount(ctx context.Context, userID int) (int, error) {
	query := "SELECT COUNT(*) FROM post_favorites WHERE user_id = ? AND status = 1"
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepositoryImpl) CreatePostImagesTx(ctx context.Context, tx *sql.Tx, postImages []*entity.PostImage) error {
	query := `INSERT INTO post_images (post_id, image_url) 
			  VALUES (?, ?)`

	for _, image := range postImages {
		_, err := tx.ExecContext(ctx, query,
			image.PostID,
			image.ImageURL)
		if err != nil {
			return err
		}
	}

	return nil
}
