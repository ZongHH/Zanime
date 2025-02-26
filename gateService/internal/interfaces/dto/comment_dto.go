package dto

// SubmitCommentRequest 提交评论请求结构体
// 包含用户提交评论所需的基本信息
type SubmitCommentRequest struct {
	UserID  int    `json:"user_id"`  // 评论用户ID
	VideoID int    `json:"video_id"` // 被评论视频ID
	Content string `json:"content"`  // 评论内容
}

// SubmitCommentResponse 提交评论响应结构体
type SubmitCommentResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
}

// GetCommentsRequest 获取评论列表请求结构体
type GetCommentsRequest struct {
	VideoID  int `form:"video_id"`  // 视频ID
	Page     int `form:"page"`      // 当前页码,从1开始
	PageSize int `form:"page_size"` // 每页评论数量
}

// GetCommentsResponse 获取评论列表响应结构体
type GetCommentsResponse struct {
	Code      int        `json:"code"`       // 响应状态码,200表示成功
	TotalPage int        `json:"total_page"` // 总页数
	Comments  []*Comment `json:"comments"`   // 评论列表
}

// Comment 评论信息结构体
// 包含评论的详细信息以及对应的回复列表
type Comment struct {
	CommentID int      `json:"comment_id"`        // 评论ID
	Content   string   `json:"content"`           // 评论内容
	UserInfo  UserInfo `json:"userInfo"`          // 评论用户信息
	CreatedAt string   `json:"created_at"`        // 评论创建时间
	ReplyNum  int      `json:"reply_num"`         // 回复数量
	Replies   []*Reply `json:"replies,omitempty"` // 回复列表,可选字段
}

// Reply 回复信息结构体
// 包含对评论的回复信息
type Reply struct {
	CommentID   int      `json:"comment_id"`   // 回复ID
	RootID      int      `json:"root_id"`      // 根评论ID
	Content     string   `json:"content"`      // 回复内容
	UserInfo    UserInfo `json:"userInfo"`     // 回复用户信息
	CreatedAt   string   `json:"created_at"`   // 回复创建时间
	RepliedName string   `json:"replied_name"` // 被回复用户的用户名
}

// SubmitReplyRequest 提交回复请求结构体
type SubmitReplyRequest struct {
	UserID   int    `json:"user_id"`                       // 回复用户ID
	VideoID  int    `json:"video_id" binding:"required"`   // 被回复视频ID
	Content  string `json:"content" binding:"required"`    // 回复内容
	ParentID int    `json:"parent_id" binding:"required"`  // 父评论ID
	RootID   int    `json:"root_id" binding:"required"`    // 根评论ID
	ToUserID int    `json:"to_user_id" binding:"required"` // 被回复用户ID
}

// SubmitReplyResponse 提交回复响应结构体
type SubmitReplyResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
}

// GetReplyRequest 获取回复请求结构体
type GetReplyRequest struct {
	RootID   int `form:"root_id"`   // 根评论ID
	Page     int `form:"page"`      // 当前页码,从1开始
	PageSize int `form:"page_size"` // 每页回复数量
}

// GetReplyResponse 获取回复响应结构体
type GetReplyResponse struct {
	Code    int      `json:"code"`    // 响应状态码,200表示成功
	Replies []*Reply `json:"replies"` // 回复列表
}
