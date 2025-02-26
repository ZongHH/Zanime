package nsqpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nsqio/go-nsq"
)

// ProducerOptions 生产者配置选项
type ProducerOptions struct {
	// 基础配置
	PoolSize    int
	NSQDAddress string

	// 生产者配置
	WriteTimeout      time.Duration // 写超时时间
	ReadTimeout       time.Duration // 读超时时间
	DialTimeout       time.Duration // 连接超时时间
	MaxRetryCount     int           // 最大重试次数
	RetryInterval     time.Duration // 重试间隔
	HeartbeatInterval time.Duration // 心跳间隔

	// 监控配置
	EnableMetrics bool // 是否启用监控
}

// ProducerPool NSQ生产者连接池
type ProducerPool struct {
	opts      *ProducerOptions
	config    *nsq.Config
	producers chan *nsq.Producer // 使用channel替代slice，更安全的并发访问
	doneChan  chan *nsq.ProducerTransaction

	// 监控指标
	messageCount    int64 // 发送消息总数
	errorCount      int64 // 错误总数
	processingTime  int64 // 处理时间
	activeProducers int32 // 活跃生产者数量
	availableCount  int32 // 可用连接数

	// 状态控制
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.RWMutex
	started bool
}

// NewProducerPool 创建生产者池
func NewProducerPool(opts *ProducerOptions) (*ProducerPool, error) {
	if err := validateProducerOptions(opts); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建NSQ配置
	config := nsq.NewConfig()
	config.WriteTimeout = opts.WriteTimeout
	config.ReadTimeout = opts.ReadTimeout
	config.DialTimeout = opts.DialTimeout
	config.HeartbeatInterval = opts.HeartbeatInterval

	pool := &ProducerPool{
		opts:      opts,
		config:    config,
		producers: make(chan *nsq.Producer, opts.PoolSize),
		doneChan:  make(chan *nsq.ProducerTransaction, 1000),
		ctx:       ctx,
		cancel:    cancel,
	}

	// 初始化生产者
	if err := pool.initProducers(); err != nil {
		return nil, err
	}

	// 启动异步结果处理
	go pool.handleAsyncResults()

	// 启动健康检查
	if opts.HeartbeatInterval > 0 {
		go pool.healthCheck()
	}

	return pool, nil
}

// initProducers 初始化生产者实例
func (p *ProducerPool) initProducers() error {
	for i := 0; i < p.opts.PoolSize; i++ {
		producer, err := nsq.NewProducer(p.opts.NSQDAddress, p.config)
		if err != nil {
			return fmt.Errorf("创建NSQ生产者实例失败: %w", err)
		}
		p.producers <- producer
		atomic.AddInt32(&p.availableCount, 1)
	}
	return nil
}

// getProducer 获取生产者实例
func (p *ProducerPool) getProducer(ctx context.Context) (*nsq.Producer, error) {
	atomic.AddInt32(&p.activeProducers, 1)

	select {
	case producer := <-p.producers:
		atomic.AddInt32(&p.availableCount, -1)
		return producer, nil
	case <-ctx.Done():
		atomic.AddInt32(&p.activeProducers, -1)
		return nil, ctx.Err()
	}
}

// releaseProducer 释放生产者实例
func (p *ProducerPool) releaseProducer(producer *nsq.Producer) {
	p.producers <- producer
	atomic.AddInt32(&p.availableCount, 1)
	atomic.AddInt32(&p.activeProducers, -1)
}

// Publish 同步发布消息
func (p *ProducerPool) Publish(ctx context.Context, topic string, msg []byte) error {
	startTime := time.Now()

	for i := 0; i <= p.opts.MaxRetryCount; i++ {
		producer, err := p.getProducer(ctx)
		if err != nil {
			continue
		}

		err = producer.Publish(topic, msg)
		p.releaseProducer(producer)

		if err == nil {
			p.updateMetrics(startTime, nil)
			return nil
		}

		// 重试等待
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(p.opts.RetryInterval):
			continue
		}
	}

	err := fmt.Errorf("发布消息失败，已重试%d次", p.opts.MaxRetryCount)
	p.updateMetrics(startTime, err)
	return err
}

