// Package service 实现了视频服务的具体业务逻辑
package service

import (
	"context"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/grpc/client/recommend"
	"gateService/internal/grpc/client/scrapeClient"
	"gateService/internal/infrastructure/middleware/lock"
	"gateService/internal/interfaces/dto"
	"math/rand/v2"
	"sync"

	"github.com/redis/go-redis/v9"
)

// VideoServiceImpl 实现了VideoService接口,提供视频相关的业务功能
type VideoServiceImpl struct {
	rdb                *redis.Client                 // Redis客户端,用于创建分布式锁
	scrapeClient       *scrapeClient.GRPCClientPool  // 视频爬虫客户端池
	recommendClient    *recommend.GRPCClientPool     // 推荐客户端池
	videoRepositoty    repository.VideoRepository    // 视频仓储接口
	progressRepository repository.ProgressRepository // 观看进度仓储接口

}

// NewVideoServiceImpl 创建VideoServiceImpl的新实例
// 参数:
//   - rdb: Redis客户端
//   - scrape: 视频爬虫客户端池
//   - videoRepositoty: 视频仓储实现
//   - progressRepository: 观看进度仓储实现
//
// 返回:
//   - *VideoServiceImpl: 服务实例
func NewVideoServiceImpl(rdb *redis.Client, scrapeClient *scrapeClient.GRPCClientPool, recommendClient *recommend.GRPCClientPool, videoRepositoty repository.VideoRepository, progressRepository repository.ProgressRepository) *VideoServiceImpl {
	return &VideoServiceImpl{
		rdb:                rdb,
		scrapeClient:       scrapeClient,
		recommendClient:    recommendClient,
		videoRepositoty:    videoRepositoty,
		progressRepository: progressRepository,
	}
}

// GetVideoInfo 获取视频的详细信息
// 参数:
//   - ctx: 上下文信息
//   - request: 包含视频ID的请求参数
//
// 返回:
//   - *dto.GetVideoInfoResponse: 视频详细信息响应
//   - error: 可能的错误信息
func (v *VideoServiceImpl) GetVideoInfo(ctx context.Context, request *dto.GetVideoInfoRequest) (*dto.GetVideoInfoResponse, error) {
	response, err := v.videoRepositoty.GetVideoInfoWithEposidesByVideoID(ctx, request.VideoID)
	if err != nil {
		return &dto.GetVideoInfoResponse{
			Code:      500,
			VideoInfo: nil,
		}, fmt.Errorf("获取视频详细信息失败: %v", err)
	}

	return &dto.GetVideoInfoResponse{
		Code:      200,
		VideoInfo: response,
	}, nil
}

// GetVideoLibrary 根据筛选条件获取视频库列表
// 参数:
//   - ctx: 上下文信息
//   - request: 包含地区、年份、类型、字母等筛选条件的请求参数
//
// 返回:
//   - *dto.GetVideoLibraryRespnse: 视频库列表响应
//   - error: 可能的错误信息
func (v *VideoServiceImpl) GetVideoLibrary(ctx context.Context, request *dto.GetVideoLibraryRequest) (*dto.GetVideoLibraryRespnse, error) {
	video := &entity.Video{
		Area:        request.Region,
		ReleaseDate: request.Year,
		Genres:      request.Type,
		Initial:     request.Letter,
	}

	reponse, total, err := v.videoRepositoty.GetVideosByFilters(ctx, video, request.Page, request.PageSize)
	if err != nil {
		return &dto.GetVideoLibraryRespnse{
			Code:   500,
			Total:  0,
			Videos: nil,
		}, fmt.Errorf("通过筛选获取动漫资源失败: %v", err)
	}

	return &dto.GetVideoLibraryRespnse{
		Code:   200,
		Total:  total,
		Videos: reponse,
	}, nil
}

// GetVideoFilters 获取视频的筛选条件选项
// 参数:
//   - ctx: 上下文信息
//   - request: 筛选条件请求参数
//
// 返回:
//   - *dto.GetVideoFiltersResponse: 视频筛选条件响应
//   - error: 可能的错误信息
func (v *VideoServiceImpl) GetVideoFilters(ctx context.Context, request *dto.GetVideoFiltersRequest) (*dto.GetVideoFiltersResponse, error) {
	response, err := v.videoRepositoty.GetVideoFilters(ctx)
	if err != nil {
		return &dto.GetVideoFiltersResponse{
			Code:         500,
			VideoFilters: nil,
		}, fmt.Errorf("获取动漫筛选选项失败: %v", err)
	}

	return &dto.GetVideoFiltersResponse{
		Code:         200,
		VideoFilters: response,
	}, nil
}

