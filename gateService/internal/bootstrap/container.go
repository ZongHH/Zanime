package bootstrap

import (
	"gateService/internal/infrastructure/config"
	"log"
)

// Container 包含所有服务依赖
type Container struct {
	Config       *config.Config
	Bases        *bases
	Repositories *repositories
	Services     *services
	Consumers    *consumers
	Interfaces   *interfaces
}

// NewContainer 创建并初始化容器
func NewContainer(configPath string) *Container {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化基础设施
	bases := initBases(cfg)

	// 初始化仓储层
	repositories := initRepositories(bases)

	// 初始化服务层
	services := initServices(cfg, bases, repositories)

	// 初始化消费者
	consumers := initConsumers(bases, repositories)

	// 初始化接口层
	interfaces := initInterfaces(cfg, bases, services)

	return &Container{
		Config:       cfg,
		Bases:        bases,
		Repositories: repositories,
		Services:     services,
		Consumers:    consumers,
		Interfaces:   interfaces,
	}
}
