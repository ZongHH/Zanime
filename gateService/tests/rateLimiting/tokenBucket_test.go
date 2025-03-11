package ratelimiting

import (
	"context"
	"gateService/internal/interfaces/http/middleware/security"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTokenBucketBasicOperations 测试令牌桶的基本操作
func TestTokenBucketBasicOperations(t *testing.T) {
	// 创建一个令牌桶，每秒生成10个令牌，最大容量为20
	limiter := security.NewTokenBucket(
		security.WithBucketSize(20),
		security.WithGenerationRate(10.0),
		security.WithBurstable(true),
	)
	defer limiter.Stop()

	t.Run("初始状态检查", func(t *testing.T) {
		metrics := limiter.GetMetrics()
		assert.Equal(t, int32(20), metrics["bucket_size"])
		assert.Equal(t, float64(10), metrics["generation_rate"])
		assert.True(t, metrics["running"].(bool))
	})

	t.Run("成功获取令牌", func(t *testing.T) {
		ctx := context.Background()
		err := limiter.GetToken(ctx)
		assert.NoError(t, err)
	})

	t.Run("尝试获取令牌", func(t *testing.T) {
		success := limiter.TryGetToken()
		assert.True(t, success)
	})
}

// TestTokenBucketRateLimiting 测试令牌桶的限流效果
func TestTokenBucketRateLimiting(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(5),
		security.WithGenerationRate(2.0), // 每秒生成2个令牌
		security.WithBurstable(true),
	)
	defer limiter.Stop()

	t.Run("限流测试", func(t *testing.T) {
		ctx := context.Background()

		// 快速消耗所有令牌
		for i := 0; i < 5; i++ {
			err := limiter.GetToken(ctx)
			assert.NoError(t, err)
		}

		// 此时应该被限流
		err := limiter.GetToken(ctx)
		assert.Equal(t, security.ErrNoToken, err)

		// 等待令牌生成
		time.Sleep(time.Second)

		// 应该可以获取到新生成的令牌
		err = limiter.GetToken(ctx)
		assert.NoError(t, err)
	})
}

// TestTokenBucketConcurrency 测试令牌桶的并发性能
func TestTokenBucketConcurrency(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(100),
		security.WithGenerationRate(50.0),
	)
	defer limiter.Stop()

	t.Run("并发获取令牌", func(t *testing.T) {
		var wg sync.WaitGroup
		var successCount int32
		workerCount := 100

		// 启动多个goroutine并发获取令牌
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if limiter.TryGetToken() {
					atomic.AddInt32(&successCount, 1)
				}
			}()
		}

		wg.Wait()

		// 验证指标
		metrics := limiter.GetMetrics()
		assert.Equal(t, int64(workerCount), metrics["total_requests"])
		assert.True(t, metrics["rejected_requests"].(int64) > 0)
	})
}

// TestTokenBucketTimeout 测试令牌桶的超时机制
func TestTokenBucketTimeout(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(5),
		security.WithGenerationRate(1.0),
		security.WithWaitTimeout(time.Millisecond*100),
	)
	defer limiter.Stop()

	t.Run("等待超时", func(t *testing.T) {
		ctx := context.Background()

		// 消耗所有令牌
		for i := 0; i < 5; i++ {
			_ = limiter.GetToken(ctx)
		}

		// 应该发生超时
		err := limiter.GetToken(ctx)
		assert.Equal(t, security.ErrWaitTimeout, err)
	})
}

// TestTokenBucketBurst 测试突发流量处理
func TestTokenBucketBurst(t *testing.T) {
	t.Run("允许突发", func(t *testing.T) {
		limiter := security.NewTokenBucket(
			security.WithBucketSize(10),
			security.WithGenerationRate(1.0),
			security.WithBurstable(true),
		)
		defer limiter.Stop()

		ctx := context.Background()

		// 由于启用了突发模式，应该可以立即获取多个令牌
		for i := 0; i < 10; i++ {
			err := limiter.GetToken(ctx)
			assert.NoError(t, err)
		}
	})

	t.Run("禁止突发", func(t *testing.T) {
		limiter := security.NewTokenBucket(
			security.WithBucketSize(10),
			security.WithGenerationRate(1.0),
			security.WithBurstable(false),
		)
		defer limiter.Stop()

		// 由于禁用了突发模式，应该无法立即获取多个令牌
		success := 0
		for i := 0; i < 10; i++ {
			if limiter.TryGetToken() {
				success++
			}
		}
		assert.True(t, success < 10)
	})
}

// TestTokenBucketDynamicUpdate 测试动态更新配置
func TestTokenBucketDynamicUpdate(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(10),
		security.WithGenerationRate(1.0),
	)
	defer limiter.Stop()

	t.Run("更新生成速率", func(t *testing.T) {
		oldRate := limiter.GetMetrics()["generation_rate"].(float64)
		newRate := 5.0

		limiter.UpdateRate(newRate)
		time.Sleep(time.Millisecond * 100) // 等待更新生效

		currentRate := limiter.GetMetrics()["generation_rate"].(float64)
		assert.NotEqual(t, oldRate, currentRate)
		assert.Equal(t, newRate, currentRate)
	})

	t.Run("更新桶大小", func(t *testing.T) {
		oldSize := limiter.GetMetrics()["bucket_size"].(int32)
		newSize := int32(20)

		limiter.UpdateBucketSize(newSize)
		time.Sleep(time.Millisecond * 100) // 等待更新生效

		currentSize := limiter.GetMetrics()["bucket_size"].(int32)
		assert.NotEqual(t, oldSize, currentSize)
		assert.Equal(t, newSize, currentSize)
	})
}

