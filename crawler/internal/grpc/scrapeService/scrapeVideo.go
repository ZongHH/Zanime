package scrapeService

import (
	"context"
	"crawler/internal/domain/service"
	"fmt"
	"log"
	"math/rand/v2"
	"sync/atomic"
	"time"
)

type Server struct {
	UnimplementedVideoServer
	Searcher service.SearchService
}

var requestQueue = make(chan struct{}, 1)
var totalRequestCount int64 = 0

// NewServer 创建一个新的Server实例
func NewServer(searcher service.SearchService) *Server {
	return &Server{
		Searcher: searcher,
	}
}

func (s *Server) ScrapeVideoUrl(ctx context.Context, parms *VideoParms) (*VideoMsg, error) {
	requestID := atomic.AddInt64(&totalRequestCount, 1)

	select {
	case requestQueue <- struct{}{}:
		defer func() { <-requestQueue }() // 确保信号量被释放
	case <-ctx.Done():
		log.Printf("请求ID: %v, 动漫: %v, 选集: %v等待排队请求超时\n", requestID, parms.Name, parms.Episode)
		return nil, fmt.Errorf("等待排队超时: %v", ctx.Err())
	}

	randomDelay := time.Duration(500)*time.Millisecond + time.Duration(rand.IntN(500))*time.Millisecond
	// 随机延迟一段时间
	time.Sleep(randomDelay)
	fmt.Printf("随机延迟: %v\n", randomDelay)

	// 执行爬取操作
	videoUrl, err := s.Searcher.SearchAnime(
		ctx,
		parms.Name,
		parms.Release,
		parms.Area,
		parms.Episode,
	)

	// 处理错误
	if err != nil {
		log.Printf("gRpc请求错误: %v\n", err)
		return nil, err
	}

	return &VideoMsg{Url: videoUrl}, nil
}
