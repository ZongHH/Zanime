package nsqpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nsqio/go-nsq"
)

// ConsumerOptions 消费者配置选项
type ConsumerOptions struct {
	// 基础配置
	Topic       string
	Channel     string
	PoolSize    int
	NSQDAddress string

	// 消费者配置
	MaxInFlight       int           // 最大处理中的消息数
	MaxAttempts       int           // 最大重试次数
	RequeueDelay      time.Duration // 重试延迟时间
	MaxProcessTimeout time.Duration // 单条消息处理超时时间

	// 监控配置
	EnableMetrics bool // 是否启用监控
}

// ConsumerPool NSQ消费者连接池
type ConsumerPool struct {
	opts         *ConsumerOptions
	config       *nsq.Config
	consumers    []*nsq.Consumer
	callBackFunc func(context.Context, []byte) error

	// 监控指标
	messageCount   int64
	errorCount     int64
	processingTime int64
	activeWorkers  int32

	// 状态控制
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.RWMutex
	started bool
}

// NewConsumerPool 创建消费者池
func NewConsumerPool(opts *ConsumerOptions) (*ConsumerPool, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建NSQ配置
	config := nsq.NewConfig()
	config.MaxInFlight = opts.MaxInFlight
	config.MaxAttempts = uint16(opts.MaxAttempts)
	config.DefaultRequeueDelay = opts.RequeueDelay
	config.ReadTimeout = opts.MaxProcessTimeout
	config.WriteTimeout = opts.MaxProcessTimeout
	config.MsgTimeout = opts.MaxProcessTimeout

	pool := &ConsumerPool{
		opts:      opts,
		config:    config,
		consumers: make([]*nsq.Consumer, 0, opts.PoolSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	// 初始化消费者实例
	if err := pool.initConsumers(); err != nil {
		return nil, err
	}

	return pool, nil
}

// initConsumers 初始化消费者实例
func (p *ConsumerPool) initConsumers() error {
	for i := 0; i < p.opts.PoolSize; i++ {
		consumer, err := nsq.NewConsumer(p.opts.Topic, p.opts.Channel, p.config)
		if err != nil {
			return fmt.Errorf("创建NSQ消费者实例失败: %w", err)
		}
		consumer.AddHandler(p)
		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

// HandleMessage 消息处理
func (p *ConsumerPool) HandleMessage(msg *nsq.Message) error {
	atomic.AddInt32(&p.activeWorkers, 1)
	defer atomic.AddInt32(&p.activeWorkers, -1)

	// 创建消息处理上下文
	ctx, cancel := context.WithTimeout(p.ctx, p.opts.MaxProcessTimeout)
	defer cancel()

	startTime := time.Now()

	// 执行用户回调
	err := p.processMessage(ctx, msg)

	// 更新监控指标
	p.updateMetrics(startTime, err)

	if err != nil {
		// 处理失败，进行重试
		if msg.Attempts < uint16(p.opts.MaxAttempts) {
			msg.RequeueWithoutBackoff(p.opts.RequeueDelay)
			return nil
		}
		// 超过最大重试次数
		return err
	}

	msg.Finish()
	return nil
}

// processMessage 处理单条消息
func (p *ConsumerPool) processMessage(ctx context.Context, msg *nsq.Message) error {
	if p.callBackFunc == nil {
		return fmt.Errorf("callback function not registered")
	}

	return p.callBackFunc(ctx, msg.Body)
}

// updateMetrics 更新监控指标
func (p *ConsumerPool) updateMetrics(startTime time.Time, err error) {
	if !p.opts.EnableMetrics {
		return
	}

	atomic.AddInt64(&p.messageCount, 1)
	if err != nil {
		atomic.AddInt64(&p.errorCount, 1)
	}
	atomic.AddInt64(&p.processingTime, time.Since(startTime).Milliseconds())
}

// RegisterCallback 注册消息处理回调
func (p *ConsumerPool) RegisterCallback(fn func(context.Context, []byte) error) {
	p.callBackFunc = fn
}

// Start 启动消费者池
func (p *ConsumerPool) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return fmt.Errorf("consumer pool already started")
	}

	for _, consumer := range p.consumers {
		if err := consumer.ConnectToNSQD(p.opts.NSQDAddress); err != nil {
			return fmt.Errorf("connect to NSQD failed: %w", err)
		}
	}

	p.started = true
	return nil
}

// Stop 停止消费者池
func (p *ConsumerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.started {
		return
	}

	// 取消上下文
	p.cancel()

	// 停止所有消费者
	for _, consumer := range p.consumers {
		consumer.Stop()
	}

	p.started = false
}

// GetMetrics 获取监控指标
func (p *ConsumerPool) GetMetrics() map[string]interface{} {
	if !p.opts.EnableMetrics {
		return nil
	}

	return map[string]interface{}{
		"message_count":    atomic.LoadInt64(&p.messageCount),
		"error_count":      atomic.LoadInt64(&p.errorCount),
		"processing_time":  atomic.LoadInt64(&p.processingTime),
		"active_workers":   atomic.LoadInt32(&p.activeWorkers),
		"error_rate":       float64(p.errorCount) / float64(p.messageCount),
		"avg_process_time": float64(p.processingTime) / float64(p.messageCount),
	}
}

// validateOptions 验证配置选项
func validateOptions(opts *ConsumerOptions) error {
	if opts.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if opts.Channel == "" {
		return fmt.Errorf("channel is required")
	}
	if opts.PoolSize <= 0 {
		opts.PoolSize = 1
	}
	if opts.MaxInFlight <= 0 {
		opts.MaxInFlight = 300
	}
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = 5
	}
	if opts.RequeueDelay <= 0 {
		opts.RequeueDelay = 5 * time.Second
	}
	if opts.MaxProcessTimeout <= 0 {
		opts.MaxProcessTimeout = 60 * time.Second
	}
	if opts.NSQDAddress == "" {
		opts.NSQDAddress = "127.0.0.1:5150"
	}
	return nil
}
