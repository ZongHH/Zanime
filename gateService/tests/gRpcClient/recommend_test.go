package gRpcClient

import (
	"context"
	"gateService/internal/grpc/client/recommend"
	"gateService/internal/infrastructure/config"
	"sync"
	"testing"
	"time"
)

func TestRecommendGRPCClientPool(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	t.Run("创建连接池", func(t *testing.T) {
		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}
		defer pool.Close()

		metrics := pool.GetMetrics()
		if metrics["active_connections"].(int32) != 0 {
			t.Errorf("期望初始活跃连接数为0，实际为: %d", metrics["active_connections"])
		}
	})

	t.Run("并发获取推荐", func(t *testing.T) {
		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}
		defer pool.Close()

		var wg sync.WaitGroup
		concurrency := 10
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
				defer cancel()

				response, err := pool.GetListRecommend(ctx, 1) // 测试用户ID
				if err != nil {
					t.Errorf("协程 %d 获取推荐列表失败: %v", index, err)
					return
				}

				t.Logf("协程 %d 获取推荐列表成功: %+v", index, response)
			}(i)
		}

		wg.Wait()

		metrics := pool.GetMetrics()
		t.Logf("并发测试完成，最终指标: %+v", metrics)
		if metrics["error_count"].(int64) > 0 {
			t.Errorf("期望错误数为0，实际为: %d", metrics["error_count"])
		}
	})

	t.Run("连接超时处理", func(t *testing.T) {
		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}
		defer pool.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		_, err = pool.GetListRecommend(ctx, 1001)
		if err == nil {
			t.Error("期望超时错误，实际没有错误")
		}
	})

	t.Run("连接池指标监控", func(t *testing.T) {
		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}
		defer pool.Close()

		// 发送多个请求
		for i := 0; i < 5; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := pool.GetListRecommend(ctx, 1001)
			cancel()
			if err != nil {
				t.Errorf("请求失败: %v", err)
			}
		}

		metrics := pool.GetMetrics()
		if metrics["request_count"].(int64) != 5 {
			t.Errorf("期望请求数为5，实际为: %d", metrics["request_count"])
		}
	})

	t.Run("连接池关闭", func(t *testing.T) {
		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}

		if err := pool.Close(); err != nil {
			t.Errorf("关闭连接池失败: %v", err)
		}

		// 连接池关闭后应该无法发送请求
		ctx := context.Background()
		_, err = pool.GetListRecommend(ctx, 1001)
		if err == nil {
			t.Error("期望连接池关闭错误，实际没有错误")
		}
	})

	t.Run("长时间运行测试", func(t *testing.T) {
		if testing.Short() {
			t.Skip("跳过长时间运行测试")
		}

		pool, err := recommend.NewGRPCClientPool(cfg)
		if err != nil {
			t.Fatalf("创建连接池失败: %v", err)
		}
		defer pool.Close()

		// 运行5分钟，每10秒发送一次请求
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, err := pool.GetListRecommend(ctx, 1001)
				if err != nil {
					t.Errorf("请求失败: %v", err)
				}

				metrics := pool.GetMetrics()
				t.Logf("当前指标: %+v", metrics)
			}
		}
	})
}

func BenchmarkRecommendGRPCClientPool(b *testing.B) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	if err != nil {
		b.Fatalf("加载配置失败: %v", err)
	}

	pool, err := recommend.NewGRPCClientPool(cfg)
	if err != nil {
		b.Fatalf("创建连接池失败: %v", err)
	}
	defer pool.Close()

	b.Run("并发请求性能", func(b *testing.B) {
		b.SetParallelism(100) // 设置并发数
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctx := context.Background()
				_, err := pool.GetListRecommend(ctx, 1001)
				if err != nil {
					b.Errorf("请求失败: %v", err)
				}
			}
		})
	})
}
