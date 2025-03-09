package entity

// UserInfo 用户信息结构体
// 对应数据库表 user_infos
type UserInfo struct {
	UserID      int    `json:"user_id"`       // 用户ID,自增主键
	Username    string `json:"username"`      // 用户名,不可为空
	Email       string `json:"email"`         // 邮箱,唯一索引,不可为空
	Password    string `json:"password"`      // 密码,不可为空
	AccountType string `json:"account_type"`  // 账户类型:admin/regular/trial,默认regular
	Status      int8   `json:"status"`        // 用户状态:0-删除,1-正常,默认1
	FullName    string `json:"full_name"`     // 用户全名,可为空
	Gender      string `json:"gender"`        // 性别:Male/Female/Other,默认Other
	BirthDate   string `json:"birth_date"`    // 出生日期,格式:YYYY-MM-DD,默认2000-01-01
	AvatarURL   string `json:"avatar_url"`    // 头像URL,不可为空
	Signature   string `json:"signature"`     // 个性签名,默认"这个人很懒，什么都没留下"
	CreatedAt   string `json:"created_at"`    // 创建时间,自动生成,默认当前时间戳
	LastLoginAt string `json:"last_login_at"` // 最后登录时间,自动更新,默认当前时间戳
}

// UserNotification 用户通知结构体
// 对应数据库表 user_notifications
type UserNotification struct {
	NotificationID   int64  `json:"notification_id"`   // 通知ID,自增主键
	UserID           int    `json:"user_id"`           // 接收通知的用户ID
	FromUserID       int    `json:"from_user_id"`      // 发送通知的用户ID
	PostID           *int64 `json:"post_id"`           // 相关的帖子ID,可为空
	CommentID        *int64 `json:"comment_id"`        // 相关的评论ID,可为空
	NotificationType int8   `json:"notification_type"` // 通知类型: 1-点赞评论, 2-回复评论, 3-收藏帖子, 4-点赞帖子, 5-关注
	Content          string `json:"content"`           // 通知内容,可为空
	IsRead           bool   `json:"is_read"`           // 是否已读,默认为0(未读)
	CreatedAt        string `json:"created_at"`        // 创建时间,自动生成
}
