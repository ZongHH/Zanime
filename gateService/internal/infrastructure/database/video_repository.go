package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gateService/internal/domain/entity"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type VideoRepositoryImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewVideoRepositoryImpl(db *sql.DB, rdb *redis.Client) *VideoRepositoryImpl {
	return &VideoRepositoryImpl{
		db:  db,
		rdb: rdb,
	}
}

func (r *VideoRepositoryImpl) CreateVideo(ctx context.Context, video *entity.Video) error {
	query := `INSERT INTO anime_videos (video_name, release_date, area, description, cover_image_url, uploader_id) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, video.Name, video.ReleaseDate, video.Area, video.Description, video.CoverImageUrl, video.UploaderID)
	return err
}

func (r *VideoRepositoryImpl) UpdateVideo(ctx context.Context, video *entity.Video) error {
	query := `UPDATE anime_videos SET video_name = ?, release_date = ?, area = ?, description = ?, cover_image_url = ?, uploader_id = ? WHERE video_id = ?`
	_, err := r.db.ExecContext(ctx, query, video.Name, video.ReleaseDate, video.Area, video.Description, video.CoverImageUrl, video.UploaderID, video.ID)
	return err
}

func (r *VideoRepositoryImpl) DeleteVideo(ctx context.Context, videoID int) error {
	query := `DELETE FROM anime_videos WHERE video_id = ?`
	_, err := r.db.ExecContext(ctx, query, videoID)
	return err
}

func (r *VideoRepositoryImpl) GetVideoByID(ctx context.Context, videoID int) (*entity.Video, error) {
	query := `SELECT * FROM anime_videos WHERE video_id = ?`
	row := r.db.QueryRowContext(ctx, query, videoID)
	var video entity.Video
	err := row.Scan(&video.ID, &video.Name, &video.ReleaseDate, &video.Area, &video.Description, &video.CoverImageUrl, &video.Views, &video.Likes, &video.CreatedAt, &video.UploaderID)
	return &video, err
}

func (r *VideoRepositoryImpl) GetVideosByVideoName(ctx context.Context, videoName string) ([]*entity.Video, error) {
	query := `SELECT video_id, video_name, cover_image_url 
			  FROM anime_videos
			  WHERE video_name LIKE ?
			  LIMIT 10 OFFSET 0`
	rows, err := r.db.QueryContext(ctx, query, "%"+videoName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*entity.Video
	for rows.Next() {
		var video entity.Video
		err := rows.Scan(&video.ID, &video.Name, &video.CoverImageUrl)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}
	return videos, nil
}

func (r *VideoRepositoryImpl) GetVideosALLEpisodesByVideoName(ctx context.Context, videoName string, page int) ([]*entity.Video, error) {
	// 建议设置group_concat_max_len为更大的值以确保所有episode都能连接
	// _, err := r.db.ExecContext(ctx, "SET SESSION group_concat_max_len = 4096")
	// if err != nil {
	// 	return nil, err
	// }

	query := `
		SELECT 
			v.video_id,
			v.video_name,
			v.cover_image_url,
			v.release_date,
			v.area,
			v.description,
			GROUP_CONCAT(DISTINCT g.genre) AS genres,
			GROUP_CONCAT(DISTINCT u.episode ORDER BY u.id ASC) AS episodes
		FROM anime_videos v
		LEFT JOIN anime_genres g ON v.video_id = g.anime_id 
		LEFT JOIN video_urls u ON v.video_id = u.video_id
		WHERE v.video_name LIKE ?
		GROUP BY 
			v.video_id,
			v.video_name,
			v.cover_image_url,
			v.release_date,
			v.area,
			v.description
		LIMIT 10 OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, "%"+videoName+"%", (page-1)*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*entity.Video
	for rows.Next() {
		var episodes string
		var video entity.Video
		err := rows.Scan(&video.ID, &video.Name, &video.CoverImageUrl, &video.ReleaseDate, &video.Area, &video.Description, &video.Genres, &episodes)
		if err != nil {
			return nil, err
		}
		video.Episodes = strings.Split(episodes, ",")
		videos = append(videos, &video)
	}
	return videos, nil
}

