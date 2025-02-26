// package database 提供了数据库访问层的实现
package database

import (
	"context"
	"database/sql"
	"fmt"
	"gateService/internal/domain/entity"
	"strings"
)

// UserRepositoryImpl 实现了用户仓储接口
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepositoryImpl 创建一个新的用户仓储实现实例
// 参数:
// - db: 数据库连接对象
// 返回:
// - *UserRepositoryImpl: 用户仓储实现实例
func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// CreateUser 创建新用户
// 参数:
// - ctx: 上下文
// - user: 用户信息
// 返回:
// - int: 新创建用户的ID
// - error: 错误信息
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *entity.UserInfo) (int, error) {
	query := "INSERT INTO user_infos (username, email, password, avatar_url) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.AvatarURL)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

// IsExistUser 检查用户是否存在
// 参数:
// - ctx: 上下文
// - email: 用户邮箱
// 返回:
// - bool: 用户是否存在
// - error: 错误信息
func (r *UserRepositoryImpl) IsExistUser(ctx context.Context, email string) (bool, error) {
	query := "SELECT email FROM user_infos WHERE email = ?"
	row := r.db.QueryRowContext(ctx, query, email)
	var existUser entity.UserInfo
	err := row.Scan(&existUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// VerifyUser 验证用户信息
// 参数:
// - ctx: 上下文
// - user: 待验证的用户信息
// 返回:
// - bool: 验证是否通过
// - error: 错误信息
func (r *UserRepositoryImpl) VerifyUser(ctx context.Context, user *entity.UserInfo) (bool, error) {
	query := "SELECT user_id, username, avatar_url, full_name, gender, birth_date FROM user_infos WHERE email = ? AND password = ?"
	row := r.db.QueryRowContext(ctx, query, user.Email, user.Password)
	err := row.Scan(&user.UserID, &user.Username, &user.AvatarURL, &user.FullName, &user.Gender, &user.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetUserByEmail 通过邮箱获取用户信息
// 参数:
// - ctx: 上下文
// - email: 用户邮箱
// 返回:
// - *entity.UserInfo: 用户信息
// - error: 错误信息
func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entity.UserInfo, error) {
	query := "SELECT * FROM user_infos WHERE email = ?"
	row := r.db.QueryRowContext(ctx, query, email)
	var user entity.UserInfo
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.Gender, &user.BirthDate, &user.CreatedAt, &user.LastLoginAt, &user.AvatarURL)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 通过用户ID获取用户信息
// 参数:
// - ctx: 上下文
// - userID: 用户ID
// 返回:
// - *entity.UserInfo: 用户信息
// - error: 错误信息
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID int) (*entity.UserInfo, error) {
	query := "SELECT * FROM user_infos WHERE user_id = ?"
	row := r.db.QueryRowContext(ctx, query, userID)
	var user entity.UserInfo
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.Signature, &user.FullName, &user.Gender, &user.BirthDate, &user.CreatedAt, &user.LastLoginAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
// 参数:
// - ctx: 上下文
// - user: 更新后的用户信息
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entity.UserInfo) error {
	query := "UPDATE user_infos SET username = ?, gender = ?, signature = ?, avatar_url = ? WHERE user_id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Gender, user.Signature, user.AvatarURL, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser 删除用户
// 参数:
// - ctx: 上下文
// - userID: 要删除的用户ID
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, userID int) error {
	query := "DELETE FROM user_infos WHERE user_id = ?"
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetUsersByIDs 通过用户ID列表批量查询用户信息
// 参数:
// - ctx: 上下文
// - userIDs: 用户ID列表指针
// 返回:
// - *[]entity.UserInfo: 用户信息列表指针,按传入的userIDs顺序返回
// - error: 错误信息
func (r *UserRepositoryImpl) GetUsersByIDs(ctx context.Context, userIDs *[]int) (*[]entity.UserInfo, error) {
	// 检查userIDs是否为空
	if len(*userIDs) == 0 {
		return &[]entity.UserInfo{}, nil
	}

	// 构建IN查询的占位符
	placeholders := make([]string, len(*userIDs))
	args := make([]interface{}, len(*userIDs))
	for i, id := range *userIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT user_id, username, avatar_url FROM user_infos WHERE user_id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建一个map用于存储查询结果,key为user_id
	userMap := make(map[int]entity.UserInfo)
	for rows.Next() {
		user := entity.UserInfo{}
		err := rows.Scan(&user.UserID, &user.Username, &user.AvatarURL)
		if err != nil {
			return nil, err
		}
		userMap[user.UserID] = user
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 按照传入的userIDs顺序构建结果
	users := make([]entity.UserInfo, len(*userIDs))
	for i, id := range *userIDs {
		if user, ok := userMap[id]; ok {
			users[i] = user
		}
	}

	return &users, nil
}

// UpdateUserAvatar 更新用户头像URL
// 参数:
// - ctx: 上下文
// - userID: 用户ID
// - avatarURL: 头像URL
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) UpdateUserAvatar(ctx context.Context, userID int, avatarURL string) error {
	query := "UPDATE user_infos SET avatar_url = ? WHERE user_id = ?"
	result, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}
