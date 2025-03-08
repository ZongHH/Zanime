package service

import (
	"context"
	"managerService/internal/interfaces/dto/request"
	"managerService/internal/interfaces/dto/response"
)

type StatisticsDataService interface {
	GetStatisticsData(context.Context, *request.StatisticsDataRequest) (*response.StatisticsDataRespnse, error)
}
