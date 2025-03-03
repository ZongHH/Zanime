package router

import (
	"context"
	"gateService/internal/domain/repository"
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/config"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/http/handler"
	"gateService/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller 路由控制器，负责管理HTTP服务的核心组件和依赖注入
// 包含以下主要职责：
// 1. 维护全局配置和基础设施组件
// 2. 初始化各个业务模块的处理器
// 3. 管理HTTP服务器生命周期
type Controller struct {
	// 全局配置信息，包含：
	// - 服务器监听地址
	// - 超时设置
	// - 安全相关配置
	cfg *config.Config

	// 路由服务组件
	engine *gin.Engine  // Gin框架路由引擎，负责路由注册和中间件处理
	srv    *http.Server // HTTP服务器实例，管理服务启停

	// 基础设施组件
	jwtManager    *auth.JWTManager    // JWT令牌管理器，负责令牌的生成与验证
	cookieManager *auth.CookieManager // Cookie管理器，处理Cookie的加密和验证

	userRepository repository.UserRepository // 用户仓储实例

	// 业务模块处理器（接口处理层）
	progressHandler *handler.ProgressHandler // 用户观看进度处理器
	postHandler     *handler.PostHandler     // 帖子管理处理器
	commentHandler  *handler.CommentHandler  // 评论管理处理器
	searchHandler   *handler.SearchHandler   // 搜索功能处理器
	userHandler     *handler.UserHandler     // 用户管理处理器
	productHandler  *handler.ProductHandler  // 商品管理处理器
	orderHandler    *handler.OrderHandler    // 订单管理处理器
	videoHandler    *handler.VideoHandler    // 视频服务处理器

	// WebSocket通信处理器
	// 功能包括：
	// - 实时消息推送
	// - 长连接管理
	// - 双向通信支持
	websocketHandler *handler.WebSocketHandler
}

// NewController 构造函数，用于创建路由控制器实例
// 参数说明：
//   - cfg: 应用配置，包含所有运行时配置参数
//   - jwtManager: JWT认证管理器实例
//   - cookieManager: Cookie管理实例
//   - userRepository: 用户仓储实例
//   - progressService ~ websocketService: 各业务领域服务实现
//
// 返回值说明：
//   - *Controller: 初始化完成的路由控制器实例，包含所有依赖组件
func NewController(
	cfg *config.Config,
	jwtManager *auth.JWTManager,
	cookieManager *auth.CookieManager,
	userRepository repository.UserRepository,
	progressService service.ProgressService,
	postService service.PostService,
	commentService service.CommentService,
	searchService service.SearchService,
	userService service.UserService,
	productService service.ProductService,
	orderService service.OrderService,
	videoService service.VideoService,
	websocketService service.WebSocketService,
) *Controller {
	return &Controller{
		cfg:              cfg,
		engine:           gin.Default(), // 使用Gin默认配置初始化路由引擎
		jwtManager:       jwtManager,
		cookieManager:    cookieManager,
		userRepository:   userRepository,
		progressHandler:  handler.NewProgressHandler(progressService),   // 初始化进度处理器
		postHandler:      handler.NewPostHandler(postService),           // 初始化帖子处理器
		commentHandler:   handler.NewCommentHandler(commentService),     // 初始化评论处理器
		searchHandler:    handler.NewSearchHandler(searchService),       // 初始化搜索处理器
		userHandler:      handler.NewUserHandler(userService),           // 初始化用户处理器
		productHandler:   handler.NewProductHandler(productService),     // 初始化商品处理器
		orderHandler:     handler.NewOrderHandler(orderService),         // 初始化订单处理器
		videoHandler:     handler.NewVideoHandler(videoService),         // 初始化视频处理器
		websocketHandler: handler.NewWebSocketHandler(websocketService), // 初始化WebSocket处理器
	}
}

func (c *Controller) Run() {
	c.initMiddlewares()
	c.setupRoutes()

	c.srv = &http.Server{
		Addr:           c.cfg.GetHTTPAddr(),
		Handler:        c.engine,
		ReadTimeout:    c.cfg.HTTP.ReadTimeout,
		WriteTimeout:   c.cfg.HTTP.WriteTimeout,
		IdleTimeout:    c.cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: c.cfg.HTTP.MaxHeaderBytes,
	}

	go func() {
		if err := c.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动 Gin 服务失败: %v\n", err)
		}
	}()
}

func (c *Controller) Stop() {
	ctxWithTime, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.srv.Shutdown(ctxWithTime); err != nil {
		log.Fatalf("关闭 Gin 服务失败: %v\n", err)
	}
	log.Println("关闭 Gin 服务成功")
}

func (c *Controller) initMiddlewares() {
	// 安全中间件
	// c.setupMiddleHandler(c.middleHandler.IpLimiter())
	c.setupMiddleHandler(errors.ErrorHandler())
	// c.setupMiddleHandler(security.XssProtection())
	// c.setupMiddleHandler(security.CsrfProtection())
	// c.setupMiddleHandler(security.SizeLimiter())
	// c.setupMiddleHandler(security.TimeoutMiddleware(10 * time.Second))
	// c.setupMiddleHandler(security.SensitiveFilter())

}

func (c *Controller) setupRoutes() {
	c.setupAuthRoutes()
	c.setupAPIRoutes()
}

func (c *Controller) setupMiddleHandler(f func(*gin.Context)) {
	c.engine.Use(f)
}
