package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressService service.ProgressService
}

func NewProgressHandler(progressService service.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

func (h *ProgressHandler) WatchHistory(c *gin.Context) {
	request := &dto.WatchHistoryRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID
	request.PageSize = 15

	if request.Page <= 0 {
		request.Page = 1
	}

	response, err := h.progressService.GetWatchHistory(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *ProgressHandler) LoadProgress(c *gin.Context) {
	request := &dto.LoadProgressRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := h.progressService.LoadProgress(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProgressHandler) SaveProgress(c *gin.Context) {
	request := &dto.SaveProgressRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.progressService.SaveProgress(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
