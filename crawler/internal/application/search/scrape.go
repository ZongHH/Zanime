package search

import (
	"context"
	"crawler/internal/infrastructure/config"
	"crawler/pkg/random"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ScraperConfig 爬虫配置结构体
// 用于配置爬虫的各项行为参数，包括超时设置、重试机制和浏览器选项等
type ScraperConfig struct {
	// Timeout 单次请求的最大超时时间
	// - 如果请求超过此时间将会被中断并触发重试机制
	// - 设置过短可能导致复杂页面加载失败
	// - 设置过长会导致错误处理延迟
	// - 建议值: 15-30秒，具体取决于目标网站的响应速度
	// - 如果目标网站加载速度慢，可适当增加
	Timeout time.Duration

	// RetryCount 请求失败后的最大重试次数
	// - 当请求失败（超时、网络错误等）时会重试指定次数
	// - 设置过大会增加总体执行时间
	// - 设置过小可能降低成功率
	// - 建议值: 1-3次，避免过多重试占用资源
	// - 需要考虑目标网站的反爬虫机制
	RetryCount int

	// RetryInterval 重试之间的时间间隔
	// - 连续重试之间等待的时间，避免频繁请求
	// - 间隔过短可能触发目标网站的频率限制
	// - 间隔过长会影响整体效率
	// - 建议值: 1-3秒，防止对目标服务器造成压力
	// - 建议配合RandomDelay使用增加随机性
	RetryInterval time.Duration

	// RandomDelay 随机延迟时间上限
	// - 每次请求之前会随机等待一个不超过此值的时间
	// - 用于模拟真实用户行为，避免被反爬虫机制检测
	// - 设置过短可能容易被识别为机器人
	// - 设置过长会显著降低爬取效率
	// - 建议值: 1-3秒，根据实际需要调整
	// - 实际延迟时间是0到此值之间的随机数
	RandomDelay time.Duration

	// ChromeOptions Chrome浏览器的启动参数
	// - 用于配置浏览器的行为和特性
	// - 包含的常用选项：
	//   * enable-automation: 是否启用自动化标志
	//   * disable-blink-features: 禁用特定的blink特性
	//   * headless: 无界面模式运行
	//   * disable-web-security: 禁用同源策略
	//   * disable-background-networking: 禁用后台网络活动
	//   * disable-background-timer-throttling: 禁用计时器限制
	//   * disable-backgrounding-occluded-windows: 禁用窗口遮挡时的后台处理
	//   * no-sandbox: 禁用沙箱模式（注意安全性）
	// - 这些选项会影响:
	//   * 浏览器的性能
	//   * 内存使用
	//   * 安全性
	//   * 反爬虫检测
	ChromeOptions []chromedp.ExecAllocatorOption
}

// DefaultConfig 默认的爬虫配置
// 提供了一组经过优化的默认参数值，适用于大多数爬取场景
// 这些配置已经过实际测试，可以在保证性能的同时避免被反爬虫机制检测
var DefaultConfig = ScraperConfig{
	// Timeout 设置10秒超时
	// - 足够处理大多数页面加载情况
	// - 包括DOM渲染和JavaScript执行时间
	// - 对于网络较慢的环境也有足够余量
	Timeout: 10 * time.Second,

	// RetryCount 设置3次重试
	// - 降低对目标服务器的压力
	// - 避免频繁请求触发反爬虫机制
	// - 失败后重试最多3次，快速失败
	RetryCount: 3,

	// RetryInterval 重试间隔1秒
	// - 给予服务器足够的恢复时间
	// - 模拟人类手动重试的间隔
	// - 避免被识别为机械性重试
	RetryInterval: 1 * time.Second,

	// RandomDelay 随机延迟最多1毫秒
	// - 增加请求的随机性
	// - 模拟真实用户行为
	// - 平衡性能和反爬虫需求
	RandomDelay: 1 * time.Millisecond, // 暂未使用

	// ChromeOptions Chrome浏览器启动参数
	// 精心调教的浏览器配置，兼顾性能和隐蔽性
	ChromeOptions: []chromedp.ExecAllocatorOption{
		// 禁用自动化标志
		// 防止网站检测到自动化测试工具
		chromedp.Flag("enable-automation", false),

		// 禁用自动化控制特性
		// 移除可能暴露浏览器自动化的特征
		chromedp.Flag("disable-blink-features", "AutomationControlled"),

		// 启用无界面模式
		// 在后台运行浏览器，节省资源
		chromedp.Headless,

		// 禁用Web安全策略
		// 允许跨域请求，访问更多资源
		chromedp.Flag("disable-web-security", true),

		// 禁用软件光栅化
		// 提高性能，减少资源占用
		chromedp.Flag("disable-software-rasterizer", true),

		// 禁用沙箱模式
		// 提高性能，但会降低安全性
		// 在可控环境中使用
		chromedp.Flag("no-sandbox", true),

		// 禁用后台网络活动，减少不必要的网络请求
		chromedp.Flag("disable-background-networking", true),

		// 禁用后台计时器限制，提高性能
		chromedp.Flag("disable-background-timer-throttling", true),

		// 禁用窗口遮挡时的后台处理，提高性能
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
	},
}

// VideoScraper 视频爬虫结构体
type VideoScraper struct {
	config *ScraperConfig
}

// NewVideoScraper 创建新的视频爬虫实例
func NewVideoScraper(config *config.DynamicCrawlerConfig) *VideoScraper {
	if config != nil {
		DefaultConfig.Timeout = config.Timeout
		DefaultConfig.RetryCount = config.RetryCount
		DefaultConfig.RetryInterval = config.RetryInterval
		DefaultConfig.RandomDelay = config.RandomDelay
	}

	return &VideoScraper{config: &DefaultConfig}
}

// createContext 创建 Chromedp 上下文
func (s *VideoScraper) createContext(ctx context.Context) (context.Context, []context.CancelFunc) {
	// 收集所有需要取消的函数
	var cancels []context.CancelFunc

	// 创建一个新的分配器
	opts := append(chromedp.DefaultExecAllocatorOptions[:], s.config.ChromeOptions...)
	opts = append(opts, chromedp.UserAgent(random.GetRandomUserAgent()))
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	cancels = append(cancels, allocCancel)

	// 创建一个新的Chrome实例
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	cancels = append(cancels, taskCancel)

	taskCtxWithTime, taskCtxWithTimeCancel := context.WithTimeout(taskCtx, s.config.Timeout)
	cancels = append(cancels, taskCtxWithTimeCancel)

	return taskCtxWithTime, cancels
}

func (s *VideoScraper) cancelContexts(cancels []context.CancelFunc) {
	// 按照反向顺序清理资源
	for i := len(cancels) - 1; i >= 0; i-- {
		cancels[i]()
	}
}

// ScrapeVideo 抓取视频URL
func (s *VideoScraper) ScrapeVideo(url string) (string, error) {
	// 创建新的上下文和取消函数列表
	ctxChromedp, cancels := s.createContext(context.Background())
	defer s.cancelContexts(cancels)

	urlChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// 设置网络请求监听器
	s.setupNetworkListener(ctxChromedp, urlChan)

	// 执行页面导航
	go s.navigateToPage(ctxChromedp, url, errChan)

	return s.waitForResult(ctxChromedp, urlChan, errChan)
}

// setupNetworkListener 设置网络请求监听
func (s *VideoScraper) setupNetworkListener(ctx context.Context, urlChan chan<- string) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*network.EventRequestWillBeSent); ok {
			if strings.Contains(e.Request.URL, "m3u8") {
				select {
				case urlChan <- e.Request.URL:
				default:
				}
			}
		}
	})
}

