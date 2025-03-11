package security

import (
	"context"
	"errors"
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

// TokenBucketOption 定义令牌桶配置选项函数类型
type TokenBucketOption func(*TokenBucket)

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	// 基本配置
	bucketSize     int32         // 令牌桶总大小
	generationRate int64         // 每秒生成令牌数量 (使用int64存储float64位模式)
	bucket         chan struct{} // 令牌桶通道
	currentTokens  int32         // 当前桶中令牌数量

	// 高级特性
	burstable   int32 // 是否允许突发流量（启动时填满令牌桶）(0=false, 1=true)
	waitTimeout int64 // 获取令牌等待超时时间，0表示不等待 (纳秒)

	// 控制相关
	ticker         unsafe.Pointer // 定时器指针，用于生成令牌
	stopCh         chan struct{}  // 停止信号通道
	running        int32          // 运行状态标志
	lastRefillTime int64          // 上次填充令牌的时间 (Unix纳秒时间戳)

	// 指标收集
	totalRequests int64 // 总请求数
	rejectedCount int64 // 被拒绝的请求数
	waitingCount  int32 // 当前等待令牌的请求数
}

// 预定义错误
var (
	ErrNoToken       = errors.New("令牌桶为空，请求被限流")
	ErrWaitTimeout   = errors.New("等待令牌超时")
	ErrBucketStopped = errors.New("令牌桶已停止运行")
)

// float64与int64转换辅助函数
func float64ToInt64(f float64) int64 {
	return int64(math.Float64bits(f))
}

func int64ToFloat64(i int64) float64 {
	return math.Float64frombits(uint64(i))
}

// WithBucketSize 设置令牌桶容量
func WithBucketSize(size int32) TokenBucketOption {
	return func(tb *TokenBucket) {
		tb.bucketSize = size
		tb.bucket = make(chan struct{}, size)
	}
}

// WithGenerationRate 设置令牌生成速率（每秒）
func WithGenerationRate(rate float64) TokenBucketOption {
	return func(tb *TokenBucket) {
		atomic.StoreInt64(&tb.generationRate, float64ToInt64(rate))
	}
}

// WithBurstable 设置是否允许突发流量（启动时填满令牌桶）
func WithBurstable(enabled bool) TokenBucketOption {
	return func(tb *TokenBucket) {
		var val int32 = 0
		if enabled {
			val = 1
		}
		atomic.StoreInt32(&tb.burstable, val)
	}
}

// WithWaitTimeout 设置获取令牌的等待超时时间
func WithWaitTimeout(timeout time.Duration) TokenBucketOption {
	return func(tb *TokenBucket) {
		atomic.StoreInt64(&tb.waitTimeout, int64(timeout))
	}
}

// NewTokenBucket 创建新的令牌桶限流器
func NewTokenBucket(opts ...TokenBucketOption) *TokenBucket {
	tb := &TokenBucket{
		bucketSize:     1000,                // 默认桶大小1000
		generationRate: float64ToInt64(500), // 默认每秒生成500个令牌
		bucket:         make(chan struct{}, 1000),
		currentTokens:  0,
		burstable:      0, // 默认不允许突发
		waitTimeout:    0, // 默认不等待
		stopCh:         make(chan struct{}),
		running:        0,
		lastRefillTime: time.Now().UnixNano(),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(tb)
	}

	// 启动令牌生成器
	tb.Start()

	return tb
}

