package entity

// Anime 表示动画视频实体
type Anime struct {
	VideoID       int    `json:"video_id"`        // 视频ID，主键
	VideoName     string `json:"video_name"`      // 视频名称
	ReleaseDate   string `json:"release_date"`    // 发布日期
	Area          string `json:"area"`            // 地区
	Description   string `json:"description"`     // 描述
	CoverImageURL string `json:"cover_image_url"` // 封面图片URL
	Views         int    `json:"views"`           // 观看次数
	Likes         int    `json:"likes"`           // 点赞数
	CreatedAt     string `json:"created_at"`      // 创建时间
	UploaderID    int    `json:"uploader_id"`     // 上传者ID
}
