package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	websocketService service.WebSocketService
}

func NewWebSocketHandler(websocketService service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		websocketService: websocketService,
	}
}

func (w *WebSocketHandler) EstablishConnection(c *gin.Context) {
	request := &dto.EstablishWebSocketRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	userInfo := c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo
	request.UserID = userInfo.UserID

	_, err := w.websocketService.EstablishConnection(c, request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	// c.JSON(http.StatusOK, response)
}