func (r *VideoRepositoryImpl) GetVideoInfoWithEposidesByVideoID(ctx context.Context, videoID int) (*entity.Video, error) {
	// 获取名字、选集信息
	query := `
		SELECT 
			a.video_id,
			a.video_name, 
			GROUP_CONCAT(v.episode ORDER BY v.id ASC SEPARATOR ',') AS episodes
		FROM anime_videos a
			JOIN video_urls v ON a.video_id = v.video_id
		WHERE a.video_id = ?
		GROUP BY a.video_name
	`
	video := &entity.Video{}
	var episodes string
	err := r.db.QueryRowContext(ctx, query, videoID).Scan(&video.ID, &video.Name, &episodes)
	if err != nil {
		return nil, err
	}
	video.Episodes = strings.Split(episodes, ",")
	return video, nil
}

func (r *VideoRepositoryImpl) GetVideosByFilters(ctx context.Context, video *entity.Video, page, limit int) ([]*entity.Video, int, error) {
	type queryParams struct {
		baseQuery  string
		conditions []string
		args       []interface{}
	}

	// 公共查询条件构建器
	buildQuery := func(base string) *queryParams {
		qp := &queryParams{
			baseQuery:  base,
			conditions: []string{"v.release_date != '未知'"},
			args:       make([]interface{}, 0),
		}

		if video.Area != "" {
			qp.conditions = append(qp.conditions, "v.area = ?")
			qp.args = append(qp.args, video.Area)
		}
		if video.ReleaseDate != "" {
			qp.conditions = append(qp.conditions, "v.release_date = ?")
			qp.args = append(qp.args, video.ReleaseDate)
		}
		if video.Genres != "" {
			qp.conditions = append(qp.conditions, "g.genre = ?")
			qp.args = append(qp.args, video.Genres)
		}
		if video.Initial != "" {
			qp.conditions = append(qp.conditions, "v.video_name LIKE ?")
			qp.args = append(qp.args, video.Initial+"%")
		}

		return qp
	}

	var (
		wg      sync.WaitGroup
		library []*entity.Video
		total   int
		errChan = make(chan error, 2) // 缓冲通道防止阻塞
	)

	// 主查询协程
	wg.Add(1)
	go func() {
		defer wg.Done()

		qp := buildQuery(`
			SELECT 
				v.video_id, 
				v.video_name, 
				v.release_date, 
				v.cover_image_url,
				GROUP_CONCAT(DISTINCT g.genre) AS genres
			FROM anime_videos v
			JOIN anime_genres g ON v.video_id = g.anime_id`)

		// 构建完整查询
		query := fmt.Sprintf("%s WHERE %s GROUP BY v.video_id ORDER BY v.release_date DESC LIMIT ? OFFSET ?",
			qp.baseQuery, strings.Join(qp.conditions, " AND "))
		args := append(qp.args, limit, (page-1)*limit)

		// 使用预编译语句防止 SQL 注入
		rows, err := r.db.QueryContext(ctx, query, args...)
		if err != nil {
			errChan <- fmt.Errorf("query failed: %w", err)
			return
		}
		defer rows.Close()

		// 使用更安全的扫描方式
		for rows.Next() {
			var video entity.Video
			if err := rows.Scan(
				&video.ID,
				&video.Name,
				&video.ReleaseDate,
				&video.CoverImageUrl,
				&video.Genres,
			); err != nil {
				errChan <- fmt.Errorf("scan failed: %w", err)
				return
			}
			video.Rating = "9.5" // 需要确认评分来源
			library = append(library, &video)
		}

		if err := rows.Err(); err != nil {
			errChan <- fmt.Errorf("rows error: %w", err)
		}
	}()

	// 总数查询协程
	wg.Add(1)
	go func() {
		defer wg.Done()

		qp := buildQuery(`
			SELECT COUNT(DISTINCT v.video_id)
			FROM anime_videos v
			JOIN anime_genres g ON v.video_id = g.anime_id`)

		query := fmt.Sprintf("%s WHERE %s", qp.baseQuery, strings.Join(qp.conditions, " AND "))

		if err := r.db.QueryRowContext(ctx, query, qp.args...).Scan(&total); err != nil {
			errChan <- fmt.Errorf("count query failed: %w", err)
		}
	}()

	// 优雅关闭处理
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 错误优先处理
	if err := <-errChan; err != nil {
		return nil, 0, err
	}

	return library, total, nil
}

