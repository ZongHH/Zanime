package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 总配置结构
type Config struct {
	Server            ServerConfig                  `yaml:"server"`
	HTTP              HTTPConfig                    `yaml:"http"`
	GRPC              GRPCConfig                    `yaml:"grpc"`
	MySQL             MySQLConfig                   `yaml:"mysql"`
	Redis             RedisConfig                   `yaml:"redis"`
	NSQ               NSQConfig                     `yaml:"nsq"`
	TargetGrpcServers map[string]*GrpcServiceConfig `yaml:"target_grpc_servers"`
	JWT               JWTConfig                     `yaml:"jwt"`
	Cookie            CookieConfig                  `yaml:"cookie"`
	Storage           StorageConfig                 `yaml:"storage"`
	Security          SecurityConfig                `yaml:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

// GRPCConfig gRPC服务配置
type GRPCConfig struct {
	Host                  string        `yaml:"host"`
	Port                  int           `yaml:"port"`
	MaxRecvMsgSize        int           `yaml:"max_recv_msg_size"`
	MaxSendMsgSize        int           `yaml:"max_send_msg_size"`
	InitialWindowSize     int           `yaml:"initial_window_size"`
	InitialConnWindowSize int           `yaml:"initial_conn_window_size"`
	Keepalive             GRPCKeepalive `yaml:"keepalive"`
}

type GRPCKeepalive struct {
	Time                time.Duration `yaml:"time"`
	Timeout             time.Duration `yaml:"timeout"`
	PermitWithoutStream bool          `yaml:"permit_without_stream"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `yaml:"host"`              // MySQL服务器主机地址
	Port            int           `yaml:"port"`              // MySQL服务器端口号
	Username        string        `yaml:"username"`          // MySQL用户名
	Password        string        `yaml:"password"`          // MySQL密码
	Database        string        `yaml:"database"`          // 数据库名称
	Charset         string        `yaml:"charset"`           // 字符集设置
	ParseTime       bool          `yaml:"parse_time"`        // 是否解析时间
	Loc             string        `yaml:"loc"`               // 时区设置
	MaxIdleConns    int           `yaml:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int           `yaml:"max_open_conns"`    // 最大打开连接数
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"` // 连接最大生命周期
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `yaml:"host"`           // Redis服务器主机地址
	Port         int           `yaml:"port"`           // Redis服务器端口号
	Password     string        `yaml:"password"`       // Redis服务器密码
	DB           int           `yaml:"db"`             // Redis数据库编号
	PoolSize     int           `yaml:"pool_size"`      // 连接池大小
	MinIdleConns int           `yaml:"min_idle_conns"` // 最小空闲连接数
	DialTimeout  time.Duration `yaml:"dial_timeout"`   // 建立连接超时时间
	ReadTimeout  time.Duration `yaml:"read_timeout"`   // 读取超时时间
	WriteTimeout time.Duration `yaml:"write_timeout"`  // 写入超时时间
	PoolTimeout  time.Duration `yaml:"pool_timeout"`   // 从连接池获取连接的超时时间
	IdleTimeout  time.Duration `yaml:"idle_timeout"`   // 空闲连接超时时间
}

// NSQConfig NSQ配置
type NSQConfig struct {
	NSQD     NSQDConfig        `yaml:"nsqd"`
	Producer NSQProducerConfig `yaml:"producer"`
	Consumer NSQConsumerConfig `yaml:"consumer"`
}

type NSQDConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type NSQProducerConfig struct {
	PoolSize int `yaml:"pool_size"`
}

type NSQConsumerConfig struct {
	MaxInFlight        int `yaml:"max_in_flight"`
	ConcurrentHandlers int `yaml:"concurrent_handlers"`
	MaxAttempts        int `yaml:"max_attempts"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey    string         `yaml:"secret_key"`
	Issuer       string         `yaml:"issuer"`
	AccessToken  JWTTokenConfig `yaml:"access_token"`
	RefreshToken JWTTokenConfig `yaml:"refresh_token"`
	TokenType    string         `yaml:"token_type"`
}

type JWTTokenConfig struct {
	ExpireTime     time.Duration `yaml:"expire_time"`
	MaxRefreshTime time.Duration `yaml:"max_refresh_time,omitempty"`
}

// CookieConfig Cookie配置
type CookieConfig struct {
	Domain   string `yaml:"domain"`
	Path     string `yaml:"path"`
	MaxAge   int    `yaml:"max_age"`
	Secure   bool   `yaml:"secure"`
	HTTPOnly bool   `yaml:"http_only"`
	SameSite string `yaml:"same_site"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORS      CORSConfig      `yaml:"cors"`
	CSRF      CSRFConfig      `yaml:"csrf"`
	XSS       XSSConfig       `yaml:"xss"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type CORSConfig struct {
	AllowedOrigins   []string      `yaml:"allowed_origins"`
	AllowedMethods   []string      `yaml:"allowed_methods"`
	AllowedHeaders   []string      `yaml:"allowed_headers"`
	ExposedHeaders   []string      `yaml:"exposed_headers"`
	AllowCredentials bool          `yaml:"allow_credentials"`
	MaxAge           time.Duration `yaml:"max_age"`
}

type CSRFConfig struct {
	Enabled      bool     `yaml:"enabled"`
	ExcludePaths []string `yaml:"exclude_paths"`
}

type XSSConfig struct {
	Enabled bool `yaml:"enabled"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerSecond int  `yaml:"requests_per_second"`
}

