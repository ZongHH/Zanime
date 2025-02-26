package dto

type CreatePostRequest struct {
	UserID     int      `json:"user_id"`
	CategoryID int      `json:"category_id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Images     []string `json:"images"`
}

type CreatePostResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PostAuthor struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type PostTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PostImage struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type PostResponse struct {
	ID            int64       `json:"id"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	Author        PostAuthor  `json:"author"`
	CreatedAt     string      `json:"created_at"`
	Tags          []PostTag   `json:"tags"`
	IsPinned      bool        `json:"is_pinned"`
	IsFeatured    bool        `json:"is_featured"`
	ViewCount     int         `json:"view_count"`
	CommentCount  int         `json:"comment_count"`
	LikeCount     int         `json:"like_count"`
	FavoriteCount int         `json:"favorite_count"`
	IsLiked       bool        `json:"is_liked"`
	IsFavorited   bool        `json:"is_favorited"`
	Images        []PostImage `json:"images"`
}

type GetPostListResponse struct {
	Code        int            `json:"code"`
	Posts       []PostResponse `json:"posts"`
	PinnedPosts []PostResponse `json:"pinned_posts"`
}

type GetPostListRequest struct {
	UserID     int `form:"user_id"`
	CategoryID int `form:"category_id"`
	Page       int `form:"page"`
	PageSize   int `form:"page_size"`
}

type PostCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	PostCount int    `json:"post_count"`
}

type GetCategoryListResponse struct {
	Code       int             `json:"code"`
	Categories []*PostCategory `json:"categories"`
}

type GetCategoryListRequest struct {
}

type GetPostByPostIDRequest struct {
	UserID int   `form:"user_id"`
	PostID int64 `form:"post_id"`
}

type GetPostByPostIDResponse struct {
	Code int          `json:"code"`
	Post PostResponse `json:"post"`
}

type PostLikeRequest struct {
	PostID int64 `json:"post_id"`
	UserID int   `json:"user_id"`
	Status int8  `json:"status"`
}

type PostLikeResponse struct {
	Code int `json:"code"`
}

type PostFavoriteRequest struct {
	PostID int64 `json:"post_id"`
	UserID int   `json:"user_id"`
	Status int8  `json:"status"`
}

type PostFavoriteResponse struct {
	Code int `json:"code"`
}

type PostBriefInfo struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	ViewCount    int    `json:"view_count"`
	CommentCount int    `json:"comment_count"`
	LikeCount    int    `json:"like_count"`
	CreatedAt    string `json:"created_at"`
}

type GetRecentPostsRequest struct {
	UserID   int `form:"user_id"`
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type GetRecentPostsResponse struct {
	Code  int             `json:"code"`
	Posts []PostBriefInfo `json:"posts"`
}