func (r *VideoRepositoryImpl) GetVideoFilters(ctx context.Context) (*entity.VideoFilters, error) {
	// 1. 并发执行三个独立查询
	var (
		wg          sync.WaitGroup
		regions     []string
		releases    []string
		genres      []string
		queryErrors [3]error
	)

	// 2. 查询地区分布
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := `
			SELECT COALESCE(GROUP_CONCAT(area ORDER BY cnt DESC SEPARATOR ','), '') 
			FROM (
				SELECT area, COUNT(*) as cnt 
				FROM anime_videos 
				WHERE release_date != '未知' 
				GROUP BY area
			) t_areas`
		var s string
		err := r.db.QueryRowContext(ctx, query).Scan(&s) // 3. 添加上下文传递
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			queryErrors[0] = fmt.Errorf("regions query failed: %w", err)
			return
		}
		regions = strings.Split(s, ",")
	}()

	// 4. 查询发行年份
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := `
			SELECT COALESCE(GROUP_CONCAT(release_date ORDER BY release_date DESC SEPARATOR ','), '') 
			FROM (
				SELECT DISTINCT release_date 
				FROM anime_videos 
				WHERE release_date != '未知'
			) t_dates`
		var s string
		err := r.db.QueryRowContext(ctx, query).Scan(&s)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			queryErrors[1] = fmt.Errorf("releases query failed: %w", err)
			return
		}
		releases = strings.Split(s, ",")
	}()

	// 5. 查询类型分布
	wg.Add(1)
	go func() {
		defer wg.Done()
		query := `
			SELECT COALESCE(GROUP_CONCAT(genre ORDER BY cnt DESC SEPARATOR ','), '') 
			FROM (
				SELECT genre, COUNT(*) as cnt 
				FROM anime_genres 
				WHERE EXISTS (
					SELECT 1 
					FROM anime_videos 
					WHERE video_id = anime_id 
					AND release_date != '未知'
				)
				GROUP BY genre
			) t_genres`
		var s string
		err := r.db.QueryRowContext(ctx, query).Scan(&s)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			queryErrors[2] = fmt.Errorf("genres query failed: %w", err)
			return
		}
		genres = strings.Split(s, ",")
	}()

	wg.Wait()

	// 6. 统一错误处理
	for _, err := range queryErrors {
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}
	}

	// 7. 字母列表动态生成
	letters := make([]string, 0, 26)
	for c := 'A'; c <= 'Z'; c++ {
		letters = append(letters, string(c))
	}

	return &entity.VideoFilters{
		Regions: regions,
		Years:   releases,
		Types:   genres,
		Letters: letters,
	}, nil
}

