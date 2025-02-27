package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构
// 包含了应用程序所需的所有配置信息,包括gRPC服务器、MySQL数据库、爬虫、NSQ消息队列和URL等配置
type Config struct {
	GRPCServer     GRPCServerConfig     `mapstructure:"gRpcServer"`      // gRPC服务器配置
	MySQL          MySQLConfig          `mapstructure:"mysql"`           // MySQL数据库配置
	Crawler        CrawlerConfig        `mapstructure:"crawler"`         // 爬虫配置
	NSQ            NSQConfig            `mapstructure:"nsq"`             // NSQ消息队列配置
	URLs           URLsConfig           `mapstructure:"urls"`            // URL相关配置
	DynamicCrawler DynamicCrawlerConfig `mapstructure:"dynamic_crawler"` // 动态链接爬虫配置
}

// DynamicCrawlerConfig 动态链接爬虫配置
// 定义了动态链接爬虫的相关参数
type DynamicCrawlerConfig struct {
	Timeout       time.Duration `mapstructure:"timeout"`        // 默认超时时间
	RetryCount    int           `mapstructure:"retry_count"`    // 默认重试次数
	RetryInterval time.Duration `mapstructure:"retry_interval"` // 默认重试间隔
	RandomDelay   time.Duration `mapstructure:"random_delay"`   // 默认随机延迟
}

// GRPCServerConfig gRPC服务器配置
// 定义了gRPC服务器的监听地址和端口
type GRPCServerConfig struct {
	Host string `mapstructure:"host"` // 服务器监听地址,如"0.0.0.0"表示监听所有网卡
	Port int    `mapstructure:"port"` // 服务器监听端口,如9094
}

// MySQLConfig MySQL配置
// 包含了MySQL数据库连接和连接池的所有相关配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`              // 数据库主机地址
	Port            int           `mapstructure:"port"`              // 数据库端口号
	User            string        `mapstructure:"user"`              // 数据库用户名
	Password        string        `mapstructure:"password"`          // 数据库密码
	Database        string        `mapstructure:"database"`          // 数据库名称
	Charset         string        `mapstructure:"charset"`           // 字符集,如utf8mb4
	ParseTime       bool          `mapstructure:"parse_time"`        // 是否解析时间类型
	Loc             string        `mapstructure:"loc"`               // 时区设置,如Asia/Shanghai
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 连接池中的最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 连接池最大连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接的最大生命周期
}

// CrawlerConfig 爬虫配置
// 定义了爬虫的行为参数,包括并发控制、延迟设置等
type CrawlerConfig struct {
	Domains     []string      `mapstructure:"domains"`     // 允许爬取的域名列表,用于限制爬虫范围
	Concurrency int           `mapstructure:"concurrency"` // 并发爬虫数量,控制同时运行的爬虫协程数
	Async       bool          `mapstructure:"async"`       // 是否启用异步模式,true表示使用异步爬取
	MaxDepth    int           `mapstructure:"max_depth"`   // 最大爬取深度,控制网页递归爬取的层级
	Timeout     time.Duration `mapstructure:"timeout"`     // 请求超时时间(秒),控制单个请求的最长等待时间
	MinDelay    time.Duration `mapstructure:"min_delay"`   // 最小请求延迟(毫秒),控制请求间隔的下限
	MaxDelay    time.Duration `mapstructure:"max_delay"`   // 最大请求延迟(毫秒),控制请求间隔的上限
	Interval    time.Duration `mapstructure:"interval"`    // 爬虫更新数据间隔,控制爬虫更新数据的时间间隔
}

// NSQConfig NSQ配置
// 包含了NSQ消息队列的连接和行为配置
type NSQConfig struct {
	Host string `mapstructure:"host"` // NSQ服务器地址
	Port int    `mapstructure:"port"` // NSQ服务器端口
}

// URLsConfig URL配置
// 存储各个网站的URL模板和基础URL
type URLsConfig struct {
	YingHuaDongMan struct {
		Search string `mapstructure:"search"` // 搜索页面URL模板,用于构建搜索请求
		Main   string `mapstructure:"main"`   // 主页URL,网站的基础URL
	} `mapstructure:"yinghuadongman"`
}

// 全局配置实例
var cfg *Config

// LoadConfig 加载配置
// 从指定路径加载配置文件并解析到全局配置实例中
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv() // 允许使用环境变量覆盖配置

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return cfg, nil
}

// GetConfig 获取配置
// 返回全局配置实例,如果配置未初始化则触发panic
func GetConfig() *Config {
	if cfg == nil {
		panic("配置未初始化")
	}
	return cfg
}

// GetMySQLDSN 获取MySQL连接字符串
// 根据MySQL配置生成标准的DSN(Data Source Name)连接字符串
func (c *Config) GetMySQLDSN() string {
	// 转义时区参数
	timezone := url.QueryEscape("Asia/Shanghai")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database,
		timezone,
	)

	return dsn
}

// GetGRPCAddress 获取gRPC服务地址
// 返回格式化的gRPC服务器地址,格式为"host:port"
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.GRPCServer.Host, c.GRPCServer.Port)
}

// GetNSQAddress 获取NSQ地址
// 返回格式化的NSQ服务器地址,格式为"host:port"
func (c *Config) GetNSQAddress() string {
	return fmt.Sprintf("%s:%d", c.NSQ.Host, c.NSQ.Port)
}

// GetSearchURL 获取搜索URL
// 返回配置的搜索页面URL模板
func (c *Config) GetSearchURL() string {
	return c.URLs.YingHuaDongMan.Search
}

// GetMainURL 获取主站URL
// 返回配置的网站主页URL
func (c *Config) GetMainURL() string {
	return c.URLs.YingHuaDongMan.Main
}

// GetCrawlerTimeout 获取爬虫超时时间
// 将配置中的超时秒数转换为time.Duration类型
func (c *Config) GetCrawlerTimeout() time.Duration {
	return time.Duration(c.Crawler.Timeout) * time.Second
}

// GetCrawlerDelay 获取爬虫延迟范围
// 返回配置的最小和最大延迟时间,用于控制爬虫请求间隔
func (c *Config) GetCrawlerDelay() (time.Duration, time.Duration) {
	return time.Duration(c.Crawler.MinDelay) * time.Millisecond,
		time.Duration(c.Crawler.MaxDelay) * time.Millisecond
}
