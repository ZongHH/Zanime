package bootstrap

import (
	"context"
	"crawler/pkg/monitor"
	"log"
	"net"
)

type Bootstrap struct {
	container *Container
	cancel    context.CancelFunc
}

func NewBootstrap(configPath string) *Bootstrap {
	container := BuildContainer(configPath)

	return &Bootstrap{
		container: container,
	}
}

func (b *Bootstrap) Start() {
	monitor.Init(monitor.NewLogConfig())

	ctx, cancel := context.WithCancel(context.Background())
	go b.container.Services.Crawler.Start(ctx)
	b.cancel = cancel

	var err error
	b.container.Interfaces.Lis, err = net.Listen("tcp", b.container.Cfg.GetGRPCAddress())
	if err != nil {
		log.Fatalf("监听GRPC地址失败: %v", err)
	}
	go func() {
		if err := b.container.Interfaces.GRPC.Serve(b.container.Interfaces.Lis); err != nil {
			log.Fatalf("GRPC服务启动失败: %v", err)
		}
	}()
}

func (b *Bootstrap) Stop() {
	b.container.Interfaces.GRPC.GracefulStop()
	b.container.Interfaces.Lis.Close()
	b.cancel()
	b.container.Bases.close()
	monitor.Close()
}
