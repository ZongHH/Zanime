package service

import (
	"context"
	"fmt"
	"managerService/internal/domain/repository"
	"managerService/internal/interfaces/dto/request"
	"managerService/internal/interfaces/dto/response"
	"sync"
)

type StatisticsServiceImpl struct {
	statisticsRepo repository.StatisticsRepository
}

func NewStatisticsServiceImpl(statisticsRepo repository.StatisticsRepository) *StatisticsServiceImpl {
	return &StatisticsServiceImpl{
		statisticsRepo: statisticsRepo,
	}
}

func (s *StatisticsServiceImpl) GetStatisticsData(ctx context.Context, request *request.StatisticsDataRequest) (*response.StatisticsDataRespnse, error) {
	errorChan := make(chan error, 4)
	statisticsItems := make([]*response.StatisticsItem, 4)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		userCount, err := s.statisticsRepo.GetUserCount(ctx)
		if err != nil {
			errorChan <- fmt.Errorf("获取用户总数失败: %v", err)
			return
		}
		statisticsItems[0] = &response.StatisticsItem{
			Label:      "总用户数",
			Value:      userCount,
			Icon:       "User",
			Color:      "#409EFF",
			Trend:      "up",
			Percentage: 15,
		}
	}()

	go func() {
		defer wg.Done()
		animeCount, err := s.statisticsRepo.GetAnimeCount(ctx)
		if err != nil {
			errorChan <- fmt.Errorf("获取动漫总数失败: %v", err)
			return
		}
		statisticsItems[1] = &response.StatisticsItem{
			Label:      "动漫数量",
			Value:      animeCount,
			Icon:       "VideoPlay",
			Color:      "#67C23A",
			Trend:      "up",
			Percentage: 8,
		}
	}()

	go func() {
		defer wg.Done()
		todayPlayCount, err := s.statisticsRepo.GetTodayPlayCount(ctx)
		if err != nil {
			errorChan <- fmt.Errorf("获取今日播放总数失败: %v", err)
			return
		}
		statisticsItems[2] = &response.StatisticsItem{
			Label:      "今日播放",
			Value:      todayPlayCount,
			Icon:       "VideoCamera",
			Color:      "#E6A23C",
			Trend:      "down",
			Percentage: 3,
		}
	}()

	go func() {
		defer wg.Done()
		activeUserCount, err := s.statisticsRepo.GetActiveUserCount(ctx)
		if err != nil {
			errorChan <- fmt.Errorf("获取活跃用户总数失败: %v", err)
			return
		}
		statisticsItems[3] = &response.StatisticsItem{
			Label:      "活跃用户",
			Value:      activeUserCount,
			Icon:       "Star",
			Color:      "#F56C6C",
			Trend:      "up",
			Percentage: 12,
		}
	}()

	closeChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(closeChan)
	}()

	select {
	case <-closeChan:
		return &response.StatisticsDataRespnse{
			Code:            200,
			StatisticsItems: statisticsItems,
		}, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("获取统计数据超时")
	case err := <-errorChan:
		return nil, err
	}
}

func (s *StatisticsServiceImpl) GetNewAnime(ctx context.Context, request *request.NewAnimeRequest) (*response.NewAnimeResponse, error) {
	animes, err := s.statisticsRepo.GetNewAnime(ctx, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取最新上线的动漫失败: %v", err)
	}
	animeItems := make([]*response.AnimeItem, len(animes))
	for i, anime := range animes {
		animeItems[i] = &response.AnimeItem{
			Title:      anime.VideoName,
			Image:      anime.CoverImageURL,
			UpdateTime: anime.CreatedAt,
		}
	}
	return &response.NewAnimeResponse{
		Code:   200,
		Animes: animeItems,
	}, nil
}
