package bootstrap

import (
	"gateService/internal/application/connection"
	serviceImpl "gateService/internal/application/service"
	"gateService/internal/domain/service"
	"gateService/internal/grpc/server/tokenService"
	"gateService/internal/infrastructure/config"
)

// services 结构体聚合所有业务领域服务实例
// 职责：
//   - 作为依赖注入容器，集中管理各领域服务的生命周期
//   - 提供清晰的领域服务依赖关系视图
type services struct {
	// UserService 用户领域服务
	// 功能包含：用户注册/登录、信息管理、权限验证等
	UserService service.UserService

	// PostService 社区帖子领域服务
	// 功能包含：帖子CRUD、标签管理、评论互动等
	PostService service.PostService

	// ProgressService 观看进度领域服务
	// 功能包含：进度记录/同步、历史记录管理、跨设备同步等
	ProgressService service.ProgressService

	// CommentService 视频评论领域服务
	// 功能包含：评论发布、树状结构管理、敏感词过滤等
	CommentService service.CommentService

	// SearchService 搜索领域服务
	// 功能包含：多条件检索、相关性排序、搜索建议生成等
	SearchService service.SearchService

	// ProductService 商品领域服务
	// 功能包含：商品信息管理、库存校验、价格策略等
	ProductService service.ProductService

	// OrderService 订单领域服务
	// 功能包含：订单创建、支付流程、状态机管理、超时取消等
	OrderService service.OrderService

	// VideoService 视频领域服务
	// 功能包含：元数据管理、推荐算法集成、资源地址生成等
	VideoService service.VideoService

	// TokenService 认证令牌gRPC服务
	// 功能包含：JWT令牌签发/验证、令牌刷新、吊销列表管理等
	TokenService *tokenService.Server

	// WebSocketService 实时通信领域服务
	// 功能包含：长连接管理、消息路由、连接状态维护等
	WebSocketService *connection.WebSocketServiceImpl
}

// initServices 服务初始化工厂方法
// 职责：
//   - 根据基础设施和仓储层组件构造领域服务实例
//   - 显式声明各服务的依赖关系，实现依赖注入
//
// 参数:
//   - cfg: 配置文件
//   - bases: 基础设施组件集合（数据库连接、中间件管理器等）
//   - repos: 数据仓储层实例集合（各领域Repository实现）
//
// 返回:
//   - *services: 完全初始化的领域服务集合
func initServices(cfg *config.Config, bases *bases, repos *repositories) *services {
	return &services{
		UserService: serviceImpl.NewUserServiceImpl(
			&cfg.Storage,          // 文件存储配置
			repos.UserRepo,        // 用户数据仓储
			repos.PostRepo,        // 帖子数据仓储
			repos.PostCommentRepo, // 帖子评论数据仓储
			bases.JwtManager,      // JWT认证组件
			bases.CookieManager,   // Cookie管理组件
			bases.ProducerPool,    // 消息队列生产者池
		),
		PostService: serviceImpl.NewPostServiceImpl(
			&cfg.Storage,              // 文件存储配置
			repos.PostRepo,            // 帖子主数据仓储
			repos.PostTagRepo,         // 标签定义仓储
			repos.PostTagRelationRepo, // 标签关系仓储
			repos.PostCommentRepo,     // 帖子评论仓储
			repos.UserRepo,            // 用户信息仓储
			bases.ProducerPool,        // NSQ消息生产者池
		),
		ProgressService: serviceImpl.NewProgressServiceImpl(
			bases.ProducerPool, // 消息队列生产者池（用于进度同步）
			repos.ProgressRepo, // 进度数据仓储
		),
		CommentService: serviceImpl.NewCommentServiceImpl(
			repos.CommentRepo, // 视频评论仓储
			repos.UserRepo,    // 用户信息仓储
		),
		SearchService: serviceImpl.NewSearchServiceImpl(
			repos.VideoRepo, // 视频元数据仓储
		),
		ProductService: serviceImpl.NewProductServiceImpl(
			repos.ProductRepo, // 商品数据仓储
		),
		OrderService: serviceImpl.NewOrderServiceImpl(
			repos.OrderRepo,    // 订单数据仓储
			repos.ProductRepo,  // 商品信息仓储
			bases.ProducerPool, // 消息队列生产者（订单事件通知）
			bases.RDB.GetRDB(), // Redis客户端（分布式锁）
			nil,                // 支付服务客户端（预留扩展）
		),
		VideoService: serviceImpl.NewVideoServiceImpl(
			bases.RDB.GetRDB(),    // Redis客户端（缓存层）
			bases.ScrapeClient,    // 爬虫服务客户端
			bases.RecommendClient, // 推荐算法服务客户端
			repos.VideoRepo,       // 视频元数据仓储
			repos.ProgressRepo,    // 进度数据仓储（关联查询）
		),
		TokenService: tokenService.NewServer(
			bases.JwtManager, // JWT管理器（签名/验证）
		),
		WebSocketService: connection.NewWebSocketServiceImpl(
			bases.WebSocketManager, // WebSocket连接管理器
		),
	}
}
