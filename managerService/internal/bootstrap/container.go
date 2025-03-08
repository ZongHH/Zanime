package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	"managerService/internal/infrastructure/config"
	"managerService/internal/infrastructure/database"
)

type Container struct {
	Bases *bases
}

func LoadContainer(configPath string) *Container {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载初始化配置失败: %v", err)
	}

	bases, err := InitBases(cfg)
	if err != nil {
		log.Fatalf("初始化基础设施失败: %v", err)
	}

	return &Container{
		Bases: bases,
	}
}

type bases struct {
	DB *sql.DB
}

func InitBases(cfg *config.Config) (*bases, error) {
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