// TestTokenBucketReset 测试重置功能
func TestTokenBucketReset(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(10),
		security.WithGenerationRate(1.0),
	)
	defer limiter.Stop()

	ctx := context.Background()

	// 先消耗一些令牌
	for i := 0; i < 5; i++ {
		_ = limiter.GetToken(ctx)
	}

	// 获取重置前的指标
	beforeMetrics := limiter.GetMetrics()

	// 重置令牌桶
	limiter.Reset()

	// 获取重置后的指标
	afterMetrics := limiter.GetMetrics()

	// 验证重置效果
	assert.NotEqual(t, beforeMetrics["current_tokens"], afterMetrics["current_tokens"])
	assert.Equal(t, int64(0), afterMetrics["total_requests"])
	assert.Equal(t, int64(0), afterMetrics["rejected_requests"])
}

// TestTokenBucketContextCancellation 测试上下文取消
func TestTokenBucketContextCancellation(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(5),
		security.WithGenerationRate(1.0),
	)
	defer limiter.Stop()

	// 消耗所有令牌
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_ = limiter.GetToken(ctx)
	}

	// 创建一个带取消的上下文
	ctxWithCancel, cancel := context.WithCancel(context.Background())

	// 在另一个goroutine中取消上下文
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	// 尝试获取令牌，应该因为上下文取消而失败
	err := limiter.GetToken(ctxWithCancel)
	assert.Equal(t, context.Canceled, err)
}

// TestTokenBucketEdgeCases 测试边界情况
func TestTokenBucketEdgeCases(t *testing.T) {
	t.Run("零容量桶", func(t *testing.T) {
		limiter := security.NewTokenBucket(
			security.WithBucketSize(0),
		)
		defer limiter.Stop()

		assert.Equal(t, int32(0), limiter.GetMetrics()["bucket_size"])
	})

	t.Run("负生成速率", func(t *testing.T) {
		limiter := security.NewTokenBucket()
		defer limiter.Stop()

		limiter.UpdateRate(-1.0)
		rate := limiter.GetMetrics()["generation_rate"].(float64)
		assert.True(t, rate > 0)
	})

	t.Run("停止后获取令牌", func(t *testing.T) {
		limiter := security.NewTokenBucket()
		limiter.Stop()

		err := limiter.GetToken(context.Background())
		assert.Equal(t, security.ErrBucketStopped, err)
	})
}

// TestTokenBucketMultipleTokens 测试批量获取令牌
func TestTokenBucketMultipleTokens(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(10),
		security.WithGenerationRate(5.0),
		security.WithBurstable(true),
	)
	defer limiter.Stop()

	ctx := context.Background()

	t.Run("批量获取成功", func(t *testing.T) {
		err := limiter.GetTokens(ctx, 3)
		assert.NoError(t, err)
	})

	t.Run("请求过多令牌", func(t *testing.T) {
		err := limiter.GetTokens(ctx, 20)
		assert.Error(t, err)
	})

	t.Run("请求零个令牌", func(t *testing.T) {
		err := limiter.GetTokens(ctx, 0)
		assert.NoError(t, err)
	})
}

// TestTokenBucketMetrics 测试指标收集
func TestTokenBucketMetrics(t *testing.T) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(10),
		security.WithGenerationRate(5.0),
	)
	defer limiter.Stop()

	ctx := context.Background()

	// 生成一些指标数据
	for i := 0; i < 5; i++ {
		_ = limiter.GetToken(ctx)
	}
	_ = limiter.GetToken(ctx) // 应该被拒绝

	metrics := limiter.GetMetrics()

	// 验证指标完整性
	assert.NotNil(t, metrics["total_requests"])
	assert.NotNil(t, metrics["rejected_requests"])
	assert.NotNil(t, metrics["waiting_requests"])
	assert.NotNil(t, metrics["current_tokens"])
	assert.NotNil(t, metrics["bucket_size"])
	assert.NotNil(t, metrics["generation_rate"])
	assert.NotNil(t, metrics["wait_timeout"])
	assert.NotNil(t, metrics["last_refill_time"])
	assert.NotNil(t, metrics["running"])
	assert.NotNil(t, metrics["burstable"])

	// 验证指标值的合理性
	assert.Equal(t, int64(6), metrics["total_requests"])
	assert.True(t, metrics["rejected_requests"].(int64) > 0)
	assert.Equal(t, int32(0), metrics["waiting_requests"])
}

// BenchmarkTokenBucket 基准测试
func BenchmarkTokenBucket(b *testing.B) {
	limiter := security.NewTokenBucket(
		security.WithBucketSize(1000),
		security.WithGenerationRate(1000),
	)
	defer limiter.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.TryGetToken()
		}
	})
}
