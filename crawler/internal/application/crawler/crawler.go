package crawler

import (
	"context"
	"crawler/internal/domain/entity"
	"crawler/internal/domain/repository"
	"crawler/internal/infrastructure/collector"
	"crawler/internal/infrastructure/config"
	"crawler/pkg/monitor"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/v2"
)

// AnimeCrawler 动漫爬虫结构体
// 负责爬取、解析和存储动漫相关数据
// 支持定时任务和增量更新功能
type AnimeCrawler struct {
	// 应用配置信息
	config *config.Config
	// 视频数据仓储接口
	videoRepo repository.VideoRepository
	// 爬虫收集器池
	collectorPool *collector.CollectorPool

	// 需要爬取的动漫分类名称列表
	categorieName []string

	// 主爬虫收集器，用于爬取首页内容
	collector *colly.Collector
	// 分类爬虫收集器，用于爬取分类页面
	categories *colly.Collector
	// 动漫库爬虫收集器，用于爬取动漫列表
	library *colly.Collector
	// 详情爬虫收集器，用于爬取动漫详细信息
	details *colly.Collector

	// 定时器，控制爬虫执行间隔
	ticker *time.Ticker
	// 爬虫启动次数计数器
	StartCount int32
	// 本次爬取更新的动漫总数
	// 使用原子操作确保并发安全
	TotalUpdateCount int32
}

// NewAnimeCrawler 创建新的动漫爬虫
func NewAnimeCrawler(cfg *config.Config, collectorPool *collector.CollectorPool, videoRepo repository.VideoRepository) *AnimeCrawler {
	return &AnimeCrawler{
		config:        cfg,
		videoRepo:     videoRepo,
		collectorPool: collectorPool,
		ticker:        time.NewTicker(cfg.Crawler.Interval),
		categorieName: []string{"日本动漫", "国产动漫", "欧美动漫", "动漫电影"},
	}
}

// Start 启动爬虫
func (c *AnimeCrawler) Start(ctx context.Context) {
	monitor.Info("启动动漫爬虫...")

	// 定时执行
	for {
		select {
		case <-ctx.Done():
			monitor.Info("停止动漫爬虫")
			c.ticker.Stop()
			return
		case <-c.ticker.C:
			c.crawl(ctx)
		}
	}
}

// crawl 执行爬虫任务
func (c *AnimeCrawler) crawl(ctx context.Context) {
	atomic.AddInt32(&c.StartCount, 1)
	monitor.Info("第%d次更新爬取开始", atomic.LoadInt32(&c.StartCount))

	errChan := make(chan error, 1)
	done := make(chan struct{})

	// 开始爬取
	go func() {
		c.startCrawling(errChan)
		// 当爬虫完成时关闭 done channel
		close(done)
	}()

	// 处理错误，直到爬虫完成或上下文取消
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			monitor.Error("爬虫错误: %v", err)
		case <-done:
			monitor.Info("第%d次更新爬取完成，共更新%d部动漫", atomic.LoadInt32(&c.StartCount), atomic.LoadInt32(&c.TotalUpdateCount))
			c.resetCount()
			return
		}
	}
}

// setupHandlers 设置所有处理器
func (c *AnimeCrawler) setupHandlers(errChan chan<- error) {
	c.setupCollectorHandler()
	c.setupCategoryHandler()
	c.setupLibraryHandler()
	c.setupDetailHandler(errChan)
	c.setupErrorHandlers(errChan)
}

// setupCollectorHandler 设置主页处理器
func (c *AnimeCrawler) setupCollectorHandler() {
	c.collector.OnHTML("ul.navbar-items.swiper-wrapper", func(h *colly.HTMLElement) {
		for _, category := range c.categorieName {
			href := h.ChildAttr(fmt.Sprintf("a[title='%s']", category), "href")
			if href != "" {
				absoluteURL := h.Request.AbsoluteURL(href)
				c.categories.Visit(absoluteURL)
			}
		}
	})
}

// setupCategoryHandler 设置分类处理器
func (c *AnimeCrawler) setupCategoryHandler() {
	c.categories.OnHTML("a.module-heading-tab-link", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		absoluteURL := h.Request.AbsoluteURL(href)
		c.library.Visit(absoluteURL)
	})
}

// setupLibraryHandler 设置库处理器
func (c *AnimeCrawler) setupLibraryHandler() {
	c.library.OnHTML("div.page.list.pianku a.module-poster-item.module-item", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		absoluteURL := h.Request.AbsoluteURL(href)
		c.details.Visit(absoluteURL)
	})
	c.library.OnHTML("div.page.list.pianku a[title='下一页']", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		absoluteURL := h.Request.AbsoluteURL(href)
		c.library.Visit(absoluteURL)
	})
}

