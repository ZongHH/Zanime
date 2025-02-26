package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/mq/nsqpool"
)

type ProgressServiceImpl struct {
	producerPool       *nsqpool.ProducerPool
	progressRepository repository.ProgressRepository
}

func NewProgressServiceImpl(producerPool *nsqpool.ProducerPool, progressRepository repository.ProgressRepository) *ProgressServiceImpl {
	return &ProgressServiceImpl{
		progressRepository: progressRepository,
		producerPool:       producerPool,
	}
}

func (p *ProgressServiceImpl) GetWatchHistory(ctx context.Context, request *dto.WatchHistoryRequest) (*dto.WatchHistoryResponse, error) {
	progress, err := p.progressRepository.GetProgress(ctx, request.UserID, request.Page, request.PageSize)
	if err != nil {
		return &dto.WatchHistoryResponse{
			Code:     500,
			Progress: nil,
		}, fmt.Errorf("获取用户观看历史记录失败: %v", err)
	}

	return &dto.WatchHistoryResponse{
		Code:     200,
		Progress: progress,
	}, nil
}

func (p *ProgressServiceImpl) SaveProgress(ctx context.Context, request *dto.SaveProgressRequest) (*dto.SaveProgressResponse, error) {
	progress := &entity.Progress{
		VideoID:  request.VideoID,
		UserID:   request.UserID,
		Episode:  request.Episode,
		Progress: request.Progress,
	}
	err := p.progressRepository.UpdateProgress(ctx, progress)
	if err != nil {
		return &dto.SaveProgressResponse{
			Code: 500,
		}, fmt.Errorf("更新用户观看进度失败: %v", err)
	}

	err = p.PublishTONsq(ctx, request)
	if err != nil {
		return &dto.SaveProgressResponse{
			Code: 500,
		}, fmt.Errorf("观看进度发送到推荐系统失败: %v", err)
	}

	return &dto.SaveProgressResponse{
		Code: 200,
	}, nil
}

func (p *ProgressServiceImpl) LoadProgress(ctx context.Context, request *dto.LoadProgressRequest) (*dto.LoadProgressResponse, error) {
	progress, err := p.progressRepository.GetUserWatchProgress(ctx, request.UserID, request.VideoID)
	if err != nil {
		return &dto.LoadProgressResponse{
			Code:     500,
			Progress: nil,
		}, fmt.Errorf("获取该视频用户历史记录失败: %v", err)
	}

	progress.UserID = request.UserID
	progress.VideoID = request.VideoID

	if !(request.Episode == "" || request.Episode == "null") && request.Episode != progress.Episode {
		progress.Episode = request.Episode
		progress.Progress = 0
	}

	return &dto.LoadProgressResponse{
		Code:     200,
		Progress: progress,
	}, nil
}

func (p *ProgressServiceImpl) PublishTONsq(ctx context.Context, progress *dto.SaveProgressRequest) error {
	progress.MsgType = "user_behavior"
	progressJson, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	err = p.producerPool.Publish(ctx, "user_video_behavior", progressJson)
	if err != nil {
		return err
	}

	return nil
}
