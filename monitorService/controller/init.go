package controller

import (
	"context"
	"log"
	"monitorService/pkg/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	engine *gin.Engine
	srv    *http.Server
}

func NewController() *Controller {
	engine := gin.Default()
	ginaddr, err := config.GetHostAndPort("monitorService")
	if err != nil {
		log.Fatal("config.GetHostAndPort() ", err)
	}
	return &Controller{
		engine: engine,
		srv: &http.Server{
			Addr:    ginaddr,
			Handler: engine,
		},
	}
}

func (c *Controller) Run() {
	// 使用 goroutine 启动服务器
	go func() {
		if err := c.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (c *Controller) Shutdown() {
	// 创建一个5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := c.srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}

var globalController *Controller

func Init() {
	globalController = NewController()
	globalController.setAPIRouters()
	globalController.Run()
}

func Close() {
	if globalController != nil {
		globalController.Shutdown()
	}
}
