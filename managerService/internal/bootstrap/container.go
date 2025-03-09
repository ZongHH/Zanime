package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	appService "managerService/internal/application/service"
	"managerService/internal/domain/repository"
	"managerService/internal/domain/service"
	"managerService/internal/infrastructure/config"
	"managerService/internal/infrastructure/database"
	"managerService/internal/interfaces/handler"
	"managerService/internal/interfaces/router"
)

type Container struct {
	Bases       *bases
	Repositorys *repositorys
	Services    *services
	Handlers    *handlers
	Controller  *router.Controller
}

func LoadContainer(configPath string) *Container {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载初始化配置失败: %v", err)
	}

	bases, err := initBases(cfg)
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v", err)
	}

	repositorys, err := initRepositorys(bases)
	if err != nil {
		log.Fatalf("初始化仓储失败: %v", err)
	}

	services, err := initServices(repositorys)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	handlers, err := initHandlers(services)
	if err != nil {
		log.Fatalf("初始化处理器失败: %v", err)
	}

	controller := router.NewController(cfg)

	return &Container{
		Bases:       bases,
		Repositorys: repositorys,
		Services:    services,
		Handlers:    handlers,
		Controller:  controller,
	}
}

type handlers struct {
	StatisticsHandler     *handler.StatisticsHandler
	UserActionLogsHandler *handler.UserActionLogsHandler
}

func initHandlers(services *services) (*handlers, error) {
	return &handlers{
		StatisticsHandler:     handler.NewStatisticsHandler(services.StatisticsService),
		UserActionLogsHandler: handler.NewUserActionLogsHandler(services.UserActionLogsService),
	}, nil
}

type services struct {
	StatisticsService     service.StatisticsDataService
	UserActionLogsService service.UserActionLogsService
}

func initServices(repositorys *repositorys) (*services, error) {
	return &services{
		StatisticsService:     appService.NewStatisticsServiceImpl(repositorys.StatisticsRepository),
		UserActionLogsService: appService.NewUserActionLogsServiceImpl(repositorys.UserActionLogRepository),
	}, nil
}

type repositorys struct {
	StatisticsRepository    repository.StatisticsRepository
	UserActionLogRepository repository.UserActionLogsRepository
}

func initRepositorys(bases *bases) (*repositorys, error) {
	return &repositorys{
		StatisticsRepository:    database.NewStatisticsRepository(bases.DB),
		UserActionLogRepository: database.NewUserActionLogRepository(bases.DB),
	}, nil
}

type bases struct {
	DB *sql.DB
}

func initBases(cfg *config.Config) (*bases, error) {
	db, err := database.NewMySQLDB(&database.Config{
		Username: cfg.MySQL.Username,
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})

	if err != nil {
		return nil, fmt.Errorf("初始化数据库连接失败: %v", err)
	}

	return &bases{
		DB: db,
	}, nil
}
