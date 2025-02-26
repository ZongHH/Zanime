package database

import (
	"context"
	"crawler/internal/domain/entity"
	"database/sql"
	"fmt"
	"log"
)

type VideoRepositoryImpl struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepositoryImpl {
	return &VideoRepositoryImpl{db: db}
}

func (v *VideoRepositoryImpl) CreateVideo(ctx context.Context, anime *entity.Video) error {
	// 参数验证
	if err := validateAnime(anime); err != nil {
		return err
	}

	// video_name, release_date, area联合唯一索引
	// SQL 语句，使用 ON DUPLICATE KEY UPDATE 处理重复情况
	query := `
        INSERT INTO anime_videos 
        (video_name, release_date, area, description, cover_image_url, uploader_id)
        VALUES (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            video_id = video_id
    `

	// 执行插入或更新
	result, err := v.db.ExecContext(ctx, query,
		anime.VideoName,
		anime.ReleaseDate,
		anime.Area,
		anime.Description,
		anime.CoverImageUrl,
		anime.UploaderID,
	)

	if err != nil {
		return fmt.Errorf("创建/更新视频失败: %v", err)
	}

	// 获取ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取插入ID失败: %v", err)
	}

	// 如果是更新操作（id = 0），需要查询现有记录的ID
	if id == 0 {
		err = v.db.QueryRowContext(ctx, `
            SELECT video_id 
            FROM anime_videos 
            WHERE video_name = ? AND release_date = ? AND area = ?`,

			anime.VideoName, anime.ReleaseDate, anime.Area,
		).Scan(&anime.VideoID)

		if err != nil {
			return fmt.Errorf("获取现有记录ID失败: %v", err)
		}
	} else {
		anime.VideoID = int(id)
	}

	return nil
}

// 验证函数
func validateAnime(anime *entity.Video) error {
	if anime == nil {
		return fmt.Errorf("anime 对象为空")
	}
	if anime.VideoName == "" {
		return fmt.Errorf("视频名称不能为空")
	}
	if anime.ReleaseDate == "" {
		return fmt.Errorf("发布日期不能为空")
	}
	if anime.Area == "" {
		return fmt.Errorf("地区不能为空")
	}
	if anime.CoverImageUrl == "" {
		return fmt.Errorf("封面图片URL不能为空")
	}
	if anime.UploaderID == 0 {
		return fmt.Errorf("上传者ID不能为空")
	}
	return nil
}

func (v *VideoRepositoryImpl) CreateGenre(ctx context.Context, genres *[]entity.AnimeGenre) error {
	if len(*genres) == 0 {
		return fmt.Errorf("类型列表为空")
	}

	// genre, anime_id联合唯一索引
	// 构建批量插入的SQL
	query := `
        INSERT INTO anime_genres 
        (genre, anime_id) 
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE
            id = id
    `

	// 开始事务
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer tx.Rollback() // 如果提交成功，这个回滚将不会执行

	// 准备预处理语句
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("预处理语句失败: %v", err)
	}
	defer stmt.Close()

	// 遍历所有类型进行插入
	for _, genre := range *genres {
		// 验证数据
		if genre.Genre == "" {
			return fmt.Errorf("类型名称不能为空")
		}
		if genre.AnimeID == 0 {
			return fmt.Errorf("动漫ID不能为空")
		}

		// 执行插入
		_, err := stmt.Exec(
			genre.Genre,
			genre.AnimeID,
		)

		if err != nil {
			return fmt.Errorf("插入类型失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

func (v *VideoRepositoryImpl) CreateVideoUrl(ctx context.Context, videoUrls *[]entity.VideoUrl) error {
	if len(*videoUrls) == 0 {
		return fmt.Errorf("视频URL列表为空")
	}

	// video_id, episode联合唯一索引
	// 准备SQL语句
	query := `
        INSERT INTO video_urls 
        (video_id, episode, video_url) 
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE 
        	video_url = VALUES(video_url)
    `

	// 开始事务
	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer tx.Rollback() // 如果提交成功，这个回滚将不会执行

	// 准备预处理语句
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("预处理语句失败: %v", err)
	}
	defer stmt.Close()

	// 记录成功和失败的数量
	successCount := 0
	skipCount := 0

	// 批量插入数据
	for _, url := range *videoUrls {
		// 检查必要字段
		if url.VideoID == 0 || url.Episode == "" || url.VideoUrl == "" {
			skipCount++
			log.Printf("跳过无效数据: VideoID=%d, Episode=%s, VideoURL=%s",
				url.VideoID, url.Episode, url.VideoUrl)
			continue // 跳过这条记录，继续处理下一条
		}

		// 执行插入
		_, err = stmt.Exec(url.VideoID, url.Episode, url.VideoUrl)
		if err != nil {
			log.Printf("插入数据失败: %v", err)
			skipCount++
			continue // 跳过这条记录，继续处理下一条
		}
		successCount++
	}

	// 如果没有成功插入任何记录
	if successCount == 0 {
		return fmt.Errorf("没有成功插入任何记录，跳过 %d 条无效数据", skipCount)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 返回处理结果
	if skipCount > 0 {
		log.Printf("成功插入 %d 条记录，跳过 %d 条无效数据", successCount, skipCount)
	}

	return nil
}

// GetVideoIDByVideoName 方法用于根据视频名称获取视频的 ID
func (v *VideoRepositoryImpl) GetVideoIDByVideoName(ctx context.Context, videoName string) (int, error) {
	var videoID int
	query := "SELECT video_id FROM anime_videos WHERE video_name = ?"
	err := v.db.QueryRowContext(ctx, query, videoName).Scan(&videoID)
	if err != nil {
		return 0, err
	}
	return videoID, nil
}

// GetVideoURLByVideoIDANDEpisode 方法用于根据视频 ID 和集数获取视频的链接
func (v *VideoRepositoryImpl) GetVideoURLByVideoIDANDEpisode(ctx context.Context, videoID int, Episode string) (string, error) {
	var videoUrl string
	query := "SELECT video_url FROM video_urls WHERE video_id = ? AND episode = ?"
	err := v.db.QueryRowContext(ctx, query, videoID, Episode).Scan(&videoUrl)
	if err != nil {
		return "", err
	}
	return videoUrl, nil
}
