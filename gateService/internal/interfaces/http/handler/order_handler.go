package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	request := &dto.CreateOrderRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	response, err := h.orderService.CreateOrder(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	request := &dto.GetOrderRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	// 暂时不从URL参数获取pageSize
	request.PageSize = 10

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	response, err := h.orderService.GetOrdersByUserID(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) CallbackPay(c *gin.Context) {
	request := &dto.PayRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.orderService.CallbackPay(c.Request.Context(), request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
