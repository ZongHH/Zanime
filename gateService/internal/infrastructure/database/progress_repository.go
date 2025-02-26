package database

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
)

type ProgressRepositoryImpl struct {
	db *sql.DB
}

func NewProgressRepositoryImpl(db *sql.DB) *ProgressRepositoryImpl {
	return &ProgressRepositoryImpl{db: db}
}

func (p *ProgressRepositoryImpl) CreateProgress(ctx context.Context, progress *entity.Progress) error {
	return nil
}

func (p *ProgressRepositoryImpl) UpdateProgress(ctx context.Context, progress *entity.Progress) error {
	query := `
		INSERT INTO user_watch_progress 
			(video_id, user_id, episode, progress) VALUES (?, ?, ?, ?) 
		ON DUPLICATE KEY 
		UPDATE episode = VALUES(episode), progress = VALUES(progress)`
	_, err := p.db.ExecContext(ctx, query, progress.VideoID, progress.UserID, progress.Episode, progress.Progress)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProgressRepositoryImpl) GetProgress(ctx context.Context, userID int, page int, pageSize int) ([]*entity.Progress, error) {
	offset := (page - 1) * pageSize
	query := `
			SELECT w.id, w.user_id, w.video_id, w.episode, w.progress, w.status, w.updated_at, w.created_at, v.cover_image_url, v.video_name
			FROM user_watch_progress w
			JOIN anime_videos v ON w.video_id = v.video_id
			WHERE w.user_id = ? 
			ORDER BY w.updated_at DESC
			LIMIT ? OFFSET ?`
	rows, err := p.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []*entity.Progress
	for rows.Next() {
		var progress entity.Progress
		err := rows.Scan(&progress.ID, &progress.UserID, &progress.VideoID, &progress.Episode, &progress.Progress, &progress.Status, &progress.UpdatedAt, &progress.CreatedAt, &progress.CoverImageURL, &progress.VideoName)
		if err != nil {
			return nil, err
		}
		progresses = append(progresses, &progress)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return progresses, nil
}

func (p *ProgressRepositoryImpl) GetUserWatchProgress(ctx context.Context, userID, videoID int) (*entity.Progress, error) {
	query := `
        SELECT 
			v.video_name,
			v.area,
			v.release_date,
			COALESCE(uwp.episode, vu.default_episode) AS episode,  -- 优先用户进度，否则默认首集
			COALESCE(uwp.progress, 0) AS progress,                 -- 优先用户进度，否则0
			v.genres,
			v.cover_image_url
		FROM 
			-- 获取视频基础信息（含合并的 genres）
			(
				SELECT 
					video_id, 
					video_name, 
					area, 
					release_date, 
					GROUP_CONCAT(g.genre) AS genres, 
					cover_image_url
				FROM anime_videos v
				JOIN anime_genres g ON v.video_id = g.anime_id
				WHERE v.video_id = ?  -- 视频ID参数
				GROUP BY video_id
			) AS v
			-- 获取用户观看进度（可能为空）
			LEFT JOIN user_watch_progress uwp 
				ON v.video_id = uwp.video_id 
				AND uwp.user_id = ?  -- 用户ID参数
			-- 获取视频默认首集（确保必有值）
			CROSS JOIN (
				SELECT episode AS default_episode
				FROM video_urls 
				WHERE video_id = ?  -- 视频ID参数
				ORDER BY id ASC 
				LIMIT 1
			) AS vu;`

	progress := &entity.Progress{}
	err := p.db.QueryRowContext(ctx, query, videoID, userID, videoID).Scan(&progress.VideoName, &progress.Area, &progress.Release, &progress.Episode, &progress.Progress, &progress.Genre, &progress.CoverImageURL)
	if err != nil {
		return nil, err
	}

	return progress, nil
}