// GetVideoURL 获取视频播放地址
// 参数:
//   - ctx: 上下文信息,用于控制请求的生命周期
//   - request: 包含用户ID、视频ID和集数的请求参数
//
// 返回:
//   - *dto.GetVideoURLResponse: 视频URL响应对象
//   - error: 可能的错误信息
func (v *VideoServiceImpl) GetVideoURL(ctx context.Context, request *dto.GetVideoURLRequest) (*dto.GetVideoURLResponse, error) {
	// 获取用户观看进度
	progress, err := v.progressRepository.GetUserWatchProgress(ctx, request.UserID, request.VideoID)
	if err != nil {
		return v.Response(500, ""), fmt.Errorf("获取用户历史记录失败: %v", err)
	}

	// 如果请求中没有指定集数,则使用用户上次观看的集数
	if request.Episode == "" || request.Episode == "null" {
		request.Episode = progress.Episode
	}

	// 构造Redis缓存key
	URLKey := fmt.Sprintf("VideoURL:Video%d:Episode%s", request.VideoID, request.Episode)
	// 尝试从缓存获取视频URL
	URL, err := v.videoRepositoty.GetVideoURL(ctx, URLKey)
	if err != nil && err != redis.Nil {
		return v.Response(500, ""), fmt.Errorf("获取缓存视频链接失败: %v", err)
	}
	// 如果缓存中存在,直接返回
	if URL != "" {
		return v.Response(200, URL), nil
	}

	// 构造分布式锁key
	key := fmt.Sprintf("scrape_lock:%d:%s", request.VideoID, request.Episode)
	redisLock := lock.NewRedisLock(v.rdb, key, nil)
	// 获取分布式锁,防止并发爬取
	if err := redisLock.WaitLock(ctx); err != nil {
		return v.Response(500, ""), fmt.Errorf("获取分布式锁%s失败: %v", key, err)
	}
	defer redisLock.Unlock(ctx)

	// 双重检查,再次尝试从缓存获取
	URL, err = v.videoRepositoty.GetVideoURL(ctx, URLKey)
	if err != nil && err != redis.Nil {
		return v.Response(500, ""), fmt.Errorf("获取缓存视频链接失败: %v", err)
	}
	if URL != "" {
		return v.Response(200, URL), nil
	}

	// 爬取视频URL
	VideoMsg, err := v.scrapeClient.ScrapeVideoUrl(ctx, progress.VideoName, progress.Release, progress.Area, request.Episode)
	if err != nil {
		return v.Response(500, ""), fmt.Errorf("爬取视频链接失败: %v", err)
	}

	// 异步缓存视频URL
	go v.videoRepositoty.CacheVideoURL(context.Background(), URLKey, VideoMsg.Url)

	return v.Response(200, VideoMsg.Url), nil
}

