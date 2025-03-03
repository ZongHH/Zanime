package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler(videoService service.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (v *VideoHandler) GetVideoInfo(c *gin.Context) {
	request := &dto.GetVideoInfoRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := v.videoService.GetVideoInfo(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetVideoLibrary(c *gin.Context) {
	request := &dto.GetVideoLibraryRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := v.videoService.GetVideoLibrary(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetVideoFilters(c *gin.Context) {
	request := &dto.GetVideoFiltersRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := v.videoService.GetVideoFilters(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetVideoURL(c *gin.Context) {
	request := &dto.GetVideoURLRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := v.videoService.GetVideoURL(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetHomeAnimes(c *gin.Context) {
	request := &dto.GetHomeAnimesRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := v.videoService.GetHomeAnimes(c.Request.Context(), request)

	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) UpdateAnimeCollection(c *gin.Context) {
	request := &dto.UpdateAnimeCollectionRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := v.videoService.UpdateAnimeCollection(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetAnimeCollection(c *gin.Context) {
	request := &dto.GetAnimeCollectionRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID
	request.Limit = 15

	response, err := v.videoService.GetAnimeCollection(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (v *VideoHandler) GetRecommend(c *gin.Context) {
	request := &dto.GetRecommendRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	response, err := v.videoService.GetRecommend(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
