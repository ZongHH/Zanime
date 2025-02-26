package router

import (
	"gateService/internal/interfaces/http/middleware/auth"
)

// setupAPIRoutes 初始化需要认证的API路由
// 包含以下功能模块的路由配置：
// - 评论管理
// - 视频服务
// - 订单处理
// - 搜索功能
// - 用户认证
// - 观看进度
// - 社区互动
// - WebSocket通信
func (c *Controller) setupAPIRoutes() {
	// 创建API路由组，所有路由需要JWT认证
	apiGroup := c.engine.Group("/api")
	apiGroup.Use(auth.JWTAuthMiddleware(c.jwtManager, c.cookieManager)) // 注入JWT和Cookie管理器

	// 路由分组注册
	{
		// ================== 评论管理模块 ==================
		// 功能：处理视频相关评论的提交和获取
		apiGroup.POST("/movie/submit", c.commentHandler.SubmitComment)    // 提交新评论（需包含：视频ID、评论内容、父评论ID）
		apiGroup.GET("/movie/comments", c.commentHandler.GetComments)     // 获取视频评论列表（参数：视频ID、分页信息）
		apiGroup.POST("/movie/submitReply", c.commentHandler.SubmitReply) // 提交回复（需包含：视频ID、回复内容、父评论ID、被回复用户ID）
		apiGroup.GET("/movie/replies", c.commentHandler.GetReply)         // 获取回复列表（参数：根评论ID、分页信息）

		// ================== 视频服务模块 ==================
		// 功能：提供视频资源访问和元数据查询
		apiGroup.GET("/video-resource", c.videoHandler.GetVideoURL)   // 获取视频播放地址（参数：视频ID、集数）
		apiGroup.GET("/video-info", c.videoHandler.GetVideoInfo)      // 获取视频详细信息（参数：视频ID）
		apiGroup.GET("/animeFilters", c.videoHandler.GetVideoFilters) // 获取动漫筛选条件（地区/年份/类型等）
		apiGroup.GET("/animeLibrary", c.videoHandler.GetVideoLibrary) // 获取动漫库列表（支持分页和条件过滤）
		apiGroup.GET("/getHomeAnime", c.videoHandler.GetHomeAnimes)   // 获取首页推荐动漫列表（根据用户ID）
		apiGroup.GET("/movie/recommend", c.videoHandler.GetRecommend) // 获取推荐动漫列表（根据当前动漫类型）

		// ================== 订单处理模块 ==================
		// 功能：处理商品购买和订单管理
		apiGroup.POST("/order", c.orderHandler.CreateOrder)           // 创建新订单（参数：商品ID、支付方式）
		apiGroup.GET("/get-products", c.productHandler.GetProducts)   // 获取可购商品列表
		apiGroup.POST("/call-pay", c.orderHandler.CallbackPay)        // 支付结果回调接口（第三方支付平台调用）
		apiGroup.GET("/get-orders", c.orderHandler.GetOrdersByUserID) // 获取用户历史订单（带分页）

		// ================== 搜索模块 ==================
		// 功能：提供动漫搜索服务
		apiGroup.GET("/search", c.searchHandler.SearchAnime)             // 动漫关键词搜索（参数：关键词、分页）
		apiGroup.GET("/searchDetail", c.searchHandler.SearchAnimeDetail) // 动漫详情搜索（参数：精确ID）

		// ================== 用户认证模块 ==================
		apiGroup.GET("/verify-token", c.userHandler.VerifyUser) // 令牌验证接口（返回最新用户信息）
		apiGroup.GET("/logout", c.userHandler.Logout)           // 用户登出（清除认证信息）

		// ================== 观看进度模块 ==================
		apiGroup.GET("/load-progress", c.progressHandler.LoadProgress)  // 加载观看进度（参数：视频ID）
		apiGroup.POST("/save-progress", c.progressHandler.SaveProgress) // 保存观看进度（参数：视频ID、时间点）
		apiGroup.GET("/watch-history", c.progressHandler.WatchHistory)  // 获取用户观看历史记录（参数：用户ID、页码、每页数量）

		// ================== 社区互动模块 ==================
		// 功能：处理用户发帖和评论互动
		apiGroup.POST("/create-post", c.postHandler.CreatePost)                 // 创建新帖子（参数：标题、内容、标签）
		apiGroup.GET("/post/comments", c.postHandler.GetPostComments)           // 获取帖子主评论列表（参数：帖子ID）
		apiGroup.GET("/comment/replies", c.postHandler.GetPostCommentsByRootID) // 获取评论回复（参数：根评论ID）
		apiGroup.POST("/comment/submit", c.postHandler.SubmitComment)           // 提交帖子评论（参数：帖子ID、内容）
		apiGroup.POST("/comment/reply", c.postHandler.SubmitPostReply)          // 提交帖子回复（参数：帖子ID、内容）
		apiGroup.POST("/comment/like", c.postHandler.CommentLike)               // 提交帖子点赞（参数：帖子ID、用户ID）
		apiGroup.GET("/post/categories", c.postHandler.GetPostCategoryList)     // 获取帖子分类列表
		apiGroup.GET("/post/list", c.postHandler.GetPostsByCategoryID)          // 获取帖子列表（参数：分类ID、页码、每页数量）
		apiGroup.GET("/post/detail", c.postHandler.GetPostByPostID)             // 获取帖子详情（参数：帖子ID）
		apiGroup.POST("/post/like", c.postHandler.PostLike)                     // 提交帖子点赞（参数：帖子ID、用户ID）
		apiGroup.POST("/post/favorite", c.postHandler.PostFavorite)             // 提交帖子收藏（参数：帖子ID、用户ID）
		apiGroup.GET("/post/recent", c.postHandler.RecentPosts)                 // 获取用户最近发布的帖子列表（参数：用户ID，默认返回最近5条）

		// ================== 用户信息模块 ==================
		apiGroup.GET("/user/current", c.userHandler.GetUserInfo)                       // 获取当前用户信息
		apiGroup.POST("/user/update-collection", c.videoHandler.UpdateAnimeCollection) // 更新动漫收藏状态（参数：视频ID、用户ID、收藏状态）
		apiGroup.GET("/user/collection", c.videoHandler.GetAnimeCollection)            // 获取动漫收藏列表（参数：用户ID、页码、每页数量）
		apiGroup.GET("/user/profile", c.userHandler.GetUserProfile)                    // 获取用户详细信息（参数：用户ID）
		apiGroup.POST("/user/update", c.userHandler.UpdateUserProfile)                 // 更新用户个人信息（参数：用户ID、用户名、邮箱、性别、个性签名、头像URL）
		apiGroup.GET("/user/stats", c.userHandler.GetUserStats)                        // 获取用户个人主页计数信息（参数：用户ID）
		apiGroup.POST("/user/upload-avatar", c.userHandler.UploadAvatar)               // 上传用户头像（参数：用户ID、头像文件）
	}

	// 创建连接路由组，所有路由需要JWT认证
	conGroup := c.engine.Group("/conn")
	conGroup.Use(auth.JWTAuthMiddleware(c.jwtManager, c.cookieManager))
	{
		// ================== WebSocket模块 ==================
		// 功能：建立实时通信连接
		conGroup.GET("/ws", c.websocketHandler.EstablishConnection) // WebSocket连接端点（协议升级）
	}
}