// setupDetailHandler 设置详情处理器
func (c *AnimeCrawler) setupDetailHandler(errChan chan<- error) {
	c.details.OnHTML("div.page.view", func(h *colly.HTMLElement) {
		anime, genres, urls := c.extractAnimeInfo(h)

		if err := c.saveAnimeData(anime, genres, urls); err != nil {
			select {
			case errChan <- fmt.Errorf("保存动漫数据失败: %v", err):
			default:
			}
			return
		}

		monitor.Info("更新动漫:name %s, realease %s, area %s\n", anime.VideoName, anime.ReleaseDate, anime.Area)

		atomic.AddInt32(&c.TotalUpdateCount, 1)
		if atomic.LoadInt32(&c.TotalUpdateCount)%100 == 0 {
			monitor.Info("第%d次更新已经更新到第%d部动漫", atomic.LoadInt32(&c.StartCount), atomic.LoadInt32(&c.TotalUpdateCount))
		}
	})
}

// extractAnimeInfo 提取动漫信息
func (c *AnimeCrawler) extractAnimeInfo(h *colly.HTMLElement) (*entity.Video, *[]entity.AnimeGenre, *[]entity.VideoUrl) {
	anime := entity.Video{
		UploaderID: 1,
	}

	// 提取基本信息
	imageUrl := h.ChildAttr("img.ls-is-cached.lazy.lazyload", "data-original")
	anime.CoverImageUrl = h.Request.AbsoluteURL(imageUrl)

	var genres []string
	h.ForEach("div.module-info-heading", func(_ int, h *colly.HTMLElement) {
		anime.VideoName = h.ChildText("h1")
		h.ForEach("div.module-info-tag-link", func(i int, h *colly.HTMLElement) {
			switch i {
			case 0:
				anime.ReleaseDate = h.Text
			case 1:
				anime.Area = h.Text
			default:
				h.ForEach("a", func(_ int, h *colly.HTMLElement) {
					genres = append(genres, h.Text)
				})
			}
		})
	})

	anime.Description = h.ChildText("div.module-info-introduction-content.show-desc")

	// 提取剧集信息
	var animeUrls []entity.VideoUrl
	h.ForEach("div#panel1", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			h.ForEach("a", func(_ int, h *colly.HTMLElement) {
				url := h.Request.AbsoluteURL(h.Attr("href"))
				animeUrls = append(animeUrls, entity.VideoUrl{
					Episode:  h.Text,
					VideoUrl: url,
				})
			})
		}
	})

	// 构建分类信息
	var animeGenres []entity.AnimeGenre
	for _, genre := range genres {
		animeGenres = append(animeGenres, entity.AnimeGenre{
			Genre: genre,
		})
	}

	return &anime, &animeGenres, &animeUrls
}

// saveAnimeData 保存动漫数据
func (c *AnimeCrawler) saveAnimeData(anime *entity.Video, genres *[]entity.AnimeGenre, urls *[]entity.VideoUrl) error {
	if err := c.videoRepo.CreateVideo(context.Background(), anime); err != nil {
		return fmt.Errorf("创建视频记录失败: %v", err)
	}

	for i := range *genres {
		(*genres)[i].AnimeID = anime.VideoID
	}

	if err := c.videoRepo.CreateGenre(context.Background(), genres); err != nil {
		return fmt.Errorf("创建分类记录失败: %v", err)
	}

	for i := range *urls {
		(*urls)[i].VideoID = anime.VideoID
	}

	if err := c.videoRepo.CreateVideoUrl(context.Background(), urls); err != nil {
		return fmt.Errorf("创建剧集记录失败: %v", err)
	}

	return nil
}

// setupErrorHandlers 设置错误处理器
func (c *AnimeCrawler) setupErrorHandlers(errChan chan<- error) {
	errorHandler := func(r *colly.Response, err error) {
		select {
		case errChan <- fmt.Errorf("触发OnError: %v", err):
		default:
		}
	}

	c.collector.OnError(errorHandler)
	c.categories.OnError(errorHandler)
	c.library.OnError(errorHandler)
	c.details.OnError(errorHandler)
}

func (c *AnimeCrawler) setupCollector() {
	c.collector = c.collectorPool.Get()
	c.categories = c.collectorPool.Get()
	c.library = c.collectorPool.Get()
	c.details = c.collectorPool.Get()
}

// startCrawling 开始爬取
func (c *AnimeCrawler) startCrawling(errChan chan<- error) {
	c.setupCollector()

	// 设置处理器
	c.setupHandlers(errChan)

	c.collector.Visit(c.config.URLs.YingHuaDongMan.Main)
	// 等待所有爬虫任务完成
	c.collector.Wait()
	c.categories.Wait()
	c.library.Wait()
	c.details.Wait()
}

func (c *AnimeCrawler) resetCount() {
	atomic.StoreInt32(&c.TotalUpdateCount, 0)
}
