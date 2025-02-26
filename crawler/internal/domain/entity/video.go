package entity

// Video 视频实体
// 用于存储视频的基本信息
type Video struct {
	VideoID       int    `json:"video_id"`        // 视频唯一标识ID
	VideoName     string `json:"video_name"`      // 视频名称
	ReleaseDate   string `json:"release_date"`    // 发布日期
	Area          string `json:"area"`            // 地区/产地
	Description   string `json:"description"`     // 视频描述/简介
	CoverImageUrl string `json:"cover_image_url"` // 封面图片URL
	Views         int    `json:"views"`           // 观看次数
	Likes         int    `json:"likes"`           // 点赞数
	CreatedAt     string `json:"created_at"`      // 创建时间
	UploaderID    int    `json:"uploader_id"`     // 上传者ID
}

// AnimeGenre 动漫分类实体
// 用于存储动漫的分类信息,实现动漫与分类的多对多关系
type AnimeGenre struct {
	ID        int    `json:"id"`         // 分类记录唯一标识ID
	Genre     string `json:"genre"`      // 分类名称
	AnimeID   int    `json:"anime_id"`   // 关联的动漫ID
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
	DeletedAt string `json:"deleted_at"` // 删除时间(软删除)
}

// VideoUrl 视频链接实体
// 用于存储视频的播放地址信息,支持一个视频多个播放源
type VideoUrl struct {
	ID        int    `json:"id"`         // 视频链接记录唯一标识ID
	VideoID   int    `json:"video_id"`   // 关联的视频ID
	Episode   string `json:"episode"`    // 集数信息
	VideoUrl  string `json:"video_url"`  // 视频播放地址
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
	Status    int    `json:"status"`     // 链接状态:0-无效,1-有效
}