// GetHomeAnimes 获取首页动漫列表
// 参数:
//   - ctx: 上下文信息
//   - request: 包含用户ID的请求参数
//
// 返回:
//   - *dto.GetHomeAnimesResponse: 首页动漫列表响应
//   - error: 可能的错误信息
func (v *VideoServiceImpl) GetHomeAnimes(ctx context.Context, request *dto.GetHomeAnimesRequest) (*dto.GetHomeAnimesResponse, error) {
	var (
		response = &dto.GetHomeAnimesResponse{}
		wg       sync.WaitGroup
		errChan  = make(chan error, 5)
	)

	// 设置首页推荐动漫
	response.HomeAnime = &entity.Video{
		ID:          21542,
		Name:        "英雄联盟：双城之战 第二季",
		VideoUrl:    "/src/static/videos/shuangcheng.mp4",
		Description: "黑暗梦魇，有增无减。务必认清，敌人的真面目……",
	}

	// 并发获取各类动漫列表
	wg.Add(5)

	// 获取日本动漫
	go func() {
		defer wg.Done()
		animes, err := v.videoRepositoty.GetAnimesByGenre(ctx, "日本动漫", 1, 10)
		if err != nil {
			errChan <- fmt.Errorf("获取日本动漫列表失败: %v", err)
			return
		}
		response.JapanAnime = animes
	}()

	// 获取国产动漫
	go func() {
		defer wg.Done()
		animes, err := v.videoRepositoty.GetAnimesByGenre(ctx, "国产动漫", 1, 10)
		if err != nil {
			errChan <- fmt.Errorf("获取国产动漫列表失败: %v", err)
			return
		}
		response.ChinaAnime = animes
	}()

	// 获取欧美动漫
	go func() {
		defer wg.Done()
		animes, err := v.videoRepositoty.GetAnimesByGenre(ctx, "欧美动漫", 1, 10)
		if err != nil {
			errChan <- fmt.Errorf("获取欧美动漫列表失败: %v", err)
			return
		}
		response.WesternAnime = animes
	}()

	// 获取热门动漫类型
	go func() {
		defer wg.Done()
		genres, err := v.videoRepositoty.GetTopAnimeGenres(ctx)
		if err != nil {
			errChan <- fmt.Errorf("获取热门动漫类型失败: %v", err)
			return
		}
		response.AnimeGenres = genres
	}()

	// 获取推荐列表
	go func() {
		defer wg.Done()
		tops, err := v.recommendClient.GetListRecommend(ctx, request.UserID)
		if err != nil {
			errChan <- fmt.Errorf("获取推荐列表失败: %v", err)
			return
		}

		// 构造推荐动漫列表
		for _, top := range tops.Recommendations {
			response.TopAnime = append(response.TopAnime, &entity.Video{
				ID:            int(top.VideoId),
				Name:          top.VideoName,
				CoverImageUrl: top.CoverImageUrl,
			})
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		if err != nil {
			return &dto.GetHomeAnimesResponse{
				Code: 500,
			}, err
		}
	}

	response.Code = 200
	return response, nil
}

// Response 生成视频URL响应
// 参数:
//   - code: 状态码
//   - URL: 视频播放地址
//
// 返回:
//   - *dto.GetVideoURLResponse: 视频URL响应对象
func (v *VideoServiceImpl) Response(code int, URL string) *dto.GetVideoURLResponse {
	return &dto.GetVideoURLResponse{
		Code: code,
		VideoFiles: []map[string]string{
			{
				"file":  URL,
				"label": "720p",
			},
		},
		PosterPath: "",
	}
}

// UpdateAnimeCollection 更新动漫收藏状态
func (v *VideoServiceImpl) UpdateAnimeCollection(ctx context.Context, request *dto.UpdateAnimeCollectionRequest) (*dto.UpdateAnimeCollectionResponse, error) {
	// 获取当前收藏状态
	collection := &entity.UserAnimeCollection{
		UserID:  request.UserID,
		VideoID: request.VideoID,
	}

	if request.Status {
		// 更新收藏状态
		err := v.videoRepositoty.AddAnimeCollection(ctx, collection)
		if err != nil {
			return &dto.UpdateAnimeCollectionResponse{Code: 500}, fmt.Errorf("更新收藏状态失败: %v", err)
		}
	} else {
		// 删除收藏状态
		err := v.videoRepositoty.DeleteAnimeCollection(ctx, request.UserID, request.VideoID)
		if err != nil {
			return &dto.UpdateAnimeCollectionResponse{Code: 500}, fmt.Errorf("删除收藏状态失败: %v", err)
		}
	}

	return &dto.UpdateAnimeCollectionResponse{Code: 200}, nil
}

// GetAnimeCollection 获取用户的动漫收藏列表
func (v *VideoServiceImpl) GetAnimeCollection(ctx context.Context, request *dto.GetAnimeCollectionRequest) (*dto.GetAnimeCollectionResponse, error) {
	// 获取收藏列表
	videos, total, err := v.videoRepositoty.GetAnimeCollectionByUser(ctx, request.UserID, request.Page, request.Limit)
	if err != nil {
		return &dto.GetAnimeCollectionResponse{Code: 500}, fmt.Errorf("获取收藏列表失败: %v", err)
	}

	animeCollection := make([]*dto.AnimeCollection, 0, len(videos))
	for _, video := range videos {
		animeCollection = append(animeCollection, &dto.AnimeCollection{
			VideoID:       video.ID,
			Title:         video.Name,
			CoverImageUrl: video.CoverImageUrl,
			CollectedAt:   video.CollectedAt,
		})
	}

	return &dto.GetAnimeCollectionResponse{
		Code:            200,
		Total:           total,
		AnimeCollection: animeCollection,
	}, nil
}

// GetRecommend 获取相关动漫推荐
func (v *VideoServiceImpl) GetRecommend(ctx context.Context, request *dto.GetRecommendRequest) (*dto.GetRecommendResponse, error) {
	recommendations, err := v.videoRepositoty.GetAnimeGenres(ctx, request.VideoID)
	if err != nil {
		return &dto.GetRecommendResponse{Code: 500}, fmt.Errorf("获取当前动漫类型失败: %v", err)
	}

	recommendedAnimes := make([]*dto.RecommendedAnime, 0)
	for _, genre := range recommendations {
		animes, err := v.videoRepositoty.GetAnimesByGenre(ctx, genre, 1, 100)
		if err != nil {
			return &dto.GetRecommendResponse{Code: 500}, fmt.Errorf("获取推荐动漫失败: %v", err)
		}

		visitedAnimes := make(map[int]bool)
		for i := 0; i < 3; i++ {
			pos := rand.IntN(len(animes))
			if visitedAnimes[animes[pos].ID] {
				continue
			}
			visitedAnimes[animes[pos].ID] = true
			recommendedAnimes = append(recommendedAnimes, &dto.RecommendedAnime{
				ID:       animes[pos].ID,
				Title:    animes[pos].Name,
				CoverUrl: animes[pos].CoverImageUrl,
				Rating:   fmt.Sprintf("%0.1f", 9+rand.Float64()),
			})
		}
	}

	return &dto.GetRecommendResponse{
		Code:            200,
		Recommendations: recommendedAnimes,
	}, nil
}
