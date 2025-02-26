package database

import (
	"gateService/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
)

type RDB struct {
	rdb *redis.Client
}

// NewRDB 创建新的Redis客户端连接
func NewRDB(cfg *config.Config) *RDB {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),     // Redis服务器地址,格式为host:port
		Password:     cfg.Redis.Password,     // Redis服务器密码
		DB:           cfg.Redis.DB,           // 使用的数据库编号
		PoolSize:     cfg.Redis.PoolSize,     // 连接池最大连接数
		PoolTimeout:  cfg.Redis.PoolTimeout,  // 从连接池获取连接的超时时间
		ReadTimeout:  cfg.Redis.ReadTimeout,  // 读取超时时间
		WriteTimeout: cfg.Redis.WriteTimeout, // 写入超时时间
		DialTimeout:  cfg.Redis.DialTimeout,  // 建立连接超时时间
		MinIdleConns: cfg.Redis.MinIdleConns, // 最小空闲连接数
	})
	return &RDB{rdb: client}
}

func (r *RDB) GetRDB() *redis.Client {
	return r.rdb
}

func (r *RDB) Close() {
	r.rdb.Close()
}
