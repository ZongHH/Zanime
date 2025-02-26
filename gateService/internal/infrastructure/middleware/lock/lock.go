package lock

import (
	"context"
	"time"
)

// DistributedLock 分布式锁接口
type DistributedLock interface {
	// Lock 获取锁
	Lock(ctx context.Context) error
	// WaitLock 一直等待直到获取到锁
	WaitLock(ctx context.Context) error
	// Unlock 释放锁
	Unlock(ctx context.Context) error
	// Extend 延长锁时间
	Extend(ctx context.Context) error
}

// LockOptions 锁配置选项
type LockOptions struct {
	// 锁的过期时间
	ExpireTime time.Duration
	// 重试次数
	RetryCount int
	// 重试延迟
	RetryDelay time.Duration
	// 是否自动续期
	AutoExtend bool
}
