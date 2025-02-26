package scrapeClient

import (
	"context"
	"fmt"
	"gateService/internal/infrastructure/config"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// GRPCClientPool gRPC客户端连接池
type GRPCClientPool struct {
	cfg      *config.GrpcServiceConfig
	conns    chan *grpcConn
	endpoint string

	// 监控指标
	activeConns  int32
	errorCount   int64
	requestCount int64
	responseTime int64
}

// grpcConn 封装的gRPC连接
type grpcConn struct {
	conn      *grpc.ClientConn
	client    VideoClient
	createAt  time.Time    // 连接创建时间
	lastCheck atomic.Value // 上次检查时间
}

// NewGRPCClientPool 创建一个新的gRPC客户端连接池
func NewGRPCClientPool(cfg *config.Config) (*GRPCClientPool, error) {
	if !cfg.TargetGrpcServers["scrape_service"].Enabled {
		return nil, fmt.Errorf("scrape service is disabled")
	}

	pool := &GRPCClientPool{
		cfg:      cfg.TargetGrpcServers["scrape_service"],
		conns:    make(chan *grpcConn, cfg.TargetGrpcServers["scrape_service"].Pool.MaxConns),
		endpoint: cfg.GetScrapeAddr(),
	}

	// 初始化连接池
	for i := 0; i < cfg.TargetGrpcServers["scrape_service"].Pool.MinIdleConns; i++ {
		conn, err := pool.createConn()
		if err != nil {
			return nil, fmt.Errorf("初始化连接池失败: %w", err)
		}
		pool.conns <- conn
	}

	// 启动健康检查，定期检查连接健康状态并补充连接
	go pool.healthCheck()
	// 启动连接有效性检查，定期检查连接有效性
	go pool.validityCheck()

	return pool, nil
}

// createConn 创建新的gRPC连接
func (p *GRPCClientPool) createConn() (*grpcConn, error) {
	// gRPC连接选项配置
	opts := []grpc.DialOption{
		// 使用不安全的传输凭证(禁用TLS)
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置客户端保活参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                p.cfg.Connection.KeepAlive.Time,                // 空闲时ping服务器的时间间隔
			Timeout:             p.cfg.Connection.KeepAlive.Timeout,             // ping后等待响应的超时时间
			PermitWithoutStream: p.cfg.Connection.KeepAlive.PermitWithoutStream, // 是否允许在无活动流时发送ping
		}),
	}

	conn, err := grpc.NewClient(p.endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("创建gRPC连接失败: %w", err)
	}

	gConn := &grpcConn{
		conn:     conn,
		client:   NewVideoClient(conn),
		createAt: time.Now(),
	}
	gConn.lastCheck.Store(time.Now())
	return gConn, nil
}

// checkConn 检查连接是否有效
func (p *GRPCClientPool) checkConn(conn *grpcConn) bool {
	// 检查连接是否过期
	if p.cfg.Pool.MaxConnAge > 0 && time.Since(conn.createAt) > p.cfg.Pool.MaxConnAge {
		return false
	}

	// 使用连接状态检查
	state := conn.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		return false
	}

	// 更新最后检查时间
	conn.lastCheck.Store(time.Now())
	return true
}

// getConn 从连接池获取连接
func (p *GRPCClientPool) getConn(ctx context.Context) (*grpcConn, error) {
	atomic.AddInt32(&p.activeConns, 1)

	for {
		select {
		case conn := <-p.conns:
			// 检查连接是否有效
			if p.checkConn(conn) {
				return conn, nil
			}
			// 连接无效，关闭并尝试获取新连接
			conn.conn.Close()
			continue

		case <-ctx.Done():
			atomic.AddInt32(&p.activeConns, -1)
			return nil, ctx.Err()

		default:
			// 如果没有可用连接且未达到最大连接数，创建新连接
			if len(p.conns) < p.cfg.Pool.MaxConns {
				conn, err := p.createConn()
				if err != nil {
					atomic.AddInt32(&p.activeConns, -1)
					return nil, err
				}
				return conn, nil
			}

			// 等待可用连接
			select {
			case conn := <-p.conns:
				if p.checkConn(conn) {
					return conn, nil
				}
				conn.conn.Close()
				continue
			case <-ctx.Done():
				atomic.AddInt32(&p.activeConns, -1)
				return nil, ctx.Err()
			}
		}
	}
}

