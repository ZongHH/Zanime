// package service 提供了与用户相关的业务逻辑服务
package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// UserService 定义了用户服务的接口
// 提供用户注册、登录、登出和验证等基本功能
type UserService interface {
	// Register 用户注册
	// 参数:
	// - ctx: 上下文信息
	// - user: 注册请求数据,包含用户名、密码等注册信息
	// 返回:
	// - *dto.RegisterResponse: 注册响应数据,包含用户ID等信息
	// - error: 注册过程中的错误信息
	Register(ctx context.Context, user *dto.RegisterRequest) (*dto.RegisterResponse, error)

	// Login 用户登录
	// 参数:
	// - ctx: 上下文信息
	// - user: 登录请求数据,包含用户名和密码
	// 返回:
	// - *dto.LoginResponse: 登录响应数据,包含token等信息
	// - error: 登录过程中的错误信息
	Login(ctx context.Context, user *dto.LoginRequest) (*dto.LoginResponse, error)

	// Logout 用户登出
	// 参数:
	// - ctx: 上下文信息
	// 返回:
	// - error: 登出过程中的错误信息
	Logout(ctx context.Context) error

	// VerifyUser 验证用户信息
	// 参数:
	// - ctx: 上下文信息
	// - user: 验证请求数据,包含需要验证的用户信息
	// 返回:
	// - *dto.VerifyUserResponse: 验证响应数据,包含验证结果
	// - error: 验证过程中的错误信息
	VerifyUser(ctx context.Context, user *dto.VerifyUserRequest) (*dto.VerifyUserResponse, error)

	// GetUserInfo 获取用户信息
	// 参数:
	// - ctx: 上下文信息
	// - user: 用户信息请求参数,包含用户ID
	// 返回:
	// - *dto.UserInfoResponse: 用户信息响应数据,包含用户基本信息
	// - error: 获取用户信息过程中的错误信息
	GetUserInfo(ctx context.Context, user *dto.UserInfoRequest) (*dto.UserInfoResponse, error)

	// GetUserProfile 获取用户详细信息
	// 参数:
	// - ctx: 上下文信息
	// - user: 用户详细信息请求参数,包含用户ID
	// 返回:
	// - *dto.GetUserProfileResponse: 用户详细信息响应数据,包含用户VIP状态、注册时间等详细信息
	// - error: 获取用户详细信息过程中的错误信息
	GetUserProfile(ctx context.Context, user *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error)

	// UpdateUserProfile 更新用户个人信息
	// 参数:
	// - ctx: 上下文信息
	// - user: 更新用户个人信息请求参数,包含用户ID、用户名、邮箱、性别、个性签名、头像URL
	// 返回:
	// - *dto.UpdateUserProfileResponse: 更新用户个人信息响应数据,包含更新结果
	// - error: 更新用户个人信息过程中的错误信息
	UpdateUserProfile(ctx context.Context, user *dto.UpdateUserProfileRequest) (*dto.UpdateUserProfileResponse, error)

	// GetUserStats 获取用户个人主页计数信息
	// 参数:
	// - ctx: 上下文信息
	// - user: 用户个人主页计数信息请求参数,包含用户ID
	// 返回:
	// - *dto.UserStatsResponse: 用户个人主页计数信息响应数据,包含用户关注数量、帖子数量、收藏帖子数量、评论数量
	// - error: 获取用户个人主页计数信息过程中的错误信息
	GetUserStats(ctx context.Context, user *dto.UserStatsRequest) (*dto.UserStatsResponse, error)

	// UploadAvatar 上传用户头像
	// 参数:
	// - ctx: 上下文信息
	// - user: 上传用户头像请求参数,包含用户ID和头像文件
	// 返回:
	// - *dto.UploadAvatarResponse: 上传用户头像响应数据,包含头像文件的访问URL
	// - error: 上传用户头像过程中的错误信息
	UploadAvatar(ctx context.Context, user *dto.UploadAvatarRequest) (*dto.UploadAvatarResponse, error)

	// GetUserNotifications 获取用户通知
	// 参数:
	// - ctx: 上下文信息
	// - user: 获取用户通知请求参数,包含用户ID和通知类型
	// 返回:
	// - *dto.UserNotificationResponse: 获取用户通知响应数据,包含通知列表
	// - error: 获取用户通知过程中的错误信息
	GetUserNotifications(ctx context.Context, user *dto.UserNotificationRequest) (*dto.UserNotificationResponse, error)

	// GetTestAccount 获取体验账号
	// 参数:
	// - ctx: 上下文信息
	// - user: 获取体验账号请求参数,包含用户IP地址
	// 返回:
	// - *dto.TestAccountResponse: 获取体验账号响应数据,包含体验账号邮箱和密码
	// - error: 获取体验账号过程中的错误信息
	GetTestAccount(ctx context.Context, user *dto.TestAccountRequest) (*dto.TestAccountResponse, error)
}
