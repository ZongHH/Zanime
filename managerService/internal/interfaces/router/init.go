package controller

import (
	"context"
	"log"
	"managerService/internal/infrastructure/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	engine *gin.Engine
	srv    *http.Server
}

func NewController(cfg *config.Config) *Controller {
	engine := gin.Default()

	controller := &Controller{
		engine: engine,
		srv: &http.Server{
			Addr:    cfg.GetManagerAddr(),
			Handler: engine,
		},
	}

	controller.setAPIRouters()

	return controller
}

func (c *Controller) Run() {
	// 使用 goroutine 启动服务器
	go func() {
		if err := c.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动 Gin 服务失败: %v\n", err)
		}
	}()
}

func (c *Controller) Shutdown() {
	// 创建一个5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := c.srv.Shutdown(ctx); err != nil {
		log.Printf("关闭 Gin 服务失败: %v", err)
	}

	log.Println("成功关闭 Gin 服务")
}