// 单个gRPC服务配置
type GrpcServiceConfig struct {
	Enabled    bool        `yaml:"enabled"`      // 是否启用该服务
	Endpoints  []Endpoint  `yaml:"endpoints"`    // 服务端点列表
	Connection Connection  `yaml:"connection"`   // 连接参数配置
	Retry      RetryPolicy `yaml:"retry_policy"` // 重试策略配置
	Pool       PoolConfig  `yaml:"pool"`         // 连接池配置
}

// 端点地址配置
type Endpoint struct {
	Address string `yaml:"address"` // 服务地址(IP或域名)
	Port    int    `yaml:"port"`    // 服务端口
}

// 连接参数配置
type Connection struct {
	Timeout   time.Duration `yaml:"timeout"`   // 连接超时时间
	KeepAlive KeepAlive     `yaml:"keepalive"` // 保活策略配置
}

// 保活策略
type KeepAlive struct {
	Time                time.Duration `yaml:"time"`                  // 发送保活探测的时间间隔
	Timeout             time.Duration `yaml:"timeout"`               // 保活探测超时时间
	PermitWithoutStream bool          `yaml:"permit_without_stream"` // 是否允许在没有活动流的情况下发送保活探测
}

// 重试策略
type RetryPolicy struct {
	MaxAttempts       int           `yaml:"max_attempts"`       // 最大重试次数
	InitialBackoff    time.Duration `yaml:"initial_backoff"`    // 初始重试等待时间
	MaxBackoff        time.Duration `yaml:"max_backoff"`        // 最大重试等待时间
	BackoffMultiplier float64       `yaml:"backoff_multiplier"` // 重试等待时间增长因子
}

// 连接池配置
type PoolConfig struct {
	MaxConns     int           `yaml:"max_connections"`      // 最大连接数
	MaxIdleConns int           `yaml:"max_idle_connections"` // 最大空闲连接数
	MinIdleConns int           `yaml:"min_idle_connections"` // 最小空闲连接数
	MaxConnAge   time.Duration `yaml:"max_connection_age"`   // 连接最大存活时间
}

// StorageConfig 文件存储配置
type StorageConfig struct {
	Avatar AvatarConfig `yaml:"avatar"`     // 用户头像存储配置
	Post   PostConfig   `yaml:"post_image"` // 帖子图片存储配置
}

type AvatarConfig struct {
	Path         string   `yaml:"path"`          // 头像文件存储路径
	URL          string   `yaml:"url"`           // 头像文件访问URL
	MaxSize      int      `yaml:"max_size"`      // 最大文件大小(2MB)
	AllowedTypes []string `yaml:"allowed_types"` // 允许的文件类型
}

type PostConfig struct {
	Path         string   `yaml:"path"`          // 帖子图片存储路径
	URL          string   `yaml:"url"`           // 帖子图片访问URL
	MaxSize      int      `yaml:"max_size"`      // 最大文件大小(5MB)
	AllowedTypes []string `yaml:"allowed_types"` // 允许的文件类型
	MaxFiles     int      `yaml:"max_files"`     // 单个帖子最大图片数量
}

var globalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = config

	log.Printf("加载配置参数: \n"+
		"服务配置: 名称=%s 环境=%s 版本=%s\n"+
		"HTTP服务: %s:%d (超时: 读%v/写%v/空闲%v 最大头字节:%d)\n"+
		"gRPC服务: %s:%d (消息大小: 收%d/发%d 窗口: 流%d/连接%d 保活: 时间%v/超时%v/无流%v)\n"+
		"MySQL连接: %s:%s@tcp(%s:%d)/%s 字符集=%s 解析时间=%v 时区=%s 连接池: 空闲%d/最大%d/生命周期%v\n"+
		"Redis连接: %s:%d DB=%d 密码=%t 池: 大小%d/空闲%d 超时: 连接%v/读%v/写%v/池%v/空闲%v\n"+
		"NSQ配置: 节点%s:%d 生产者池=%d 消费者: 最大处理%d/并发%d/最大尝试%d\n"+
		"目标服务: 爬虫[启用=%t %s:%d] 推荐[启用=%t %s:%d]\n"+
		"JWT配置: 密钥长度=%d 签发者=%s 访问令牌[过期%v/最大刷新%v] 刷新令牌[过期%v] 类型=%s\n"+
		"Cookie配置: 域=%s 路径=%s 最大年龄=%d 安全=%t HTTPOnly=%t SameSite=%s\n"+
		"文件存储配置: 头像[路径=%s 最大大小=%d 类型%v] 帖子[路径=%s 最大大小=%d 类型%v 最大文件数=%d]\n"+
		"安全配置: CORS[源%v 方法%v 头%v 暴露头%v 凭证=%t 缓存%v] CSRF[启用=%t 排除%v] XSS[启用=%t] 限流[启用=%t %d/s]",
		config.Server.Name, config.Server.Env, config.Server.Version,
		config.HTTP.Host, config.HTTP.Port, config.HTTP.ReadTimeout, config.HTTP.WriteTimeout, config.HTTP.IdleTimeout, config.HTTP.MaxHeaderBytes,
		config.GRPC.Host, config.GRPC.Port, config.GRPC.MaxRecvMsgSize, config.GRPC.MaxSendMsgSize,
		config.GRPC.InitialWindowSize, config.GRPC.InitialConnWindowSize, config.GRPC.Keepalive.Time, config.GRPC.Keepalive.Timeout, config.GRPC.Keepalive.PermitWithoutStream,
		config.MySQL.Username, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database, config.MySQL.Charset,
		config.MySQL.ParseTime, config.MySQL.Loc, config.MySQL.MaxIdleConns, config.MySQL.MaxOpenConns, config.MySQL.ConnMaxLifetime,
		config.Redis.Host, config.Redis.Port, config.Redis.DB, config.Redis.Password != "",
		config.Redis.PoolSize, config.Redis.MinIdleConns, config.Redis.DialTimeout, config.Redis.ReadTimeout, config.Redis.WriteTimeout,
		config.Redis.PoolTimeout, config.Redis.IdleTimeout,
		config.NSQ.NSQD.Host, config.NSQ.NSQD.Port, config.NSQ.Producer.PoolSize,
		config.NSQ.Consumer.MaxInFlight, config.NSQ.Consumer.ConcurrentHandlers, config.NSQ.Consumer.MaxAttempts,
		config.TargetGrpcServers["scrape_service"].Enabled, config.TargetGrpcServers["scrape_service"].Endpoints[0].Address, config.TargetGrpcServers["scrape_service"].Endpoints[0].Port,
		config.TargetGrpcServers["recommend_service"].Enabled, config.TargetGrpcServers["recommend_service"].Endpoints[0].Address, config.TargetGrpcServers["recommend_service"].Endpoints[0].Port,
		len(config.JWT.SecretKey), config.JWT.Issuer, config.JWT.AccessToken.ExpireTime, config.JWT.AccessToken.MaxRefreshTime,
		config.JWT.RefreshToken.ExpireTime, config.JWT.TokenType,
		config.Cookie.Domain, config.Cookie.Path, config.Cookie.MaxAge, config.Cookie.Secure, config.Cookie.HTTPOnly, config.Cookie.SameSite,
		config.Storage.Avatar.Path, config.Storage.Avatar.MaxSize, config.Storage.Avatar.AllowedTypes,
		config.Storage.Post.Path, config.Storage.Post.MaxSize, config.Storage.Post.AllowedTypes, config.Storage.Post.MaxFiles,
		config.Security.CORS.AllowedOrigins, config.Security.CORS.AllowedMethods, config.Security.CORS.AllowedHeaders,
		config.Security.CORS.ExposedHeaders, config.Security.CORS.AllowCredentials, config.Security.CORS.MaxAge,
		config.Security.CSRF.Enabled, config.Security.CSRF.ExcludePaths,
		config.Security.XSS.Enabled,
		config.Security.RateLimit.Enabled, config.Security.RateLimit.RequestsPerSecond)

	return config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// GetMySQLDSN 获取MySQL连接字符串
func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		c.MySQL.Username,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database,
		c.MySQL.Charset,
		c.MySQL.ParseTime,
		c.MySQL.Loc,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetNSQDAddr 获取NSQD地址
func (c *Config) GetNSQDAddr() string {
	return fmt.Sprintf("%s:%d", c.NSQ.NSQD.Host, c.NSQ.NSQD.Port)
}

// GetHTTPAddr 获取HTTP服务地址
func (c *Config) GetHTTPAddr() string {
	return fmt.Sprintf("%s:%d", c.HTTP.Host, c.HTTP.Port)
}

// GetGRPCAddr 获取gRPC服务地址
func (c *Config) GetGRPCAddr() string {
	return fmt.Sprintf("%s:%d", c.GRPC.Host, c.GRPC.Port)
}

// GetScrapeAddr 获取爬虫服务器的地址
func (c *Config) GetScrapeAddr() string {
	return fmt.Sprintf("%s:%d", c.TargetGrpcServers["scrape_service"].Endpoints[0].Address, c.TargetGrpcServers["scrape_service"].Endpoints[0].Port)
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}
