package dto

import "gateService/internal/domain/entity"

// GetVideoInfoRequest 获取视频信息的请求参数
type GetVideoInfoRequest struct {
	VideoID int `form:"videoId"` // 视频ID
}

// GetVideoInfoResponse 获取视频信息的响应
type GetVideoInfoResponse struct {
	Code      int           `json:"code"`       // 响应状态码
	VideoInfo *entity.Video `json:"video_info"` // 视频详细信息
}

// GetVideoLibraryRequest 获取视频库的请求参数
type GetVideoLibraryRequest struct {
	Region   string `form:"region"`    // 地区
	Year     string `form:"year"`      // 年份
	Type     string `form:"type"`      // 类型
	Letter   string `form:"letter"`    // 首字母
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页数量
}

// GetVideoLibraryRespnse 获取视频库的响应
type GetVideoLibraryRespnse struct {
	Code   int             `json:"code"`   // 响应状态码
	Total  int             `json:"total"`  // 总数量
	Videos []*entity.Video `json:"videos"` // 视频列表
}

// GetVideoFiltersRequest 获取视频过滤条件的请求参数
type GetVideoFiltersRequest struct {
}

// GetVideoFiltersResponse 获取视频过滤条件的响应
type GetVideoFiltersResponse struct {
	Code         int                  `json:"code"`          // 响应状态码
	VideoFilters *entity.VideoFilters `json:"video_filters"` // 视频过滤条件
}

// GetVideoURLRequest 获取视频URL的请求参数
type GetVideoURLRequest struct {
	UserID  int    // 用户ID
	VideoID int    `form:"videoId"` // 视频ID
	Episode string `form:"episode"` // 集数
}

// GetVideoURLResponse 获取视频URL的响应
type GetVideoURLResponse struct {
	Code       int                 `json:"code"`       // 响应状态码
	VideoFiles []map[string]string `json:"VideoFiles"` // 视频文件信息
	PosterPath string              `json:"PosterPath"` // 海报路径
}

// GetHomeAnimesRequest 获取首页动漫的请求参数
type GetHomeAnimesRequest struct {
	UserID int // 用户ID
}

// GetHomeAnimesResponse 获取首页动漫的响应
type GetHomeAnimesResponse struct {
	Code         int             `json:"code"`          // 响应状态码
	HomeAnime    *entity.Video   `json:"home_anime"`    // 首页推荐动漫
	JapanAnime   []*entity.Video `json:"japan_anime"`   // 日本动漫列表
	ChinaAnime   []*entity.Video `json:"china_anime"`   // 中国动漫列表
	WesternAnime []*entity.Video `json:"western_anime"` // 欧美动漫列表
	TopAnime     []*entity.Video `json:"top_anime"`     // 热门动漫列表
	AnimeGenres  []string        `json:"anime_genres"`  // 动漫类型列表
}

// UpdateAnimeCollectionRequest 更新动漫收藏的请求参数
type UpdateAnimeCollectionRequest struct {
	UserID  int  `json:"user_id"`  // 用户ID
	VideoID int  `json:"video_id"` // 视频ID
	Status  bool `json:"status"`   // 收藏状态
}

// UpdateAnimeCollectionResponse 更新动漫收藏的响应
type UpdateAnimeCollectionResponse struct {
	Code int `json:"code"` // 响应状态码
}

// GetAnimeCollectionRequest 获取动漫收藏的请求参数
type GetAnimeCollectionRequest struct {
	UserID int `form:"user_id"` // 用户ID
	Page   int `form:"page"`    // 页码
	Limit  int `form:"limit"`   // 每页数量
}

type AnimeCollection struct {
	VideoID       int    `json:"video_id"`        // 视频ID
	Title         string `json:"title"`           // 视频标题
	CoverImageUrl string `json:"cover_image_url"` // 封面图片URL
	CollectedAt   string `json:"collected_at"`    // 收藏时间
}

// GetAnimeCollectionResponse 获取动漫收藏的响应
type GetAnimeCollectionResponse struct {
	Code            int                `json:"code"`             // 响应状态码
	Total           int                `json:"total"`            // 总数量
	AnimeCollection []*AnimeCollection `json:"anime_collection"` // 动漫收藏列表
}

// GetRecommendRequest 获取相关动漫推荐的请求参数
type GetRecommendRequest struct {
	UserID  int `form:"user_id"`  // 用户ID
	VideoID int `form:"video_id"` // 视频ID
}

// GetRecommendResponse 获取相关动漫推荐的响应
type GetRecommendResponse struct {
	Code            int                 `json:"code"`            // 响应状态码
	Recommendations []*RecommendedAnime `json:"recommendations"` // 推荐动漫列表
}

// RecommendedAnime 推荐动漫信息
type RecommendedAnime struct {
	ID       int    `json:"id"`       // 视频ID
	Title    string `json:"title"`    // 视频标题
	CoverUrl string `json:"coverUrl"` // 封面图片URL
	Rating   string `json:"rating"`   // 评分
}
