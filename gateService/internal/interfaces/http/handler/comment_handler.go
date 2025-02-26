package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) SubmitComment(c *gin.Context) {
	var request dto.SubmitCommentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.commentService.SubmitComment(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	var request dto.GetCommentsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.PageSize = 20

	response, err := h.commentService.GetCommentsByVideoID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) SubmitReply(c *gin.Context) {
	var request dto.SubmitReplyRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.commentService.SubmitReply(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) GetReply(c *gin.Context) {
	var request dto.GetReplyRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.PageSize = 10

	response, err := h.commentService.GetCommentsByRootID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
