package bootstrap

import (
	"gateService/internal/grpc/client/recommend"
	"gateService/internal/grpc/client/scrapeClient"
	"gateService/internal/infrastructure/config"
	"gateService/internal/infrastructure/database"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/infrastructure/middleware/websocket"
	"gateService/pkg/logger"
	"gateService/pkg/mq/nsqpool"
	"log"
)

// bases 结构体聚合应用程序的核心基础设施组件
// 职责：
//   - 作为基础设施层容器，统一管理各组件的生命周期
//   - 提供数据库、中间件、消息队列等基础设施的访问入口
//   - 确保各组件之间的依赖关系清晰可见
type bases struct {
	// DB 关系型数据库连接池（MySQL）
	// 功能包含：
	// - 执行SQL查询/事务
	// - 连接池管理
	// - 数据库健康检查
	DB *database.DB

	// RDB Redis数据库客户端
	// 功能包含：
	// - 缓存读写操作
	// - 分布式锁实现
	// - 会话状态存储
	RDB *database.RDB

	// JwtManager JWT令牌管理器
	// 功能包含：
	// - 令牌生成与签名验证
	// - 令牌过期时间管理
	// - 自定义声明解析
	JwtManager *auth.JWTManager

	// CookieManager Cookie安全管理器
	// 功能包含：
	// - Cookie加密/解密
	// - 防篡改验证
	// - 安全属性设置（HttpOnly/Secure等）
	CookieManager *auth.CookieManager

	// ProducerPool NSQ消息队列生产者池
	// 功能包含：
	// - 异步消息发布
	// - 连接复用管理
	// - 失败重试机制
	ProducerPool *nsqpool.ProducerPool

	// ScrapeClient 爬虫服务gRPC客户端池
	// 功能包含：
	// - 视频元数据抓取
	// - 资源地址解析
	// - 连接负载均衡
	ScrapeClient *scrapeClient.GRPCClientPool

	// RecommendClient 推荐服务gRPC客户端池
	// 功能包含：
	// - 个性化推荐算法调用
	// - 视频相似度计算
	// - 热点内容获取
	RecommendClient *recommend.GRPCClientPool

	// WebSocketManager WebSocket连接管理器
	// 功能包含：
	// - 连接生命周期管理
	// - 消息广播路由
	// - 心跳检测机制
	WebSocketManager *websocket.Manager
}

// initBases 基础设施初始化工厂方法
// 职责：
//   - 根据配置创建并配置所有基础设施组件
//   - 处理组件间的依赖关系
//   - 保证组件初始化顺序和错误处理
//
// 参数:
//   - cfg: 应用程序配置对象，包含各基础设施的连接参数
//
// 返回值:
//   - *bases: 完全初始化的基础设施组件集合
func initBases(cfg *config.Config) *bases {
	// 初始化NSQ消息生产者池（用于异步任务处理）
	// 配置参数：
	// - NSQDAddress: 消息队列服务地址
	// - PoolSize: 连接池容量（控制并发吞吐量）
	producerPool, err := nsqpool.NewProducerPool(&nsqpool.ProducerOptions{
		NSQDAddress: cfg.GetNSQDAddr(),
		PoolSize:    10,
	})
	if err != nil {
		log.Fatalf("初始化NSQ生产者池失败: %v\n", err)
	}

	// 初始化爬虫服务gRPC客户端池（元数据获取）
	scrapeClient, err := scrapeClient.NewGRPCClientPool(cfg)
	if err != nil {
		log.Fatalf("初始化爬虫服务客户端失败: %v\n", err)
	}

	// 初始化推荐服务gRPC客户端池（个性化推荐）
	recommendClient, err := recommend.NewGRPCClientPool(cfg)
	if err != nil {
		log.Fatalf("初始化推荐服务客户端失败: %v\n", err)
	}

	// 初始化WebSocket连接管理器（实时通信）
	websocketManager := websocket.NewManager(logger.Log)

	return &bases{
		DB:               database.NewDB(cfg),                // MySQL数据库连接（业务主存储）
		RDB:              database.NewRDB(cfg),               // Redis连接（缓存/会话）
		JwtManager:       auth.NewJWTManager(&cfg.JWT),       // JWT认证组件
		CookieManager:    auth.NewCookieManager(&cfg.Cookie), // Cookie安全组件
		ProducerPool:     producerPool,                       // 消息队列生产者
		ScrapeClient:     scrapeClient,                       // 爬虫服务客户端
		RecommendClient:  recommendClient,                    // 推荐服务客户端
		WebSocketManager: websocketManager,                   // WebSocket管理器
	}
}

// Close 安全关闭所有基础设施连接
// 执行顺序说明：
// 1. 先关闭数据库连接（保证数据持久化）
// 2. 再停止消息生产者（防止消息丢失）
// 3. 最后关闭gRPC连接池（释放网络资源）
func (b *bases) Close() {
	b.DB.Close()              // 关闭MySQL连接
	b.RDB.Close()             // 关闭Redis连接
	b.ProducerPool.Close()    // 停止NSQ消息生产者
	b.ScrapeClient.Close()    // 关闭爬虫服务gRPC连接池
	b.RecommendClient.Close() // 关闭推荐服务gRPC连接池
}
