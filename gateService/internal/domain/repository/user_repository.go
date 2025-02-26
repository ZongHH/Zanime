package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

// UserRepository 定义了用户数据访问层的接口
type UserRepository interface {
	// CreateUser 创建新用户
	// 参数:
	// - ctx: 上下文
	// - user: 用户信息
	// 返回:
	// - int: 新创建用户的ID
	// - error: 错误信息
	CreateUser(ctx context.Context, user *entity.UserInfo) (int, error)

	// GetUserByEmail 通过邮箱查询用户信息
	// 参数:
	// - ctx: 上下文
	// - email: 用户邮箱
	// 返回:
	// - *entity.UserInfo: 用户信息
	// - error: 错误信息
	GetUserByEmail(ctx context.Context, email string) (*entity.UserInfo, error)

	// GetUserByID 通过用户ID查询用户信息
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// 返回:
	// - *entity.UserInfo: 用户信息
	// - error: 错误信息
	GetUserByID(ctx context.Context, userID int) (*entity.UserInfo, error)

	// UpdateUser 更新用户信息
	// 参数:
	// - ctx: 上下文
	// - user: 更新后的用户信息
	// 返回:
	// - error: 错误信息
	UpdateUser(ctx context.Context, user *entity.UserInfo) error

	// DeleteUser 删除用户
	// 参数:
	// - ctx: 上下文
	// - userID: 要删除的用户ID
	// 返回:
	// - error: 错误信息
	DeleteUser(ctx context.Context, userID int) error

	// IsExistUser 检查用户是否存在
	// 参数:
	// - ctx: 上下文
	// - email: 用户邮箱
	// 返回:
	// - bool: 用户是否存在
	// - error: 错误信息
	IsExistUser(ctx context.Context, email string) (bool, error)

	// VerifyUser 验证用户信息
	// 参数:
	// - ctx: 上下文
	// - user: 待验证的用户信息
	// 返回:
	// - bool: 验证是否通过
	// - error: 错误信息
	VerifyUser(ctx context.Context, user *entity.UserInfo) (bool, error)

	// GetUsersByIDs 通过用户ID列表批量查询用户信息
	// 参数:
	// - ctx: 上下文
	// - userIDs: 用户ID列表指针
	// 返回:
	// - *[]entity.UserInfo: 用户信息列表指针,按传入的userIDs顺序返回
	// - error: 错误信息
	GetUsersByIDs(ctx context.Context, userIDs *[]int) (*[]entity.UserInfo, error)

	// UpdateUserAvatar 更新用户头像URL
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// - avatarURL: 头像URL
	// 返回:
	// - error: 错误信息
	UpdateUserAvatar(ctx context.Context, userID int, avatarURL string) error
}
