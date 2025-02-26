package scrapeService

import (
	"context"
	"crawler/internal/domain/service"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	successID  int64 = 0 // 最新完成的请求ID
	totalCount int64 = 0 // 总请求次数
)

type Server struct {
	UnimplementedVideoServer
	Searcher service.SearchService
}

func (s *Server) ScrapeVideoUrl(ctx context.Context, parms *VideoParms) (*VideoMsg, error) {
	requestID := atomic.AddInt64(&totalCount, 1)

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("等待排队爬取超时: %v", ctx.Err())
		default:
		}
		if requestID == atomic.LoadInt64(&successID)+1 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	videoUrl, err := s.Searcher.SearchAnime(ctx,
		parms.Name,
		parms.Release,
		parms.Area,
		parms.Episode,
	)

	// successID更新为requestID
	atomic.StoreInt64(&successID, requestID)

	if err != nil {
		return nil, err
	}

	var videoMsg = &VideoMsg{
		Url: videoUrl,
	}

	return videoMsg, nil
}
