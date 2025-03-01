package dto

import (
	"gateService/internal/domain/entity"
	"mime/multipart"
)

// LoginRequest 用户登录请求参数
type LoginRequest struct {
	Email    string `form:"email" binding:"required"`    // 用户邮箱,必填
	Password string `form:"password" binding:"required"` // 用户密码,必填
	UserIP   string `form:"user_ip"`                     // 用户IP地址
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	Code     int              `json:"code"`      // 响应状态码,200表示成功
	UserInfo *entity.UserInfo `json:"user_info"` // 用户信息
}

// RegisterRequest 用户注册请求参数
type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`           // 用户邮箱,必填且需符合邮箱格式
	Password string `form:"password" binding:"required,min=8,max=32"` // 用户密码,必填且长度在8-32位之间
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	Code    int    `json:"code"`    // 响应状态码,200表示成功
	Message string `json:"message"` // 响应消息
}

// VerifyUserRequest 用户验证请求参数
type VerifyUserRequest struct {
	UserID int `form:"user_id" binding:"required"` // 用户ID,必填
}

// VerifyUserResponse 用户验证响应
type VerifyUserResponse struct {
	Code    int    `json:"code"`    // 响应状态码,200表示成功
	Message string `json:"message"` // 响应消息
}

// UserInfoRequest 获取用户信息请求参数
type UserInfoRequest struct {
	UserID int // 用户ID
}

// UserInfoResponse 获取用户信息响应
type UserInfoResponse struct {
	Code int      `json:"code"` // 响应状态码,200表示成功
	User UserInfo `json:"user"` // 用户基本信息
}

// GetUserProfileRequest 获取用户个人主页信息的请求参数
type GetUserProfileRequest struct {
	UserID int `form:"user_id"` // 用户ID
}

// GetUserProfileResponse 获取用户个人主页信息的响应
type GetUserProfileResponse struct {
	Code    int         `json:"code"`    // 响应状态码,200表示成功
	Profile *UserDetail `json:"profile"` // 用户详细信息
}

// UserDetail 用户详细信息
type UserDetail struct {
	UserID        int    `json:"user_id"`         // 用户唯一标识
	Username      string `json:"username"`        // 用户昵称
	Email         string `json:"email"`           // 邮箱
	AvatarURL     string `json:"avatar_url"`      // 头像URL
	IsVIP         bool   `json:"is_vip"`          // 是否是VIP用户
	VIPLevel      int    `json:"vip_level"`       // VIP等级
	VIPExpireDate string `json:"vip_expire_date"` // VIP到期时间
	RegisterTime  string `json:"register_time"`   // 注册时间
	Signature     string `json:"signature"`       // 个性签名
	Gender        string `json:"gender"`          // 性别
}

// UpdateUserProfileRequest 更新用户个人信息的请求参数
type UpdateUserProfileRequest struct {
	UserID    int    `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户昵称
	Email     string `json:"email"`      // 邮箱
	Gender    string `json:"gender"`     // 性别
	Signature string `json:"signature"`  // 个性签名
	AvatarURL string `json:"avatar_url"` // 头像URL
}

// UpdateUserProfileResponse 更新用户个人信息的响应
type UpdateUserProfileResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
}

// UserStatsRequest 用户个人主页计数信息的请求参数
type UserStatsRequest struct {
	UserID int `form:"user_id"` // 用户ID
}

// UserStatsResponse 用户个人主页计数信息的响应
type UserStatsResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
	Data struct {
		FollowingCount    int `json:"following_count"`     // 关注数量
		PostCount         int `json:"post_count"`          // 发布的帖子数量
		FavoritePostCount int `json:"favorite_post_count"` // 收藏的帖子数量
		CommentCount      int `json:"comment_count"`       // 评论数量
	} `json:"data"`
}

// UploadAvatarResponse 上传头像的响应
type UploadAvatarResponse struct {
	Code int    `json:"code"` // 响应状态码,200表示成功
	URL  string `json:"url"`  // 头像文件的访问URL
}

// UploadAvatarRequest 上传头像的请求参数
type UploadAvatarRequest struct {
	UserID int                   `form:"user_id"` // 用户ID
	Avatar *multipart.FileHeader `form:"avatar"`  // 头像文件,使用multipart/form-data格式上传
}

// UserNotificationRequest 获取用户通知的请求参数
type UserNotificationRequest struct {
	UserID   int  `form:"user_id"`                                     // 用户ID
	Type     int8 `form:"type" binding:"required,oneof=0 1 2 3 4 5 6"` // 通知类型,限制有效值
	Page     int  `form:"page" binding:"required,min=1"`               // 页码,必须大于等于1
	PageSize int  `form:"page_size" binding:"omitempty,min=1,max=50"`  // 每页数量,限制范围1-50
}

// UserNotificationResponse 获取用户通知的响应
type UserNotificationResponse struct {
	Code          int                `json:"code"`          // 响应状态码,200表示成功
	Notifications []UserNotification `json:"notifications"` // 通知数据列表
}

// UserNotification 用户通知数据结构
type UserNotification struct {
	ID       int64  `json:"id"`       // 通知ID
	Avatar   string `json:"avatar"`   // 头像URL
	Username string `json:"username"` // 用户名
	Message  string `json:"message"`  // 通知消息内容
	Type     int8   `json:"type"`     // 通知类型
	Time     string `json:"time"`     // 通知时间
}
