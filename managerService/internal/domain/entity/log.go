package entity

// UserActionLog 用户行为日志结构，对应数据库表user_action_logs
type UserActionLog struct {
	ID        int64  `json:"id"`
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Module    string `json:"module"`
	IPAddress string `json:"ip_address"`
	Status    int8   `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`

	// 用户名
	UserName string `json:"user_name"`
	// 用户类型
	UserType string `json:"user_type"`
}
