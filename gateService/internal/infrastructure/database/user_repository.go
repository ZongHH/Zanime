// package database 提供了数据库访问层的实现
package database

import (
	"context"
	"database/sql"
	"fmt"
	"gateService/internal/domain/entity"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserRepositoryImpl 实现了用户仓储接口
type UserRepositoryImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewUserRepositoryImpl 创建一个新的用户仓储实现实例
// 参数:
// - db: 数据库连接对象
// - rdb: Redis连接对象
// 返回:
// - *UserRepositoryImpl: 用户仓储实现实例
func NewUserRepositoryImpl(db *sql.DB, rdb *redis.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:  db,
		rdb: rdb,
	}
}

// BeginTx 开始事务
// 参数:
// - ctx: 上下文
// 返回:
// - *sql.Tx: 事务
// - error: 错误信息
func (r *UserRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
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

// CreateUserWithTx 使用事务创建新用户
// 参数:
// - ctx: 上下文
// - tx: 事务
// - user: 用户信息
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) CreateUserWithTx(ctx context.Context, tx *sql.Tx, user *entity.UserInfo) error {
	query := "INSERT INTO user_infos (username, email, password, avatar_url) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.AvatarURL)
	if err != nil {
		return err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.UserID = int(userID)
	return nil
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
	query := "SELECT user_id, username, avatar_url, full_name, gender, birth_date FROM user_infos WHERE email = ? AND password = ? AND status = 1"
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
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.Signature, &user.FullName, &user.Gender, &user.BirthDate, &user.CreatedAt, &user.LastLoginAt, &user.Status)
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
// - tx: 事务
// - userID: 要删除的用户ID
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, tx *sql.Tx, userID int) error {
	query := "UPDATE user_infos SET status = 0 WHERE user_id = ?"
	_, err := tx.ExecContext(ctx, query, userID)
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
func (r *UserRepositoryImpl) GetUserNotifications(ctx context.Context, userID int, notificationType int8, page int, pageSize int) (*[]entity.UserNotification, error) {
	offset := (page - 1) * pageSize

	var query string
	var args []interface{}

	if notificationType == 0 {
		// 如果通知类型为0，不筛选通知类型
		query = `
			SELECT *
			FROM user_notifications 
			WHERE user_id = ? 
			ORDER BY created_at DESC 
			LIMIT ? OFFSET ?
			`
		args = []interface{}{userID, pageSize, offset}
	} else {
		// 否则按通知类型筛选
		query = `
			SELECT *
			FROM user_notifications 
			WHERE user_id = ? AND notification_type = ? 
			ORDER BY created_at DESC 
			LIMIT ? OFFSET ?
			`
		args = []interface{}{userID, notificationType, pageSize, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []entity.UserNotification{}
	for rows.Next() {
		notification := entity.UserNotification{}
		err := rows.Scan(
			&notification.NotificationID,
			&notification.UserID,
			&notification.FromUserID,
			&notification.PostID,
			&notification.CommentID,
			&notification.NotificationType,
			&notification.Content,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &notifications, nil
}

// CreateUserNotification 创建用户通知
// 参数:
// - ctx: 上下文
// - tx: 事务
// - notification: 用户通知
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) CreateUserNotification(ctx context.Context, tx *sql.Tx, notification *entity.UserNotification) error {
	// 动态构建字段列表和占位符
	fields := []string{"user_id", "from_user_id", "notification_type", "content"}
	placeholders := []string{"?", "?", "?", "?"}
	args := []interface{}{notification.UserID, notification.FromUserID, notification.NotificationType, notification.Content}

	// 根据是否有效ID添加post_id和comment_id
	if notification.PostID != nil {
		fields = append(fields, "post_id")
		placeholders = append(placeholders, "?")
		args = append(args, notification.PostID)
	}
	if notification.CommentID != nil {
		fields = append(fields, "comment_id")
		placeholders = append(placeholders, "?")
		args = append(args, notification.CommentID)
	}

	query := fmt.Sprintf(`
		INSERT INTO user_notifications (
			%s
		) VALUES (%s)
	`, strings.Join(fields, ", "), strings.Join(placeholders, ", "))

	_, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

// CheckInRedis 检查键是否存在于Redis中
// 参数:
// - ctx: 上下文
// - key: 键
// 返回:
// - bool: 键是否存在
// - time.Duration: 剩余过期时间
// - error: 错误信息
func (r *UserRepositoryImpl) CheckInRedis(ctx context.Context, key string) (bool, time.Duration, error) {
	// 使用Redis检查键是否存在
	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	// 如果key不存在,TTL返回-2
	// 如果key存在但没有设置过期时间,TTL返回-1
	if ttl == -2 {
		return false, 0, nil
	}

	return true, ttl, nil
}

// SetInRedis 设置键在Redis中
// 参数:
// - ctx: 上下文
// - key: 键
// - value: 值
// - ttl: 过期时间
// 返回:
// - error: 错误信息
func (r *UserRepositoryImpl) SetInRedis(ctx context.Context, key string, value int, ttl time.Duration) error {
	// 将键和值存入Redis,并设置过期时间
	err := r.rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
