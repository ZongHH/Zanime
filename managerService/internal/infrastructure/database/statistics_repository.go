package database

import (
	"context"
	"database/sql"
	"managerService/internal/domain/entity"
)

// StatisticsRepositoryImpl 实现了 StatisticsRepository 接口
type StatisticsRepositoryImpl struct {
	db *sql.DB
}

// NewStatisticsRepository 创建一个新的统计数据仓储实现
func NewStatisticsRepository(db *sql.DB) *StatisticsRepositoryImpl {
	return &StatisticsRepositoryImpl{
		db: db,
	}
}

// GetUserCount 获取系统中的总用户数量
func (r *StatisticsRepositoryImpl) GetUserCount(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM user_infos"
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetAnimeCount 获取系统中的动画总数量
func (r *StatisticsRepositoryImpl) GetAnimeCount(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM anime_videos"
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTodayPlayCount 获取今日的播放总次数
func (r *StatisticsRepositoryImpl) GetTodayPlayCount(ctx context.Context) (int, error) {
	return 0, nil
}

// GetActiveUserCount 获取活跃用户数量
func (r *StatisticsRepositoryImpl) GetActiveUserCount(ctx context.Context) (int, error) {
	return 0, nil
}

// GetNewAnime 获取最新上线的动漫
func (r *StatisticsRepositoryImpl) GetNewAnime(ctx context.Context, page int, pageSize int) ([]*entity.Anime, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT video_id, video_name, release_date, area, description, cover_image_url, 
		       views, likes, created_at, uploader_id
		FROM anime_videos
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	animes := make([]*entity.Anime, 0)
	for rows.Next() {
		anime := &entity.Anime{}
		err := rows.Scan(
			&anime.VideoID,
			&anime.VideoName,
			&anime.ReleaseDate,
			&anime.Area,
			&anime.Description,
			&anime.CoverImageURL,
			&anime.Views,
			&anime.Likes,
			&anime.CreatedAt,
			&anime.UploaderID,
		)
		if err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return animes, nil
}
