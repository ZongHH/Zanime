package search

import (
	"context"
	"crawler/internal/domain/repository"
	"crawler/internal/infrastructure/collector"
	"crawler/internal/infrastructure/config"
	"crawler/pkg/monitor"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// AnimeMatcher 动漫匹配器
type AnimeMatcher struct {
	Name    string
	Release string
	Area    string
	Episode string
}

// AnimeSearcher 动漫搜索器
type AnimeSearcher struct {
	searchURL       string
	collectorPool   *collector.CollectorPool
	scraper         *VideoScraper
	videoRepository repository.VideoRepository
}

// NewAnimeSearcher 创建新的动漫搜索器
func NewAnimeSearcher(cfg *config.Config, collectorPool *collector.CollectorPool, scraper *VideoScraper, videoRepository repository.VideoRepository) *AnimeSearcher {
	return &AnimeSearcher{
		searchURL:       cfg.URLs.YingHuaDongMan.Search,
		collectorPool:   collectorPool,
		scraper:         scraper,
		videoRepository: videoRepository,
	}
}

// SearchAnime 搜索动漫视频
func (s *AnimeSearcher) SearchAnime(ctx context.Context, name, release, area, episode string) (string, error) {
	matcher := AnimeMatcher{
		Name:    name,
		Release: release,
		Area:    area,
		Episode: episode,
	}

	videoID, err := s.videoRepository.GetVideoIDByVideoName(ctx, matcher.Name)
	if err != nil {
		return "", fmt.Errorf("数据库获取动漫ID失败: %v", err)
	}

	videoURL, err := s.videoRepository.GetVideoURLByVideoIDANDEpisode(ctx, videoID, matcher.Episode)
	if err != nil {
		return "", fmt.Errorf("数据库获取动漫连接失败: %v", err)
	}

	resultChan := make(chan string, 1)

	if videoURL == "" {
		monitor.Info("开始搜索:name %v release: %s area: %s episode: %s", matcher.Name, matcher.Release, matcher.Area, matcher.Episode)

		c := s.collectorPool.Get()
		detailCollector := c.Clone()

		// 设置搜索结果处理器
		s.setupSearchHandler(c, detailCollector)

		// 设置详情页处理器
		s.setupDetailHandler(detailCollector, matcher, resultChan)

		// 开始搜索
		go s.startSearch(c, matcher.Name)
	} else {
		resultChan <- videoURL
		log.Printf("数据库获取到Name: %s, Episode: %s地址: %s\n", matcher.Name, matcher.Episode, videoURL)
	}

	return s.waitForResult(ctx, resultChan, name, episode)
}

// setupSearchHandler 设置搜索结果处理器
func (s *AnimeSearcher) setupSearchHandler(c *colly.Collector, detailCollector *colly.Collector) {
	c.OnHTML("a.module-card-item-poster", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		absoluteURL := h.Request.AbsoluteURL(href)
		detailCollector.Visit(absoluteURL)
	})
}

// setupDetailHandler 设置详情页处理器
func (s *AnimeSearcher) setupDetailHandler(detailCollector *colly.Collector, matcher AnimeMatcher, resultChan chan<- string) {
	detailCollector.OnHTML("div.page.view", func(h *colly.HTMLElement) {
		if details := s.extractAnimeDetails(h); s.isMatchingAnime(details, matcher) {
			s.processEpisodes(h, matcher.Episode, resultChan)
		}
	})
}

// extractAnimeDetails 提取动漫详情
func (s *AnimeSearcher) extractAnimeDetails(h *colly.HTMLElement) map[string]string {
	details := make(map[string]string)
	details["name"] = h.ChildText("div.module-info-heading h1")

	h.ForEach("div.module-info-tag-link", func(i int, h *colly.HTMLElement) {
		switch i {
		case 0:
			details["release"] = h.Text
		case 1:
			details["area"] = h.Text
		}
	})

	return details
}

// isMatchingAnime 检查是否匹配目标动漫
func (s *AnimeSearcher) isMatchingAnime(details map[string]string, matcher AnimeMatcher) bool {
	return details["name"] == matcher.Name &&
		details["release"] == matcher.Release &&
		details["area"] == matcher.Area
}

// processEpisodes 处理剧集
func (s *AnimeSearcher) processEpisodes(h *colly.HTMLElement, targetEpisode string, resultChan chan<- string) {
	var found bool
	h.ForEach("a.module-play-list-link", func(i int, h *colly.HTMLElement) {
		if found || h.Text != targetEpisode {
			return
		}

		href := h.Attr("href")
		absoluteURL := h.Request.AbsoluteURL(href)
		resultChan <- absoluteURL
		found = true
	})
}

// startSearch 开始搜索
func (s *AnimeSearcher) startSearch(c *colly.Collector, name string) {
	searchName := strings.ReplaceAll(name, " ", "%20")
	c.Visit(s.searchURL + searchName)
}

// waitForResult 等待搜索结果
func (s *AnimeSearcher) waitForResult(ctx context.Context, resultChan <-chan string, name, episode string) (string, error) {
	select {
	case url := <-resultChan:
		videoURL, err := s.scraper.ScrapeVideo(url) // 使用父超时ctx如果提前释放会导致chrome实例不能正确释放
		if err != nil {
			return "", fmt.Errorf("动态爬取视频链接失败: %v", err)
		}
		monitor.Info("%v:%v视频链接: %v", name, episode, videoURL)
		return videoURL, nil
	case <-ctx.Done():
		return "", fmt.Errorf("动态爬取视频链接超时")
	}
}
