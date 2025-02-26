package handler

import (
	"gateService/internal/domain/service"
	"gateService/internal/infrastructure/middleware/auth"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request dto.CreatePostRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.postService.CreatePost(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetPostComments(c *gin.Context) {
	var request dto.GetPostCommentsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.PageSize = 20
	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	comments, err := h.postService.GetPostCommentsByPostID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *PostHandler) GetPostCommentsByRootID(c *gin.Context) {
	var request dto.GetPostCommentsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.PageSize = 10
	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	replys, err := h.postService.GetPostCommentsByRootID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, replys)
}

func (h *PostHandler) SubmitComment(c *gin.Context) {
	var request dto.SubmitPostCommentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.postService.SubmitComment(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) SubmitPostReply(c *gin.Context) {
	var request dto.SubmitPostReplyRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.postService.SubmitPostReply(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) CommentLike(c *gin.Context) {
	var request dto.CommentLikeRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.postService.CommentLike(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetPostCategoryList(c *gin.Context) {
	var request dto.GetCategoryListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	response, err := h.postService.GetPostCategoryList(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetPostsByCategoryID(c *gin.Context) {
	var request dto.GetPostListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	request.PageSize = 10
	response, err := h.postService.GetPostsByCategoryID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetPostByPostID(c *gin.Context) {
	var request dto.GetPostByPostIDRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID

	response, err := h.postService.GetPostByPostID(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) PostLike(c *gin.Context) {
	var request dto.PostLikeRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	response, err := h.postService.PostLike(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) PostFavorite(c *gin.Context) {
	var request dto.PostFavoriteRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	response, err := h.postService.PostFavorite(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) RecentPosts(c *gin.Context) {
	var request dto.GetRecentPostsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.Error(errors.NewAppError(errors.ErrParamInvalid.Code, err.Error(), err))
		return
	}

	request.UserID = c.MustGet("UserInfo").(*auth.CustomClaims).UserInfo.UserID
	request.Page = 1
	request.PageSize = 5

	response, err := h.postService.GetRecentPosts(c.Request.Context(), &request)
	if err != nil {
		c.Error(errors.NewAppError(errors.ErrInternalError.Code, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, response)
}
