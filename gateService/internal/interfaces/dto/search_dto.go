package dto

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Query string `form:"query" binding:"required"` // 搜索关键词，必填
}

// SearchAnime 搜索结果中的动漫基本信息
type SearchAnime struct {
	VideoID  int    `json:"video_id"`        // 视频ID
	Title    string `json:"video_name"`      // 视频标题
	CoverUrl string `json:"cover_image_url"` // 封面图片URL
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Code   int            `json:"code"`   // 响应状态码
	Animes []*SearchAnime `json:"animes"` // 搜索结果动漫列表
}

// SearchDetailRequest 搜索详情请求参数
type SearchDetailRequest struct {
	UserID int    `form:"user_id"`                   // 用户ID
	Params string `form:"params" binding:"required"` // 搜索参数，必填
	Page   int    `form:"page" binding:"required"`   // 页码，必填
}

// SearchDetailAnime 搜索详情中的动漫完整信息
type SearchDetailAnime struct {
	VideoID     int      `json:"video_id"`        // 视频ID
	Title       string   `json:"video_name"`      // 视频标题
	CoverUrl    string   `json:"cover_image_url"` // 封面图片URL
	ReleaseDate string   `json:"release_date"`    // 发布日期
	Area        string   `json:"area"`            // 地区
	Description string   `json:"description"`     // 描述信息
	Genres      string   `json:"genres"`          // 动漫类型
	Episodes    []string `json:"episodes"`        // 剧集列表
	IsCollected bool     `json:"is_collected"`    // 是否已收藏
}

// SearchDetailResponse 搜索详情响应
type SearchDetailResponse struct {
	Code   int                  `json:"code"`   // 响应状态码
	Animes []*SearchDetailAnime `json:"animes"` // 搜索结果动漫详情列表
}
