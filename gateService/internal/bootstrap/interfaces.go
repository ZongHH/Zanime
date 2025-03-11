package bootstrap

import (
	"gateService/internal/grpc/server/tokenService"
	"gateService/internal/infrastructure/config"
	"gateService/internal/interfaces/http/router"
	"log"
	"net"

	"google.golang.org/grpc"
)

// interfaces 聚合应用的所有接口组件
// 包含 HTTP 路由、gRPC 服务实例和网络监听器
type interfaces struct {
	Router *router.Controller // HTTP 路由控制器
	GRPC   *grpc.Server       // gRPC 服务实例
	Lis    net.Listener       // 网络监听器
}

// initInterfaces 初始化应用的所有接口组件
// 参数:
//   - cfg: 应用配置信息
//   - bases: 基础组件（JWT/Cookie管理器等）
//   - services: 领域服务实例
//
// 返回:
//   - 初始化完成的接口组件集合
func initInterfaces(cfg *config.Config, bases *bases, services *services) *interfaces {
	// 初始化 HTTP 路由控制器，注入所有依赖服务
	router := router.NewController(cfg, bases.JwtManager, bases.CookieManager, bases.RateLimit,
		services.ProgressService, services.PostService, services.CommentService,
		services.SearchService, services.UserService, services.ProductService,
		services.OrderService, services.VideoService, services.WebSocketService)

	// 创建 gRPC 服务器并注册 Token 服务
	grpcServer := grpc.NewServer()
	tokenService.RegisterTokenServer(grpcServer, services.TokenService)

	return &interfaces{
		Router: router,
		GRPC:   grpcServer,
	}
}

// Start 启动所有接口服务
// 参数:
//   - gRpcAddr: gRPC 服务监听地址
func (i *interfaces) Start(gRpcAddr string) {
	i.Router.Run()      // 启动 HTTP 服务
	i.gRpcRun(gRpcAddr) // 启动 gRPC 服务
}

// Close 优雅关闭所有接口服务
func (i *interfaces) Close() {
	i.Router.Stop() // 停止 HTTP 服务
	i.gRpcStop()    // 停止 gRPC 服务
}

// gRpcRun 启动 gRPC 服务监听
// 参数:
//   - gRpcAddr: gRPC 服务监听地址
func (i *interfaces) gRpcRun(gRpcAddr string) {
	var err error
	// 创建 TCP 监听器
	i.Lis, err = net.Listen("tcp", gRpcAddr)
	if err != nil {
		log.Fatalf("监听 gRpc 端口失败: %v", err)
	}

	// 在协程中启动 gRPC 服务
	go func() {
		if err := i.GRPC.Serve(i.Lis); err != nil {
			log.Fatalf("启动 gRpc 服务失败: %v", err)
		}
	}()
}

// gRpcStop 优雅停止 gRPC 服务
func (i *interfaces) gRpcStop() {
	i.GRPC.GracefulStop() // 优雅关闭 gRPC 服务（等待现有请求完成）
	i.Lis.Close()         // 关闭监听器
	log.Println("关闭 gRpc 服务成功")
}
