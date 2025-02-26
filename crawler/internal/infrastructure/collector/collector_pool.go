package collector

import (
	"crawler/internal/infrastructure/config"
	"crawler/pkg/random"
	"math/rand/v2"
	"regexp"
	"sync"
	"time"

	"log"

	"github.com/gocolly/colly/v2"
)

// CollectorOptions 爬虫配置选项
type CollectorOptions struct {
	Domains     []string
	MaxDepth    int
	Timeout     time.Duration
	Parallelism int
	Async       bool
	MinDelay    time.Duration
	MaxDelay    time.Duration
}

// CollectorPool 爬虫池结构体
type CollectorPool struct {
	pool sync.Pool
	opts CollectorOptions
}

// NewCollectorPool 使用配置创建新的爬虫池
func NewCollectorPool(cfg *config.Config) *CollectorPool {
	opts := CollectorOptions{
		Domains:     cfg.Crawler.Domains,
		MaxDepth:    cfg.Crawler.MaxDepth,
		Timeout:     cfg.Crawler.Timeout,
		Parallelism: cfg.Crawler.Concurrency,
		Async:       cfg.Crawler.Async,
		MinDelay:    cfg.Crawler.MinDelay,
		MaxDelay:    cfg.Crawler.MaxDelay,
	}

	p := &CollectorPool{
		opts: opts,
		pool: sync.Pool{
			New: func() interface{} {
				c := colly.NewCollector(
					colly.UserAgent(random.GetRandomUserAgent()),
					colly.AllowedDomains(opts.Domains...),
					colly.Async(opts.Async),
					colly.MaxDepth(opts.MaxDepth),
					// colly.AllowURLRevisit(),
					colly.DisallowedURLFilters(
						regexp.MustCompile(`\.(jpg|jpeg|png|gif|ico|css|js)$`),
					),
				)

				// 设置请求超时时间和并发数
				c.SetRequestTimeout(opts.Timeout)
				c.Limit(&colly.LimitRule{
					DomainGlob:  "*",
					RandomDelay: time.Duration(float64(opts.MinDelay) + rand.Float64()*float64(opts.MaxDelay-opts.MinDelay)),
					Parallelism: opts.Parallelism,
				})

				// 添加通用请求头
				c.OnRequest(func(r *colly.Request) {
					r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
					r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
					r.Headers.Set("Cache-Control", "no-cache")
					r.Headers.Set("Pragma", "no-cache")
					r.Headers.Set("DNT", "1")
					r.Headers.Set("Upgrade-Insecure-Requests", "1")
					host := r.URL.Host
					r.Headers.Set("Referer", "https://"+host)
				})

				// 错误处理
				c.OnError(func(r *colly.Response, err error) {
					if r != nil {
						log.Printf("Error scraping %v: %v (Status code: %d)",
							r.Request.URL, err, r.StatusCode)
					} else {
						log.Printf("Error scraping: %v", err)
					}
				})

				return c
			},
		},
	}

	return p
}

// Get 从池中获取一个爬虫实例
func (p *CollectorPool) Get() *colly.Collector {
	return p.pool.Get().(*colly.Collector)
}

// 有问题的代码，不能清除之前的回调
// Put 将爬虫实例放回池中
func (p *CollectorPool) Put(c *colly.Collector) {
	// 使用空函数替代回调
	c.OnRequest(func(r *colly.Request) {})
	c.OnResponse(func(r *colly.Response) {})
	c.OnHTML("*", func(e *colly.HTMLElement) {})
	c.OnXML("*", func(e *colly.XMLElement) {})
	c.OnScraped(func(r *colly.Response) {})
	c.OnError(func(r *colly.Response, err error) {})
	c.SetProxyFunc(nil)
	p.pool.Put(c)
}