// Start 启动令牌桶
func (t *TokenBucket) Start() {
	// 确保只启动一次
	if !atomic.CompareAndSwapInt32(&t.running, 0, 1) {
		return
	}

	// 获取生成速率
	generationRate := int64ToFloat64(atomic.LoadInt64(&t.generationRate))

	// 计算填充间隔，确保填充频率适合配置的生成速率
	interval := time.Duration(float64(time.Second) / generationRate)
	if interval < time.Millisecond {
		interval = time.Millisecond // 最小间隔1毫秒
	}

	// 创建ticker并存储其指针
	ticker := time.NewTicker(interval)
	atomic.StorePointer(&t.ticker, unsafe.Pointer(&ticker))

	tokensPerInterval := generationRate * float64(interval) / float64(time.Second)

	// 初始化时如果允许突发，则填满令牌桶
	if atomic.LoadInt32(&t.burstable) != 0 {
		t.fillBucket()
	}

	// 启动令牌生成goroutine
	go func() {
		fractionalTokens := 0.0 // 处理小数部分令牌

		for {
			select {
			case <-ticker.C:
				// 计算此次应生成的令牌数（包含小数部分累积）
				tokensToAdd := tokensPerInterval + fractionalTokens
				newTokens := int(tokensToAdd)
				fractionalTokens = tokensToAdd - float64(newTokens) // 保存小数部分

				// 添加令牌
				t.addTokens(int32(newTokens))

				// 更新最后填充时间
				atomic.StoreInt64(&t.lastRefillTime, time.Now().UnixNano())

			case <-t.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止令牌桶
func (t *TokenBucket) Stop() {
	if atomic.CompareAndSwapInt32(&t.running, 1, 0) {
		// 获取ticker指针并停止
		tickerPtr := atomic.LoadPointer(&t.ticker)
		if tickerPtr != nil {
			ticker := (*time.Ticker)(tickerPtr)
			ticker.Stop()
		}
		close(t.stopCh)
	}
}

// GetToken 获取一个令牌，可设置超时
func (t *TokenBucket) GetToken(ctx context.Context) error {
	// 检查令牌桶是否运行
	if atomic.LoadInt32(&t.running) == 0 {
		atomic.AddInt64(&t.rejectedCount, 1)
		return ErrBucketStopped
	}

	atomic.AddInt64(&t.totalRequests, 1)
	atomic.AddInt32(&t.waitingCount, 1)
	defer atomic.AddInt32(&t.waitingCount, -1)

	// 根据waitTimeout决定是否等待
	waitTimeout := time.Duration(atomic.LoadInt64(&t.waitTimeout))
	if waitTimeout > 0 {
		return t.getTokenWithTimeout(ctx, waitTimeout)
	}

	// 非等待模式，立即尝试获取令牌
	select {
	case <-t.bucket:
		atomic.AddInt32(&t.currentTokens, -1)
		return nil
	case <-ctx.Done():
		atomic.AddInt64(&t.rejectedCount, 1)
		return ctx.Err()
	default:
		atomic.AddInt64(&t.rejectedCount, 1)
		return ErrNoToken
	}
}

// getTokenWithTimeout 带超时的令牌获取
func (t *TokenBucket) getTokenWithTimeout(ctx context.Context, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-t.bucket:
		atomic.AddInt32(&t.currentTokens, -1)
		return nil
	case <-timer.C:
		atomic.AddInt64(&t.rejectedCount, 1)
		return ErrWaitTimeout
	case <-ctx.Done():
		atomic.AddInt64(&t.rejectedCount, 1)
		return ctx.Err()
	}
}

// TryGetToken 尝试获取令牌，不阻塞，立即返回是否成功
func (t *TokenBucket) TryGetToken() bool {
	select {
	case <-t.bucket:
		atomic.AddInt32(&t.currentTokens, -1)
		atomic.AddInt64(&t.totalRequests, 1)
		return true
	default:
		atomic.AddInt64(&t.totalRequests, 1)
		atomic.AddInt64(&t.rejectedCount, 1)
		return false
	}
}

// GetTokens 获取多个令牌 - 优化版本
func (t *TokenBucket) GetTokens(ctx context.Context, count int32) error {
	if count <= 0 {
		return nil
	}

	if count > t.bucketSize {
		atomic.AddInt64(&t.rejectedCount, 1)
		return errors.New("请求的令牌数超过桶容量")
	}

	// 先检查是否有足够的令牌
	if atomic.LoadInt32(&t.currentTokens) < count {
		atomic.AddInt64(&t.rejectedCount, 1)
		return ErrNoToken
	}

	atomic.AddInt64(&t.totalRequests, 1)
	atomic.AddInt32(&t.waitingCount, 1)
	defer atomic.AddInt32(&t.waitingCount, -1)

	// 创建一个临时通道来接收令牌
	tokens := make([]struct{}, 0, count)

	// 尝试获取所需数量的令牌
	for i := int32(0); i < count; i++ {
		select {
		case token := <-t.bucket:
			tokens = append(tokens, token)
		case <-ctx.Done():
			// 上下文被取消，将已获取的令牌放回
			t.returnTokens(tokens)
			atomic.AddInt64(&t.rejectedCount, 1)
			return ctx.Err()
		default:
			// 令牌不足，将已获取的令牌放回
			t.returnTokens(tokens)
			atomic.AddInt64(&t.rejectedCount, 1)
			return ErrNoToken
		}
	}

	// 成功获取所有令牌
	atomic.AddInt32(&t.currentTokens, -count)
	return nil
}

// returnTokens 将令牌返回到桶中
func (t *TokenBucket) returnTokens(tokens []struct{}) {
	for range tokens {
		select {
		case t.bucket <- struct{}{}:
			// 令牌成功返回
		default:
			// 桶已满，丢弃多余的令牌
		}
	}
}

// addTokens 添加指定数量的令牌到桶中
func (t *TokenBucket) addTokens(count int32) {
	if count <= 0 {
		return
	}

	// 计算可以添加的令牌数量
	available := t.bucketSize - atomic.LoadInt32(&t.currentTokens)
	toAdd := int(math.Min(float64(count), float64(available)))

	// 向桶中添加令牌
	for i := 0; i < toAdd; i++ {
		select {
		case t.bucket <- struct{}{}:
			atomic.AddInt32(&t.currentTokens, 1)
		default:
			// 桶已满，停止添加
			return
		}
	}
}

// fillBucket 填满令牌桶
func (t *TokenBucket) fillBucket() {
	t.addTokens(t.bucketSize)
}

// UpdateRate 动态更新令牌生成速率
func (t *TokenBucket) UpdateRate(newRate float64) {
	if newRate <= 0 {
		return
	}

	// 原子方式更新生成速率
	atomic.StoreInt64(&t.generationRate, float64ToInt64(newRate))

	// 获取ticker指针
	tickerPtr := atomic.LoadPointer(&t.ticker)
	if tickerPtr != nil {
		// 停止旧ticker
		oldTicker := (*time.Ticker)(tickerPtr)
		oldTicker.Stop()

		// 创建新ticker
		interval := time.Duration(float64(time.Second) / newRate)
		if interval < time.Millisecond {
			interval = time.Millisecond
		}
		newTicker := time.NewTicker(interval)

		// 原子方式更新ticker指针
		atomic.StorePointer(&t.ticker, unsafe.Pointer(&newTicker))
	}
}

// UpdateBucketSize 动态更新桶容量 - 优化版本
func (t *TokenBucket) UpdateBucketSize(newSize int32) {
	if newSize <= 0 || newSize == t.bucketSize {
		return
	}

	// 创建新桶
	newBucket := make(chan struct{}, newSize)
	oldBucket := t.bucket

	// 获取当前令牌数
	currentTokens := atomic.LoadInt32(&t.currentTokens)

	// 计算新桶中应填充的令牌数
	tokensToAdd := currentTokens
	if tokensToAdd > newSize {
		tokensToAdd = newSize
	}

	// 填充新桶
	for i := int32(0); i < tokensToAdd; i++ {
		newBucket <- struct{}{}
	}

	// 原子方式更新桶引用和容量
	t.bucket = newBucket
	t.bucketSize = newSize
	atomic.StoreInt32(&t.currentTokens, tokensToAdd)

	// 清理旧桶 (异步进行，避免阻塞)
	go func(oldChan chan struct{}) {
		// 尽可能清空旧桶，避免资源泄漏
		for {
			select {
			case <-oldChan:
				// 取出旧桶中的令牌
				continue
			default:
				return
			}
		}
	}(oldBucket)
}

// GetMetrics 获取限流指标
func (t *TokenBucket) GetMetrics() map[string]interface{} {
	// 获取当前配置和状态
	generationRate := int64ToFloat64(atomic.LoadInt64(&t.generationRate))
	waitTimeout := time.Duration(atomic.LoadInt64(&t.waitTimeout))
	lastRefillTimeNano := atomic.LoadInt64(&t.lastRefillTime)

	return map[string]interface{}{
		"total_requests":    atomic.LoadInt64(&t.totalRequests),
		"rejected_requests": atomic.LoadInt64(&t.rejectedCount),
		"waiting_requests":  atomic.LoadInt32(&t.waitingCount),
		"current_tokens":    atomic.LoadInt32(&t.currentTokens),
		"bucket_size":       t.bucketSize,
		"generation_rate":   generationRate,
		"wait_timeout":      waitTimeout.String(),
		"last_refill_time":  time.Unix(0, lastRefillTimeNano),
		"running":           atomic.LoadInt32(&t.running) == 1,
		"burstable":         atomic.LoadInt32(&t.burstable) != 0,
	}
}

// Reset 重置令牌桶 - 优化版本
func (t *TokenBucket) Reset() {
	// 重置计数器
	atomic.StoreInt64(&t.totalRequests, 0)
	atomic.StoreInt64(&t.rejectedCount, 0)
	atomic.StoreInt32(&t.waitingCount, 0)

	// 创建新的令牌桶通道代替清空操作
	oldBucket := t.bucket
	t.bucket = make(chan struct{}, t.bucketSize)
	atomic.StoreInt32(&t.currentTokens, 0)

	// 异步清理旧桶
	go func(oldChan chan struct{}) {
		for {
			select {
			case <-oldChan:
				continue
			default:
				return
			}
		}
	}(oldBucket)

	// 如果允许突发流量，则重新填满
	if atomic.LoadInt32(&t.burstable) != 0 {
		t.fillBucket()
	}
}
