package entity

// Video 表示视频实体
type Video struct {
	// 数据库原生字段
	ID            int    `json:"video_id"`
	Name          string `json:"video_name"`
	ReleaseDate   string `json:"release_date"`
	Area          string `json:"area"`
	Description   string `json:"description"`
	CoverImageUrl string `json:"cover_image_url"`
	Views         int    `json:"views"`
	Likes         int    `json:"likes"`
	CreatedAt     string `json:"created_at"`
	UploaderID    int    `json:"uploader_id"`

	// 额外字段
	Genres   string   `json:"genres"`
	Episodes []string `json:"episodes"`
	Rating   string   `json:"rating"`
	VideoUrl string   `json:"video_url"`

	// 筛选选项
	Initial string `json:"initial"` // 首字母

	// 收藏时间
	CollectedAt string `json:"collected_at"`
}

// VideoFilters 表示视频筛选选项
type VideoFilters struct {
	Regions []string `json:"regions"` // 可用地区
	Years   []string `json:"years"`   // 可用年份
	Types   []string `json:"types"`   // 视频类型
	Letters []string `json:"letters"` // 字母筛选
}

// UserAnimeCollection 表示用户的动漫收藏记录
type UserAnimeCollection struct {
	ID        int64  `json:"id"`
	UserID    int    `json:"user_id"`
	VideoID   int    `json:"video_id"`
	Status    int8   `json:"status"` // 收藏状态: 1-已收藏 0-取消收藏
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
