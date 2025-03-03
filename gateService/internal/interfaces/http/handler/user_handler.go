package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.userService.Register(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.userService.Login(c, &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) VerifyUser(c *gin.Context) {
	var request dto.VerifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.userService.VerifyUser(c, &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Logout(c *gin.Context) {
	err := h.userService.Logout(c)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	var request dto.UserInfoRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.userService.GetUserInfo(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	var request dto.GetUserProfileRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.userService.GetUserProfile(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var request dto.UpdateUserProfileRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.userService.UpdateUserProfile(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	var request dto.UserStatsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.userService.GetUserStats(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	var request dto.UploadAvatarRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.userService.UploadAvatar(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserNotifications(c *gin.Context) {
	var request dto.UserNotificationRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	request.PageSize = 10

	response, err := h.userService.GetUserNotifications(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetTestAccount(c *gin.Context) {
	var request dto.TestAccountRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserIP = c.ClientIP()

	response, err := h.userService.GetTestAccount(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
