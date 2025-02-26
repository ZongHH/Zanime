// Package bootstrap 实现了应用程序的依赖注入容器
// 负责管理和协调所有核心组件之间的依赖关系
// 采用分层架构，从基础设施层开始逐步构建直至接口层
package bootstrap

import (
	"crawler/internal/application/crawler"
	"crawler/internal/application/search"
	"crawler/internal/domain/repository"
	"crawler/internal/domain/service"
	"crawler/internal/grpc/scrapeService"
	"crawler/internal/infrastructure/collector"
	"crawler/internal/infrastructure/config"
	"crawler/internal/infrastructure/database"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Container 应用程序的主容器结构体
// 统一管理所有层级的组件实例，实现依赖注入
// 包括：
// - 接口层：对外提供的各种接口服务
// - 服务层：核心业务逻辑的实现
// - 仓储层：数据持久化的抽象
// - 基础层：提供底层基础设施支持
type Container struct {
	Interfaces  *Interfaces    // 接口层组件集合
	Services    *services      // 服务层组件集合
	Repositorys *repositotys   // 仓储层组件集合
	Bases       *bases         // 基础设施层组件集合
	Cfg         *config.Config // 应用配置信息
}

// BuildContainer 构建完整的依赖注入容器
// 按照从下到上的顺序初始化各层组件:
// 1. 加载配置
// 2. 初始化基础设施
// 3. 构建仓储层
// 4. 装配服务层
// 5. 配置接口层
// 如果在构建过程中发生错误，将立即终止程序运行
func BuildContainer(configPath string) *Container {
	// 加载应用配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化基础设施
	bases := initBases(cfg)

	// 初始化仓储层
	repositorys := initRepositotys(bases)

	// 初始化服务层
	services := initServices(cfg, bases, repositorys)

	// 初始化接口层
	interfaces := initInterfaces(services)

	return &Container{
		Interfaces:  interfaces,
		Services:    services,
		Repositorys: repositorys,
		Bases:       bases,
		Cfg:         cfg,
	}
}

// Interfaces 接口层组件集合
// 包含所有对外提供服务的接口实现
type Interfaces struct {
	GRPC *grpc.Server // gRPC服务器实例
	Lis  net.Listener // 网络监听器
}

// initInterfaces 初始化接口层组件
// 配置和启动所有对外提供的服务接口
// 目前支持的接口类型：
// - gRPC服务：用于视频资源的远程调用
func initInterfaces(services *services) *Interfaces {
	server := grpc.NewServer()
	scrapeService.RegisterVideoServer(server, &scrapeService.Server{Searcher: services.Search})
	return &Interfaces{
		GRPC: server,
	}
}

// services 服务层组件集合
// 实现核心业务逻辑，处理业务规则和流程
type services struct {
	Crawler service.CrawlerService // 爬取静态动漫资源服务
	Search  service.SearchService  // 搜索动漫视频服务
	Scrape  service.ScrapeService  // 爬取视频链接服务
}

// initServices 初始化所有业务服务层实例
// - 创建各种业务服务的具体实现
// - 注入所需的配置和依赖
// - 确保服务之间的协作正常
func initServices(cfg *config.Config, bases *bases, repositorys *repositotys) *services {
	scrape := search.NewVideoScraper(nil)
	return &services{
		Crawler: crawler.NewAnimeCrawler(cfg, bases.CollectorPool, repositorys.VideoRepo),
		Search:  search.NewAnimeSearcher(cfg, bases.CollectorPool, scrape, repositorys.VideoRepo),
		Scrape:  scrape,
	}
}

// repositorys 仓储层组件集合
// 处理数据的持久化和访问逻辑
type repositotys struct {
	VideoRepo repository.VideoRepository // 视频资源的仓储接口实现
}

// initRepositotys 初始化所有仓储层实例
// - 配置数据访问层
// - 创建仓储接口的具体实现
// - 注入数据库连接等基础设施
func initRepositotys(bases *bases) *repositotys {
	return &repositotys{
		VideoRepo: database.NewVideoRepository(bases.DB.GetDB()),
	}
}

// bases 基础设施层组件集合
// 提供底层技术支持和基础服务
type bases struct {
	DB            *database.DB             // 数据库连接管理器
	CollectorPool *collector.CollectorPool // 爬虫收集器连接池
}

// initBases 初始化所有基础设施层实例
// - 建立数据库连接
// - 初始化爬虫收集器池
// - 配置其他基础服务
// 任何初始化错误都被视为致命错误，将导致程序终止
func initBases(cfg *config.Config) *bases {
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v\n", err)
	}
	collectorPool := collector.NewCollectorPool(cfg)
	return &bases{
		DB:            db,
		CollectorPool: collectorPool,
	}
}

// close 关闭所有基础设施连接
// - 清理数据库连接
// - 释放爬虫收集器资源
// - 确保所有资源被正确释放
// 在应用程序退出时调用此方法
func (b *bases) close() {
	b.DB.Close()
}
