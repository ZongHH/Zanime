package entity

// PostCategory 帖子分类实体
// 用于管理和组织帖子的分类信息
type PostCategory struct {
	CategoryID int    `json:"category_id"` // 分类ID,唯一标识
	Name       string `json:"name"`        // 分类名称
	Icon       string `json:"icon"`        // 分类图标URL
	PostCount  int    `json:"post_count"`  // 该分类下的帖子数量
	SortOrder  int    `json:"sort_order"`  // 排序顺序,数字越小越靠前
	Status     int8   `json:"status"`      // 分类状态:0-禁用,1-启用
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
}

// PostTag 帖子标签实体
// 用于标记和分类帖子的关键词标签
type PostTag struct {
	TagID     int    `json:"tag_id"`     // 标签ID,唯一标识
	Name      string `json:"name"`       // 标签名称
	PostCount int    `json:"post_count"` // 使用该标签的帖子数量
	CreatedAt string `json:"created_at"` // 创建时间
}

// Post 帖子实体
// 论坛的核心内容实体,包含帖子的基本信息
type Post struct {
	PostID        int64  `json:"post_id"`        // 帖子ID,唯一标识
	UserID        int    `json:"user_id"`        // 发帖用户ID
	CategoryID    int    `json:"category_id"`    // 所属分类ID
	Title         string `json:"title"`          // 帖子标题
	Content       string `json:"content"`        // 帖子内容
	ViewCount     int    `json:"view_count"`     // 浏览次数
	LikeCount     int    `json:"like_count"`     // 点赞数量
	CommentCount  int    `json:"comment_count"`  // 评论数量
	FavoriteCount int    `json:"favorite_count"` // 收藏数量
	IsPinned      bool   `json:"is_pinned"`      // 是否置顶
	IsFeatured    bool   `json:"is_featured"`    // 是否精华帖
	Status        int8   `json:"status"`         // 帖子状态:0-删除,1-正常,2-待审核
	CreatedAt     string `json:"created_at"`     // 创建时间
	UpdatedAt     string `json:"updated_at"`     // 更新时间

	// 额外字段
	Images []PostImage `json:"images"` // 帖子图片URL列表
}

// PostImage 帖子图片实体
// 存储帖子中包含的图片信息
type PostImage struct {
	ImageID   int64  `json:"image_id"`   // 图片ID,唯一标识
	PostID    int64  `json:"post_id"`    // 所属帖子ID
	ImageURL  string `json:"image_url"`  // 图片URL地址
	SortOrder int    `json:"sort_order"` // 图片排序顺序
	CreatedAt string `json:"created_at"` // 创建时间
}

// PostTagRelation 帖子标签关联实体
// 维护帖子和标签之间的多对多关系
type PostTagRelation struct {
	PostID    int64  `json:"post_id"`    // 帖子ID
	TagID     int    `json:"tag_id"`     // 标签ID
	CreatedAt string `json:"created_at"` // 创建时间
}

// PostComment 帖子评论实体
// 存储帖子的评论信息,支持多级评论
type PostComment struct {
	CommentID  int64  `json:"comment_id"`  // 评论ID,唯一标识
	PostID     int64  `json:"post_id"`     // 所属帖子ID
	UserID     int    `json:"user_id"`     // 评论用户ID
	ToUserID   *int   `json:"to_user_id"`  // 被回复的评论用户ID
	ParentID   *int64 `json:"parent_id"`   // 父评论ID,用于回复评论
	RootID     *int64 `json:"root_id"`     // 根评论ID,用于标识评论层级
	Content    string `json:"content"`     // 评论内容
	LikeCount  int    `json:"like_count"`  // 点赞数量
	ReplyCount int    `json:"reply_count"` // 回复数量
	Level      int    `json:"level"`       // 评论层级:1-一级评论,2-二级评论
	Status     int8   `json:"status"`      // 评论状态:0-删除,1-正常
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
}

// PostLike 帖子点赞实体
// 记录用户对帖子的点赞信息
type PostLike struct {
	UserID    int    `json:"user_id"`    // 用户ID
	PostID    int64  `json:"post_id"`    // 帖子ID
	Status    int8   `json:"status"`     // 点赞状态:0-取消点赞,1-已点赞
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// PostFavorite 帖子收藏实体
// 记录用户收藏的帖子信息
type PostFavorite struct {
	UserID    int    `json:"user_id"`    // 用户ID
	PostID    int64  `json:"post_id"`    // 帖子ID
	Status    int8   `json:"status"`     // 收藏状态:0-取消收藏,1-已收藏
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// PostNotification 帖子通知实体
// 存储与帖子相关的用户通知信息
type PostNotification struct {
	NotificationID int64  `json:"notification_id"` // 通知ID,唯一标识
	UserID         int    `json:"user_id"`         // 接收通知的用户ID
	FromUserID     int    `json:"from_user_id"`    // 触发通知的用户ID
	PostID         *int64 `json:"post_id"`         // 相关帖子ID
	CommentID      *int64 `json:"comment_id"`      // 相关评论ID
	Type           int8   `json:"type"`            // 通知类型:1-点赞,2-评论,3-回复,4-系统通知
	Content        string `json:"content"`         // 通知内容
	IsRead         bool   `json:"is_read"`         // 是否已读
	CreatedAt      string `json:"created_at"`      // 创建时间
}

// PostCommentLike 帖子评论点赞实体
// 记录用户对评论的点赞信息
type PostCommentLike struct {
	UserID    int    `json:"user_id"`    // 点赞用户ID
	CommentID int64  `json:"comment_id"` // 被点赞的评论ID
	Status    int8   `json:"status"`     // 点赞状态:0-取消点赞,1-已点赞
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}
