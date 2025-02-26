package entity

type Comment struct {
	CommentID int    `json:"comment_id"`
	VideoID   int    `json:"video_id"`
	UserID    int    `json:"user_id"`
	ToUserID  *int   `json:"to_user_id"`
	Content   string `json:"content"`
	RootID    *int   `json:"root_id"`
	ParentID  *int   `json:"parent_id"`
	ReplyNum  int    `json:"reply_num"`
	Level     int    `json:"level"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`

	// 额外字段,用于连表查询
	UserName  string `json:"user_name"`
	AvatarURL string `json:"avatar_url"`

	// 被回复用户名
	ToUserName string `json:"to_user_name"`
}
