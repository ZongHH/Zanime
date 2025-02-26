package entity

// UserInfo 用户信息结构体
// 对应数据库表 user_infos
type UserInfo struct {
	UserID      int    `json:"user_id"`       // 用户ID,自增主键
	Username    string `json:"username"`      // 用户名,不可为空
	Email       string `json:"email"`         // 邮箱,唯一索引,不可为空
	Password    string `json:"password"`      // 密码,不可为空
	AvatarURL   string `json:"avatar_url"`    // 头像URL,不可为空
	Signature   string `json:"signature"`     // 个性签名,默认"这个人很懒，什么都没留下"
	FullName    string `json:"full_name"`     // 用户全名
	Gender      string `json:"gender"`        // 性别:Male/Female/Other,默认Other
	BirthDate   string `json:"birth_date"`    // 出生日期,格式:YYYY-MM-DD,默认2000-01-01
	CreatedAt   string `json:"created_at"`    // 创建时间,自动生成
	LastLoginAt string `json:"last_login_at"` // 最后登录时间,自动更新
}
