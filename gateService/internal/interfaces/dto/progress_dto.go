package dto

import "gateService/internal/domain/entity"

type LoadProgressRequest struct {
	UserID  int
	VideoID int    `form:"videoId"`
	Episode string `form:"episode"`
}

type LoadProgressResponse struct {
	Code     int              `json:"code"`
	Progress *entity.Progress `json:"progress"`
}

type SaveProgressRequest struct {
	UserID        int    `json:"user_id"`
	VideoID       int    `json:"video_id"`
	Episode       string `json:"episode"`
	Progress      int    `json:"progress"`
	Area          string `json:"area"`
	VideoName     string `json:"video_name"`
	Release       string `json:"release"`
	Genre         string `json:"genre"`
	CoverImageURL string `json:"cover_image_url"`
	MsgType       string `json:"msg_type"`
}

type SaveProgressResponse struct {
	Code int `json:"code"`
}

type WatchHistoryRequest struct {
	UserID   int `form:"user_id"`
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type WatchHistoryResponse struct {
	Code     int                `json:"code"`
	Progress []*entity.Progress `json:"progress"`
}
