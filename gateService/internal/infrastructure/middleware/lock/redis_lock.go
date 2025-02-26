package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLockAcquired = errors.New("锁已被占用")
	ErrLockExpired  = errors.New("锁已过期")
)

// RedisLock Redis分布式锁实现
type RedisLock struct {
	mutex      *redsync.Mutex
	autoExtend bool
	stopExtend chan struct{}
	opts       *LockOptions
}

// NewRedisLock 创建Redis分布式锁
func NewRedisLock(client *redis.Client, key string, opts *LockOptions) *RedisLock {
	if opts == nil {
		opts = &LockOptions{
			ExpireTime: 30 * time.Second,
			RetryCount: 3,
			RetryDelay: 200 * time.Millisecond,
			AutoExtend: true,
		}
	}

	// 创建redis连接池
	pool := goredis.NewPool(client)

	// 创建redsync实例
	rs := redsync.New(pool)

	// 创建互斥锁,设置以下参数:
	// - key: 锁的唯一标识
	// - WithExpiry: 锁的过期时间
	// - WithTries: 获取锁的重试次数
	// - WithRetryDelay: 重试间隔时间
	// - WithDriftFactor: 时钟漂移因子,用于补偿不同机器之间的时钟差异
	mutex := rs.NewMutex(key,
		redsync.WithExpiry(opts.ExpireTime),
		redsync.WithTries(opts.RetryCount),
		redsync.WithRetryDelay(opts.RetryDelay),
		redsync.WithDriftFactor(0.01), // 设置1%的时钟漂移补偿
	)

	return &RedisLock{
		mutex:      mutex,
		autoExtend: opts.AutoExtend,
		stopExtend: make(chan struct{}),
		opts:       opts,
	}
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	err := l.mutex.LockContext(ctx)
	if err != nil {
		return fmt.Errorf("获取锁失败: %w", err)
	}
	// 如果配置了自动续期，启动续期协程
	if l.autoExtend {
		go l.startAutoExtend(ctx)
	}
	return nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	if l.autoExtend {
		// 停止自动续期
		close(l.stopExtend)
	}

	// 释放锁
	ok, err := l.mutex.UnlockContext(ctx)
	if err != nil {
		return fmt.Errorf("释放锁失败: %w", err)
	}
	if !ok {
		return fmt.Errorf("锁已过期或被其他进程释放")
	}
	return nil
}

// Extend 延长锁时间
func (l *RedisLock) Extend(ctx context.Context) error {
	ok, err := l.mutex.ExtendContext(ctx)
	if err != nil {
		return fmt.Errorf("延长锁时间失败: %w", err)
	}
	if !ok {
		return fmt.Errorf("锁已过期或被其他进程获取")
	}
	return nil
}

// startAutoExtend 开启自动续期
func (l *RedisLock) startAutoExtend(ctx context.Context) {
	// 计算续期间隔，设置为过期时间的1/3
	interval := l.opts.ExpireTime / 3
	if interval <= 0 {
		interval = time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-l.stopExtend:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := l.Extend(ctx); err != nil {
				// 续期失败，停止续期
				return
			}
		}
	}
}

// WaitLock 一直等待直到获取到锁或者上下文取消
func (l *RedisLock) WaitLock(ctx context.Context) error {
	for {
		// 尝试获取锁
		err := l.mutex.LockContext(ctx)
		if err == nil {
			// 获取锁成功
			if l.autoExtend {
				go l.startAutoExtend(ctx)
			}
			return nil
		}

		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.opts.RetryDelay):
			// 继续尝试
			continue
		}
	}
}
