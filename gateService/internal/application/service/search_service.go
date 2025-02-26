package service

import (
	"context"
	"fmt"
	"gateService/internal/domain/repository"
	"gateService/internal/interfaces/dto"
)

type SearchServiceImpl struct {
	videoRepo repository.VideoRepository
}

func NewSearchServiceImpl(videoRepo repository.VideoRepository) *SearchServiceImpl {
	return &SearchServiceImpl{videoRepo: videoRepo}
}

func (s *SearchServiceImpl) SearchVideos(ctx context.Context, request *dto.SearchRequest) (*dto.SearchResponse, error) {
	animes, err := s.videoRepo.GetVideosByVideoName(ctx, request.Query)
	if err != nil {
		return nil, fmt.Errorf("获取搜索列表失败: %v", err)
	}
	searchAnimes := make([]*dto.SearchAnime, 0, len(animes))
	for _, anime := range animes {
		searchAnimes = append(searchAnimes, &dto.SearchAnime{
			VideoID:  anime.ID,
			Title:    anime.Name,
			CoverUrl: anime.CoverImageUrl,
		})
	}
	return &dto.SearchResponse{
		Code:   200,
		Animes: searchAnimes,
	}, nil
}

func (s *SearchServiceImpl) SearchVideosDetail(ctx context.Context, request *dto.SearchDetailRequest) (*dto.SearchDetailResponse, error) {
	animes, err := s.videoRepo.GetVideosALLEpisodesByVideoName(ctx, request.Params, request.Page)
	if err != nil {
		return nil, fmt.Errorf("获取搜索详情失败: %v", err)
	}

	videoIDs := make([]int, 0, len(animes))
	for _, anime := range animes {
		videoIDs = append(videoIDs, anime.ID)
	}
	collection, err := s.videoRepo.GetAnimeCollectionByUserAndVideoIDs(ctx, request.UserID, videoIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户收藏状态失败: %v", err)
	}

	searchDetailAnimes := make([]*dto.SearchDetailAnime, 0, len(animes))
	for _, anime := range animes {
		searchDetailAnimes = append(searchDetailAnimes, &dto.SearchDetailAnime{
			VideoID:     anime.ID,
			Title:       anime.Name,
			CoverUrl:    anime.CoverImageUrl,
			ReleaseDate: anime.ReleaseDate,
			Area:        anime.Area,
			Description: anime.Description,
			Genres:      anime.Genres,
			Episodes:    anime.Episodes,
			IsCollected: (*collection)[anime.ID],
		})
	}
	return &dto.SearchDetailResponse{
		Code:   200,
		Animes: searchDetailAnimes,
	}, nil
}
