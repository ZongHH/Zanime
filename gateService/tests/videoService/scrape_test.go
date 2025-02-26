package videoService

import (
	"context"
	"gateService/internal/bootstrap"
	"gateService/internal/interfaces/dto"
	"sync"
	"testing"
)

func TestScrapeService(t *testing.T) {
	// 初始化容器
	container := bootstrap.NewContainer("../../configs/config.yaml")
	if container == nil {
		t.Fatal("初始化容器失败")
	}

	// 获取视频服务
	videoService := container.Services.VideoService
	if videoService == nil {
		t.Fatal("获取视频服务失败")
	}

	// 构造请求参数
	request := &dto.GetVideoURLRequest{
		UserID:  1,
		VideoID: 4858,
		Episode: "第05集",
	}

	t.Run("并发获取视频URL", func(t *testing.T) {
		var wg sync.WaitGroup
		concurrency := 10
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()

				// 设置上下文
				ctx := context.Background()
				response, err := videoService.GetVideoURL(ctx, request)
				if err != nil {
					t.Errorf("获取视频URL失败: %v", err)
					return
				}

				if response == nil {
					t.Error("获取视频URL响应为空")
					return
				}

				t.Logf("获取视频URL成功: %+v", response)
			}()
		}

		wg.Wait()
	})
}