// navigateToPage 导航到目标页面
func (s *VideoScraper) navigateToPage(ctx context.Context, url string, errChan chan<- error) {
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	)
	if err != nil {
		select {
		case errChan <- err:
		default:
		}
	}
}

// waitForResult 等待抓取结果
func (s *VideoScraper) waitForResult(ctx context.Context, urlChan <-chan string, errChan <-chan error) (string, error) {
	select {
	case videoUrl := <-urlChan:
		return s.scrapeVideoUrl(videoUrl)
	case err := <-errChan:
		return "", fmt.Errorf("获取m3u8链接失败: %v", err)
	case <-ctx.Done():
		return "", fmt.Errorf("scrapeVideoWithRetry方法超时: %v", ctx.Err())
	}
}

// scrapeVideoUrl 抓取视频URL
// 该函数的作用是通过给定的URL导航到视频页面，并提取视频元素的源URL。
// 它使用chromedp库来控制浏览器，执行一系列操作以确保页面加载完成并且视频元素可见。
// 最终，它将视频的URL返回给调用者。
func (s *VideoScraper) scrapeVideoUrl(url string) (string, error) {
	log.Printf("开始监听m3u8链接: %v\n", url)

	for i := 0; i < s.config.RetryCount; i++ {
		videoURL, err := s.scrapeVideoUrlWithRetry(url)
		if err != nil {
			log.Printf("第%d次尝试失败: %v\n", i+1, err)
		} else {
			return videoURL, nil
		}
		if i != s.config.RetryCount-1 {
			time.Sleep(s.config.RetryInterval)
		}
	}
	return "", fmt.Errorf("%d次尝试全部失败", s.config.RetryCount)
}

func (s *VideoScraper) scrapeVideoUrlWithRetry(url string) (string, error) {
	taskCtx, cancels := s.createContext(context.Background())
	defer s.cancelContexts(cancels)

	var videoUrl string
	err := chromedp.Run(taskCtx,
		network.Enable(),                                                  // 启用网络请求监听
		chromedp.Navigate(url),                                            // 导航到指定的URL
		chromedp.WaitReady("body", chromedp.ByQuery),                      // 等待页面的body元素准备好
		chromedp.WaitReady("video#lelevideo", chromedp.ByQuery),           // 等待视频元素准备好
		chromedp.AttributeValue("video#lelevideo", "src", &videoUrl, nil), // 获取视频元素的src属性值
	)

	if err != nil {
		if taskCtx.Err() != nil {
			chromedp.Cancel(taskCtx)
			return "", fmt.Errorf("chromedp.Run访问m3u8链接超时: %v", err) // 如果上下文超时，返回超时错误
		}
		return "", fmt.Errorf("chromedp.Run访问m3u8链接失败: %v", err) // 返回提取视频URL失败的错误
	}

	if videoUrl == "" {
		return "", fmt.Errorf("视频链接抓取为空") // 如果没有找到视频URL，返回错误
	}

	return videoUrl, nil // 返回提取到的视频URL
}
