package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	App            AppConfig            `mapstructure:"app"`
	ManagerService ManagerServiceConfig `mapstructure:"managerService"`
	MySQL          MySQLConfig          `mapstructure:"mysql"`
	Redis          RedisConfig          `mapstructure:"redis"`
	NSQ            NSQConfig            `mapstructure:"nsq"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

type ManagerServiceConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderBytes  int           `mapstructure:"maxHeaderBytes"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

type NSQConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Auth     NSQAuthConfig `mapstructure:"auth"`
	Producer NSQProducer   `mapstructure:"producer"`
	Consumer NSQConsumer   `mapstructure:"consumer"`
}

type NSQAuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type NSQProducer struct {
	MaxInFlight       int           `mapstructure:"maxInFlight"`
	ReconnectInterval time.Duration `mapstructure:"reconnectInterval"`
}

type NSQConsumer struct {
	MaxInFlight       int           `mapstructure:"maxInFlight"`
	Concurrency       int           `mapstructure:"concurrency"`
	ReconnectInterval time.Duration `mapstructure:"reconnectInterval"`
}

type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parseTime"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime"`
	SlowThreshold   time.Duration `mapstructure:"slowThreshold"`
	LogLevel        string        `mapstructure:"logLevel"`
}

type RedisConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	Password           string        `mapstructure:"password"`
	Database           int           `mapstructure:"database"`
	DialTimeout        time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	PoolSize           int           `mapstructure:"poolSize"`
	MinIdleConns       int           `mapstructure:"minIdleConns"`
	MaxConnAge         time.Duration `mapstructure:"maxConnAge"`
	PoolTimeout        time.Duration `mapstructure:"poolTimeout"`
	IdleTimeout        time.Duration `mapstructure:"idleTimeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idleCheckFrequency"`
}

var (
	GlobalConfig Config
	v            *viper.Viper
)

// LoadConfig 初始化配置
func LoadConfig(configPath string) (*Config, error) {
	var err error
	v = viper.New()

	// 设置默认值
	setDefaults()

	// 配置 Viper
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("Manager")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err = v.Unmarshal(&GlobalConfig); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err = validateConfig(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &GlobalConfig, nil
}

// setDefaults 设置默认配置
func setDefaults() {
	v.SetDefault("app.name", "Manager-service")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "development")

	v.SetDefault("managerService.host", "127.0.0.1")
	v.SetDefault("managerService.port", 9999)
	v.SetDefault("managerService.readTimeout", "5s")
	v.SetDefault("managerService.writeTimeout", "5s")
	v.SetDefault("managerService.maxHeaderBytes", 1048576)
	v.SetDefault("managerService.shutdownTimeout", "10s")

	v.SetDefault("mysql.host", "127.0.0.1")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.charset", "utf8mb4")
	v.SetDefault("mysql.parseTime", true)
	v.SetDefault("mysql.loc", "Local")
	v.SetDefault("mysql.maxIdleConns", 10)
	v.SetDefault("mysql.maxOpenConns", 100)
	v.SetDefault("mysql.connMaxLifetime", "3600s")
	v.SetDefault("mysql.connMaxIdleTime", "3600s")
	v.SetDefault("mysql.slowThreshold", "200ms")
	v.SetDefault("mysql.logLevel", "info")

	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.dialTimeout", "5s")
	v.SetDefault("redis.readTimeout", "3s")
	v.SetDefault("redis.writeTimeout", "3s")
	v.SetDefault("redis.poolSize", 10)
	v.SetDefault("redis.minIdleConns", 5)
	v.SetDefault("redis.maxConnAge", "3600s")
	v.SetDefault("redis.poolTimeout", "1s")
	v.SetDefault("redis.idleTimeout", "300s")
	v.SetDefault("redis.idleCheckFrequency", "60s")
}

// validateConfig 验证配置
func validateConfig() error {
	if GlobalConfig.App.Env != "development" &&
		GlobalConfig.App.Env != "testing" &&
		GlobalConfig.App.Env != "production" {
		return fmt.Errorf("无效的环境配置: %s", GlobalConfig.App.Env)
	}

	if GlobalConfig.ManagerService.Port <= 0 || GlobalConfig.ManagerService.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", GlobalConfig.ManagerService.Port)
	}

	if GlobalConfig.MySQL.Port <= 0 || GlobalConfig.MySQL.Port > 65535 {
		return fmt.Errorf("无效的 MySQL 端口号: %d", GlobalConfig.MySQL.Port)
	}

	if GlobalConfig.Redis.Port <= 0 || GlobalConfig.Redis.Port > 65535 {
		return fmt.Errorf("无效的 Redis 端口号: %d", GlobalConfig.Redis.Port)
	}

	return nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return &GlobalConfig
}

// WatchConfig 监听配置文件变化
func WatchConfig(onChange func()) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&GlobalConfig); err != nil {
			fmt.Printf("配置热重载失败: %v\n", err)
			return
		}
		if onChange != nil {
			onChange()
		}
	})
}

func (cfg *Config) GetManagerAddr() string {
	return fmt.Sprintf("%s:%d",
		cfg.ManagerService.Host, cfg.ManagerService.Port)
}

// GetMySQLDSN 获取MySQL连接DSN
func (cfg *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
		cfg.MySQL.Charset,
		cfg.MySQL.ParseTime,
		cfg.MySQL.Loc)
}

// GetRedisAddr 获取Redis地址
func (cfg *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
}

// GetNSQAddr 获取NSQ地址
func (cfg *Config) GetNSQAddr() string {
	return fmt.Sprintf("%s:%d", cfg.NSQ.Host, cfg.NSQ.Port)
}

// GetMySQLConnParams 获取MySQL连接池参数
func (cfg *Config) GetMySQLConnParams() (maxIdleConns, maxOpenConns int, maxLifetime, maxIdleTime time.Duration) {
	return cfg.MySQL.MaxIdleConns,
		cfg.MySQL.MaxOpenConns,
		cfg.MySQL.ConnMaxLifetime,
		cfg.MySQL.ConnMaxIdleTime
}

// GetRedisOptions 获取Redis选项
func (cfg *Config) GetRedisOptions() map[string]interface{} {
	return map[string]interface{}{
		"addr":               cfg.GetRedisAddr(),
		"password":           cfg.Redis.Password,
		"db":                 cfg.Redis.Database,
		"dialTimeout":        cfg.Redis.DialTimeout,
		"readTimeout":        cfg.Redis.ReadTimeout,
		"writeTimeout":       cfg.Redis.WriteTimeout,
		"poolSize":           cfg.Redis.PoolSize,
		"minIdleConns":       cfg.Redis.MinIdleConns,
		"maxConnAge":         cfg.Redis.MaxConnAge,
		"poolTimeout":        cfg.Redis.PoolTimeout,
		"idleTimeout":        cfg.Redis.IdleTimeout,
		"idleCheckFrequency": cfg.Redis.IdleCheckFrequency,
	}
}

// GetAppInfo 获取应用信息
func (cfg *Config) GetAppInfo() (name, version, env string) {
	return cfg.App.Name, cfg.App.Version, cfg.App.Env
}

// GetServerTimeouts 获取服务超时设置
func (cfg *Config) GetServerTimeouts() (read, write, shutdown time.Duration) {
	return cfg.ManagerService.ReadTimeout,
		cfg.ManagerService.WriteTimeout,
		cfg.ManagerService.ShutdownTimeout
}

// IsDevMode 是否为开发环境
func IsDevMode() bool {
	return GlobalConfig.App.Env == "development"
}

// IsTestMode 是否为测试环境
func IsTestMode() bool {
	return GlobalConfig.App.Env == "testing"
}

// IsProdMode 是否为生产环境
func IsProdMode() bool {
	return GlobalConfig.App.Env == "production"
}
