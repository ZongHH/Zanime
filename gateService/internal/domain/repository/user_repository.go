package repository

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
	"time"
)

// UserRepository 定义了用户数据访问层的接口
type UserRepository interface {
	// BeginTx 开始事务
	// 参数:
	// - ctx: 上下文
	// 返回:
	// - *sql.Tx: 事务
	// - error: 错误信息
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// CreateUser 创建新用户
	// 参数:
	// - ctx: 上下文
	// - user: 用户信息
	// 返回:
	// - int: 新创建用户的ID
	// - error: 错误信息
	CreateUser(ctx context.Context, user *entity.UserInfo) (int, error)

	// CreateUserWithTx 使用事务创建新用户
	// 参数:
	// - ctx: 上下文
	// - tx: 事务
	// - user: 用户信息
	// 返回:
	// - error: 错误信息
	CreateUserWithTx(ctx context.Context, tx *sql.Tx, user *entity.UserInfo) error

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
	// - tx: 事务
	// - userID: 要删除的用户ID
	// 返回:
	// - error: 错误信息
	DeleteUser(ctx context.Context, tx *sql.Tx, userID int) error

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

	// CreateUserNotification 创建用户通知
	// 参数:
	// - ctx: 上下文
	// - tx: 事务
	// - notification: 用户通知
	// 返回:
	// - error: 错误信息
	CreateUserNotification(ctx context.Context, tx *sql.Tx, notification *entity.UserNotification) error

	// GetUserNotifications 获取用户通知
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// - notificationType: 通知类型
	// - page: 页码
	// - pageSize: 每页数量
	// 返回:
	// - *[]entity.UserNotification: 用户通知列表
	// - error: 错误信息
	GetUserNotifications(ctx context.Context, userID int, notificationType int8, page int, pageSize int) (*[]entity.UserNotification, error)

	// CheckInRedis 检查键是否存在于Redis中
	// 参数:
	// - ctx: 上下文
	// - key: 键
	// 返回:
	// - bool: 键是否存在
	// - time.Duration: 剩余过期时间
	// - error: 错误信息
	CheckInRedis(ctx context.Context, key string) (bool, time.Duration, error)

	// SetInRedis 设置键在Redis中
	// 参数:
	// - ctx: 上下文
	// - key: 键
	// - value: 值
	// - ttl: 过期时间
	// 返回:
	// - error: 错误信息
	SetInRedis(ctx context.Context, key string, value int, ttl time.Duration) error
}
