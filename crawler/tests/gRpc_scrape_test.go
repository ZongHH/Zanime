package tests

import (
	"context"
	"crawler/internal/bootstrap"
	"crawler/internal/grpc/scrapeService"
	"crawler/pkg/monitor"
	"sync"
	"testing"
	"time"
)

func TestGRPCScrape(t *testing.T) {
	t.Run("测试gRpc接口", func(t *testing.T) {
		monitor.Init(monitor.NewLogConfig())
		defer monitor.Close()

		container := bootstrap.BuildContainer("../configs/config.yaml")

		scrape := &scrapeService.Server{Searcher: container.Services.Search}

		name := "新干线变形机器人 改变世界"
		release := "2024"
		area := "日本"
		episode := "第02集"

		wg := sync.WaitGroup{}
		wg.Add(10)

		for i := 0; i < 10; i++ {
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				_, err := scrape.ScrapeVideoUrl(ctx, &scrapeService.VideoParms{
					Name:    name,
					Release: release,
					Area:    area,
					Episode: episode,
				})
				if err != nil {
					t.Logf("Error scraping video %d: %v", i, err)
				} else {
					t.Logf("Successfully scraped video URL %d", i)
				}
				cancel()
				time.Sleep(1 * time.Second)
			}()
		}

		wg.Wait()

	})
}
