package handler

import (
	"managerService/internal/domain/service"
	"managerService/internal/interfaces/dto/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsService service.StatisticsDataService
}

func NewStatisticsHandler(statisticsService service.StatisticsDataService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (s *StatisticsHandler) GetStatisticsData(ctx *gin.Context) {
	request := &request.StatisticsDataRequest{}
	if err := ctx.ShouldBind(request); err != nil {
		ctx.Error(err)
		return
	}

	response, err := s.statisticsService.GetStatisticsData(ctx.Request.Context(), request)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
