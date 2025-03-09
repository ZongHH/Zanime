package handler

import (
	"managerService/internal/domain/service"
	"managerService/internal/interfaces/dto/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserActionLogsHandler struct {
	UserActionLogsService service.UserActionLogsService
}

func NewUserActionLogsHandler(userActionLogsService service.UserActionLogsService) *UserActionLogsHandler {
	return &UserActionLogsHandler{
		UserActionLogsService: userActionLogsService,
	}
}

func (h *UserActionLogsHandler) GetUserActionLogs(c *gin.Context) {
	request := &request.UserActionLogsRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.UserActionLogsService.GetUserActionLogs(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
