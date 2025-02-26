package entity

type Progress struct {
	// 数据库原生字段
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	VideoID       int    `json:"video_id"`
	VideoName     string `json:"video_name"`
	CoverImageURL string `json:"cover_image_url"`
	Episode       string `json:"episode"`
	Progress      int    `json:"progress"`
	Status        string `json:"status"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`

	// 额外字段
	Area    string `json:"area"`
	Release string `json:"release"`
	Genre   string `json:"genre"`
}
