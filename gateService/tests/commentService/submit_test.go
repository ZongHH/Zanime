package commentService

import (
	"context"
	"gateService/internal/bootstrap"
	"gateService/internal/interfaces/dto"
	"sync"
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	// submit_test.go:42: 测试运行时间: 10.9344471s 同步持久化
	// submit_test.go:42: 测试运行时间: 4.3067969s  异步持久化

	startTime := time.Now() // 记录开始时间
	container := bootstrap.NewContainer("../../configs/config.yaml")

	request := &dto.SubmitPostCommentRequest{
		PostID:    4,
		CommentID: 1,
		Content:   "test",
		UserID:    1,
	}

	submit := container.Services.PostService

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_, err := submit.SubmitComment(context.Background(), request)
				if err != nil {
					t.Errorf("提交评论失败: %v", err)
				}
				// t.Logf("提交评论成功: %v", response)
			}
		}()
	}

	wg.Wait()
	elapsedTime := time.Since(startTime) // 计算运行时间
	t.Logf("测试运行时间: %v", elapsedTime)
}