// releaseConn 释放连接回连接池
func (p *GRPCClientPool) releaseConn(conn *grpcConn) {
	p.conns <- conn
	atomic.AddInt32(&p.activeConns, -1)
}

// ScrapeVideoUrl 获取视频URL
func (p *GRPCClientPool) ScrapeVideoUrl(ctx context.Context, name, release, area, episode string) (*VideoMsg, error) {
	startTime := time.Now()
	atomic.AddInt64(&p.requestCount, 1)

	conn, err := p.getConn(ctx)
	if err != nil {
		atomic.AddInt64(&p.errorCount, 1)
		return nil, err
	}
	defer p.releaseConn(conn)

	request := &VideoParms{
		Name:    name,
		Release: release,
		Area:    area,
		Episode: episode,
	}

	// 设置超时时间
	ctx, cancel := context.WithTimeout(ctx, p.cfg.Connection.Timeout)
	defer cancel()

	// 执行RPC调用
	response, err := conn.client.ScrapeVideoUrl(ctx, request)
	if err != nil {
		atomic.AddInt64(&p.errorCount, 1)
		return nil, fmt.Errorf("RPC调用失败: %w", err)
	}

	// 更新监控指标
	atomic.AddInt64(&p.responseTime, time.Since(startTime).Milliseconds())

	return response, nil
}

// healthCheck 定期检查连接健康状态并补充连接
func (p *GRPCClientPool) healthCheck() {
	ticker := time.NewTicker(p.cfg.Connection.KeepAlive.Time)
	defer ticker.Stop()

	for range ticker.C {
		currentConns := len(p.conns)

		// 如果连接数小于最小空闲连接数，创建新连接
		if currentConns < p.cfg.Pool.MinIdleConns {
			for i := currentConns; i < p.cfg.Pool.MinIdleConns; i++ {
				conn, err := p.createConn()
				if err != nil {
					continue
				}
				p.conns <- conn
			}
		}
	}
}

// validityCheck 定期检查连接有效性
func (p *GRPCClientPool) validityCheck() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for range ticker.C {
		currentConns := len(p.conns)
		validConns := make([]*grpcConn, 0, currentConns)

		// 检查所有连接
		for i := 0; i < currentConns; i++ {
			select {
			case conn := <-p.conns:
				if p.checkConn(conn) {
					validConns = append(validConns, conn)
				} else {
					conn.conn.Close()
				}
			default:
				goto RESTORE
			}
		}

	RESTORE:
		// 恢复有效连接到池中
		for _, conn := range validConns {
			p.conns <- conn
		}

		// 补充连接到最小空闲连接数
		currentConns = len(p.conns)
		if currentConns < p.cfg.Pool.MinIdleConns {
			for i := currentConns; i < p.cfg.Pool.MinIdleConns; i++ {
				if conn, err := p.createConn(); err == nil {
					p.conns <- conn
				}
			}
		}
	}
}

// GetMetrics 获取监控指标
func (p *GRPCClientPool) GetMetrics() map[string]interface{} {
	requestCount := atomic.LoadInt64(&p.requestCount)
	return map[string]interface{}{
		"active_connections": atomic.LoadInt32(&p.activeConns),
		"request_count":      requestCount,
		"error_count":        atomic.LoadInt64(&p.errorCount),
		"error_rate":         float64(p.errorCount) / float64(requestCount),
		"avg_response_time":  float64(atomic.LoadInt64(&p.responseTime)) / float64(requestCount),
	}
}

// Close 关闭连接池
func (p *GRPCClientPool) Close() error {
	// 关闭所有连接
	for i := 0; i < len(p.conns); i++ {
		select {
		case conn := <-p.conns:
			conn.conn.Close()
		default:
			return nil
		}
	}

	close(p.conns)
	return nil
}