func (r *VideoRepositoryImpl) GetAnimesByGenre(ctx context.Context, genre string, page, limit int) ([]*entity.Video, error) {
	query := `
		SELECT v.video_id, v.video_name, v.cover_image_url
		FROM anime_videos v
		JOIN anime_genres g ON v.video_id = g.anime_id
		WHERE g.genre = ? AND v.release_date != '未知'
		ORDER BY v.release_date DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, genre, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		var video entity.Video
		err := rows.Scan(&video.ID, &video.Name, &video.CoverImageUrl)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}
	return videos, nil
}

func (r *VideoRepositoryImpl) GetTopAnimeGenres(ctx context.Context) ([]string, error) {
	query := `
		SELECT genre
		FROM anime_genres
		GROUP BY genre
		ORDER BY COUNT(*) DESC
		LIMIT 8
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []string
	for rows.Next() {
		var genre string
		err := rows.Scan(&genre)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *VideoRepositoryImpl) CacheVideoURL(ctx context.Context, key string, value string) error {
	err := r.rdb.Set(ctx, key, value, 4*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *VideoRepositoryImpl) GetVideoURL(ctx context.Context, key string) (string, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *VideoRepositoryImpl) AddAnimeCollection(ctx context.Context, collection *entity.UserAnimeCollection) error {
	query := `INSERT INTO user_anime_collections (user_id, video_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE status = 1`
	_, err := r.db.ExecContext(ctx, query, collection.UserID, collection.VideoID)
	return err
}

func (r *VideoRepositoryImpl) DeleteAnimeCollection(ctx context.Context, userID, videoID int) error {
	query := `UPDATE user_anime_collections SET status = 0 WHERE user_id = ? AND video_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, videoID)
	return err
}

func (r *VideoRepositoryImpl) GetAnimeCollectionByUser(ctx context.Context, userID int, page, limit int) ([]*entity.Video, int, error) {
	// 计算总数
	countQuery := `
		SELECT COUNT(*) 
		FROM user_anime_collections
		WHERE user_id = ? AND status = 1
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取收藏列表
	offset := (page - 1) * limit
	query := `
		SELECT 
			v.video_id,
			v.video_name,
			v.cover_image_url,
			uc.updated_at
		FROM user_anime_collections uc
		JOIN anime_videos v ON uc.video_id = v.video_id
		WHERE uc.user_id = ? AND uc.status = 1
		ORDER BY uc.updated_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		var video entity.Video
		err := rows.Scan(
			&video.ID,
			&video.Name,
			&video.CoverImageUrl,
			&video.CollectedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		videos = append(videos, &video)
	}

	return videos, total, nil
}

func (r *VideoRepositoryImpl) GetAnimeCollectionByUserAndVideoIDs(ctx context.Context, userID int, videoIDs []int) (*map[int]bool, error) {
	if len(videoIDs) == 0 {
		return nil, nil
	}

	// 构建IN查询的参数占位符
	placeholders := make([]string, len(videoIDs))
	args := make([]interface{}, len(videoIDs)+1)
	args[0] = userID
	for i, id := range videoIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	// 构建查询语句
	query := fmt.Sprintf(`
		SELECT video_id, status 
		FROM user_anime_collections
		WHERE user_id = ? AND video_id IN (%s)
	`, strings.Join(placeholders, ","))

	// 执行查询
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 构建结果映射
	result := make(map[int]bool)
	// 默认所有视频ID都未收藏
	for _, id := range videoIDs {
		result[id] = false
	}

	// 读取查询结果
	for rows.Next() {
		var videoID int
		var status int
		if err := rows.Scan(&videoID, &status); err != nil {
			return nil, err
		}
		// status为1表示已收藏
		result[videoID] = status == 1
	}

	return &result, nil
}

func (r *VideoRepositoryImpl) GetAnimeCollectionByUserAndVideoID(ctx context.Context, userID int, videoID int) (bool, error) {
	query := `
		SELECT status 
		FROM user_anime_collections
		WHERE user_id = ? AND video_id = ?
	`

	var status int
	err := r.db.QueryRowContext(ctx, query, userID, videoID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到记录，表示用户未收藏该视频
			return false, nil
		}
		return false, err
	}

	// status为1表示已收藏
	return status == 1, nil
}

func (r *VideoRepositoryImpl) GetAnimeGenres(ctx context.Context, videoID int) ([]string, error) {
	query := `
		SELECT genre
		FROM anime_genres
		WHERE anime_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []string
	for rows.Next() {
		var genre string
		err := rows.Scan(&genre)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, nil
}
