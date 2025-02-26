package tests

import (
	"context"
	"crawler/internal/application/search"
	"sync"
	"testing"
	"time"
)

func TestScrapeVideo(t *testing.T) {

	t.Run("测试scrape功能", func(t *testing.T) {
		// 使用默认配置创建scraper
		scraper := search.NewVideoScraper(nil)

		videoURL, err := scraper.ScrapeVideo(context.Background(), "https://www.yinghuadm.cn/play_9091-1-1.html")
		if err != nil {
			t.Logf("Error scraping video: %v", err)
			return
		}
		t.Logf("Successfully scraped video URL: %s", videoURL)
	})

	// 即时不同UserAgent，也会被封禁
	t.Run("测试scrape功能并发能力", func(t *testing.T) {
		// 使用默认配置创建scraper
		scraper := search.NewVideoScraper(nil)

		urls := []string{
			"https://www.yinghuadm.cn/play_468-1-1.html",
			"https://www.yinghuadm.cn/play_468-1-2.html",
			"https://www.yinghuadm.cn/play_468-1-3.html",
			"https://www.yinghuadm.cn/play_468-1-4.html",
			"https://www.yinghuadm.cn/play_468-1-5.html",
			"https://www.yinghuadm.cn/play_468-1-6.html",
			"https://www.yinghuadm.cn/play_468-1-7.html",
			"https://www.yinghuadm.cn/play_468-1-8.html",
			"https://www.yinghuadm.cn/play_468-1-9.html",
			"https://www.yinghuadm.cn/play_468-1-10.html",
			"https://www.yinghuadm.cn/play_468-1-11.html",
			"https://www.yinghuadm.cn/play_468-1-12.html",
		}

		wg := sync.WaitGroup{}
		wg.Add(len(urls))

		for i, url := range urls {
			go func(url string) {
				defer wg.Done()
				videoURL, err := scraper.ScrapeVideo(context.Background(), url)
				if err != nil {
					t.Logf("Error scraping video %d: %v", i, err)
				} else {
					t.Logf("Successfully scraped video URL %d: %s", i, videoURL)
				}
			}(url)
		}

		wg.Wait()
	})

	t.Run("测试scrape功能顺序执行", func(t *testing.T) {
		// 使用默认配置创建scraper
		scraper := search.NewVideoScraper(nil)

		urls := []string{
			"https://www.yinghuadm.cn/play_468-1-1.html",
			"https://www.yinghuadm.cn/play_468-1-2.html",
			"https://www.yinghuadm.cn/play_468-1-3.html",
			"https://www.yinghuadm.cn/play_468-1-4.html",
			"https://www.yinghuadm.cn/play_468-1-5.html",
			"https://www.yinghuadm.cn/play_468-1-6.html",
			"https://www.yinghuadm.cn/play_468-1-7.html",
			"https://www.yinghuadm.cn/play_468-1-8.html",
			"https://www.yinghuadm.cn/play_468-1-9.html",
			"https://www.yinghuadm.cn/play_468-1-10.html",
			"https://www.yinghuadm.cn/play_468-1-11.html",
			"https://www.yinghuadm.cn/play_468-1-12.html",
		}

		for i, url := range urls {
			// 为每个请求创建子上下文
			ctx, cancel := context.WithCancel(context.Background())
			videoURL, err := scraper.ScrapeVideo(ctx, url)
			if err != nil {
				t.Logf("Error scraping video %d: %v", i, err)
			} else {
				t.Logf("Successfully scraped video URL %d: %s", i, videoURL)
			}
			cancel() // 立即取消子上下文

			// 强制等待一小段时间确保资源释放
			time.Sleep(2 * time.Second)
		}
	})
}

func TestScraperWithCustomConfig(t *testing.T) {
	ctx := context.Background()

	// 自定义配置
	config := search.ScraperConfig{
		Timeout:       10 * time.Second,
		RetryCount:    2,
		RetryInterval: 1 * time.Second,
		RandomDelay:   2 * time.Second,
	}

	scraper := search.NewVideoScraper(&config)

	url := "https://example.com/video/test"
	videoURL, err := scraper.ScrapeVideo(ctx, url)

	if err != nil {
		t.Logf("Error with custom config: %v", err)
		return
	}
	t.Logf("Successfully scraped with custom config: %s", videoURL)
}
