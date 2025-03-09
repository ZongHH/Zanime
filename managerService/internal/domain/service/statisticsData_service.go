package service

import (
	"context"
	"managerService/internal/interfaces/dto/request"
	"managerService/internal/interfaces/dto/response"
)

type StatisticsDataService interface {
	// GetStatisticsData 获取统计数据
	// ctx 上下文
	// request 请求参数
	// 返回统计数据和可能的错误
	GetStatisticsData(ctx context.Context, request *request.StatisticsDataRequest) (*response.StatisticsDataRespnse, error)

	// GetNewAnime 获取最新上线的动漫
	// ctx 上下文
	// request 请求参数
	// 返回最新上线的动漫和可能的错误
	GetNewAnime(ctx context.Context, request *request.NewAnimeRequest) (*response.NewAnimeResponse, error)
}