// PublishAsync 异步发布消息
func (p *ProducerPool) PublishAsync(ctx context.Context, topic string, msg []byte, args ...interface{}) error {
	producer, err := p.getProducer(ctx)
	if err != nil {
		return err
	}

	err = producer.PublishAsync(topic, msg, p.doneChan, args...)
	if err != nil {
		p.releaseProducer(producer)
		return err
	}

	p.releaseProducer(producer)
	return nil
}

// DeferredPublish 延迟发布消息
func (p *ProducerPool) DeferredPublish(ctx context.Context, topic string, delay time.Duration, msg []byte) error {
	startTime := time.Now()

	producer, err := p.getProducer(ctx)
	if err != nil {
		return err
	}
	defer p.releaseProducer(producer)

	err = producer.DeferredPublish(topic, delay, msg)
	p.updateMetrics(startTime, err)
	return err
}

// handleAsyncResults 处理异步发布结果
func (p *ProducerPool) handleAsyncResults() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case trans := <-p.doneChan:
			atomic.AddInt64(&p.messageCount, 1)
			if trans.Error != nil {
				atomic.AddInt64(&p.errorCount, 1)
			}
		}
	}
}

// healthCheck 健康检查
func (p *ProducerPool) healthCheck() {
	ticker := time.NewTicker(p.opts.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.checkAndRepairConnections()
		}
	}
}

// checkAndRepairConnections 检查并修复连接
func (p *ProducerPool) checkAndRepairConnections() {
	available := atomic.LoadInt32(&p.availableCount)
	if available < int32(p.opts.PoolSize) {
		p.mu.Lock()
		defer p.mu.Unlock()

		// 补充连接
		for i := available; i < int32(p.opts.PoolSize); i++ {
			producer, err := nsq.NewProducer(p.opts.NSQDAddress, p.config)
			if err != nil {
				continue
			}
			p.producers <- producer
			atomic.AddInt32(&p.availableCount, 1)
		}
	}
}

// GetMetrics 获取监控指标
func (p *ProducerPool) GetMetrics() map[string]interface{} {
	if !p.opts.EnableMetrics {
		return nil
	}

	messageCount := atomic.LoadInt64(&p.messageCount)
	return map[string]interface{}{
		"message_count":    messageCount,
		"error_count":      atomic.LoadInt64(&p.errorCount),
		"processing_time":  atomic.LoadInt64(&p.processingTime),
		"active_producers": atomic.LoadInt32(&p.activeProducers),
		"available_count":  atomic.LoadInt32(&p.availableCount),
		"error_rate":       float64(p.errorCount) / float64(messageCount),
	}
}

// Close 关闭生产者池
func (p *ProducerPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return
	}

	p.cancel()
	close(p.doneChan)

	// 关闭所有生产者
	for i := 0; i < p.opts.PoolSize; i++ {
		if producer, err := p.getProducer(context.Background()); err == nil {
			producer.Stop()
		}
	}

	p.started = false
}

// validateProducerOptions 验证配置选项
func validateProducerOptions(opts *ProducerOptions) error {
	if opts.NSQDAddress == "" {
		opts.NSQDAddress = "127.0.0.1:5150"
	}
	if opts.PoolSize <= 0 {
		opts.PoolSize = 3
	}
	if opts.WriteTimeout <= 0 {
		opts.WriteTimeout = 50 * time.Second
	}
	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = 50 * time.Second
	}
	if opts.DialTimeout <= 0 {
		opts.DialTimeout = 50 * time.Second
	}
	if opts.MaxRetryCount <= 0 {
		opts.MaxRetryCount = 3
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = 100 * time.Millisecond
	}
	if opts.HeartbeatInterval <= 0 {
		opts.HeartbeatInterval = 30 * time.Second
	}
	return nil
}

func (p *ProducerPool) updateMetrics(startTime time.Time, err error) {
	atomic.AddInt64(&p.processingTime, int64(time.Since(startTime)))
	if err != nil {
		atomic.AddInt64(&p.errorCount, 1)
	}
	atomic.AddInt64(&p.messageCount, 1)
}
