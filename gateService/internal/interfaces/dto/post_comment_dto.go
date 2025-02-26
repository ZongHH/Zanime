package dto

// SubmitCommentRequest 提交评论的请求参数
type SubmitPostCommentRequest struct {
	PostID    int64  `json:"post_id"`    // 帖子ID
	CommentID int64  `json:"comment_id"` // 评论ID
	UserID    int    `json:"user_id"`    // 评论用户ID
	Content   string `json:"content"`    // 评论内容
	CreatedAt string `json:"created_at"` // 评论创建时间
}

// SubmitCommentResponse 提交评论的响应
type SubmitPostCommentResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
}

type SubmitPostReplyRequest struct {
	PostID    int64  `json:"post_id"`    // 帖子ID
	CommentID int64  `json:"comment_id"` // 评论ID
	UserID    int    `json:"user_id"`    // 评论用户ID
	Content   string `json:"content"`    // 评论内容
	ParentID  *int64 `json:"parent_id"`  // 父评论ID
	RootID    *int64 `json:"root_id"`    // 根评论ID
	ToUserID  *int   `json:"to_user_id"` // 被回复的评论用户ID
	CreatedAt string `json:"created_at"` // 评论创建时间
}

type SubmitPostReplyResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功
}

// GetCommentsRequest 获取评论列表的请求参数
type GetPostCommentsRequest struct {
	UserID   int   `form:"user_id"`   // 用户ID
	ID       int64 `form:"id"`        // 视频/文章ID
	Page     int   `form:"page"`      // 页码,从1开始
	PageSize int   `form:"page_size"` // 每页评论数量
	RootID   int64 `form:"root_id"`   // 根评论ID,用于获取子评论列表
}

// GetPostCommentsResponse 获取评论列表的响应
type GetPostCommentsResponse struct {
	Code      int            `json:"code"`       // 响应状态码,200表示成功
	TotalPage int            `json:"total_page"` // 总页数,用于分页
	Comments  []*CommentItem `json:"comments"`   // 评论列表数据
}

// CommentItem 评论详情
// 包含评论的基本信息、作者信息、点赞信息以及子回复列表
type CommentItem struct {
	ID        int64       `json:"id"`         // 评论唯一标识
	Content   string      `json:"content"`    // 评论文本内容
	Author    UserInfo    `json:"author"`     // 评论作者的用户信息
	CreatedAt string      `json:"created_at"` // 评论创建时间,格式为RFC3339
	LikeCount int         `json:"like_count"` // 评论获得的点赞总数
	IsLiked   bool        `json:"is_liked"`   // 当前登录用户是否对该评论点赞
	Replies   []ReplyItem `json:"replies"`    // 该评论下的回复列表
	ReplyNum  int         `json:"reply_num"`  // 该评论下的回复总数
}

// ReplyItem 回复详情
// 包含回复的基本信息、作者信息和点赞信息
type ReplyItem struct {
	ID        int64    `json:"id"`         // 回复唯一标识
	Content   string   `json:"content"`    // 回复文本内容
	Author    UserInfo `json:"author"`     // 回复作者的用户信息
	CreatedAt string   `json:"created_at"` // 回复创建时间,格式为RFC3339
	LikeCount int      `json:"like_count"` // 回复获得的点赞总数
	IsLiked   bool     `json:"is_liked"`   // 当前登录用户是否对该回复点赞
	ReplyTo   UserInfo `json:"reply_to"`   // 回复的回复
}

// UserInfo 用户信息
// 包含展示用户所需的基本信息
type UserInfo struct {
	ID        int    `json:"id"`         // 用户唯一标识
	Username  string `json:"username"`   // 用户昵称/用户名
	AvatarURL string `json:"avatar_url"` // 用户头像的URL地址
}

// CommentLikeRequest 评论点赞请求参数
// 包含用户ID、评论ID和点赞状态
type CommentLikeRequest struct {
	UserID    int   // 用户ID,用于标识点赞用户
	CommentID int64 `json:"comment_id"` // 评论ID,用于标识被点赞的评论
	Status    bool  `json:"status"`     // 点赞状态,true表示点赞,false表示取消点赞
}

// CommentLikeResponse 评论点赞响应
// 包含响应状态码
type CommentLikeResponse struct {
	Code int `json:"code"` // 响应状态码,200表示成功,其他值表示失败
}
