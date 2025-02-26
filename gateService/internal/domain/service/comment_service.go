package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// CommentService 定义评论相关的业务逻辑接口
// 提供评论的提交、获取等功能
type CommentService interface {
	// SubmitComment 提交评论
	// 参数:
	// - ctx: 上下文
	// - request: 评论提交请求,包含评论内容、视频ID等信息
	// 返回:
	// - *dto.SubmitCommentResponse: 提交评论响应
	// - error: 提交过程中的错误信息
	SubmitComment(ctx context.Context, request *dto.SubmitCommentRequest) (*dto.SubmitCommentResponse, error)

	// SubmitReply 提交回复
	// 参数:
	// - ctx: 上下文
	// - request: 回复提交请求,包含回复内容、视频ID、父评论ID等信息
	// 返回:
	// - *dto.SubmitReplyResponse: 提交回复响应
	// - error: 提交过程中的错误信息
	SubmitReply(ctx context.Context, request *dto.SubmitReplyRequest) (*dto.SubmitReplyResponse, error)

	// GetCommentsByVideoID 根据视频ID获取评论列表
	// 参数:
	// - ctx: 上下文
	// - request: 获取评论请求,包含视频ID、分页等参数
	// 返回:
	// - *dto.GetCommentsResponse: 评论列表响应
	// - error: 获取过程中的错误信息
	GetCommentsByVideoID(ctx context.Context, request *dto.GetCommentsRequest) (*dto.GetCommentsResponse, error)

	// GetCommentsByRootID 获取指定根评论下的所有子评论
	// 参数:
	// - ctx: 上下文
	// - request: 获取回复请求,包含根评论ID、分页等参数
	// 返回:
	// - *dto.GetReplyResponse: 回复列表响应
	// - error: 获取过程中的错误信息
	GetCommentsByRootID(ctx context.Context, request *dto.GetReplyRequest) (*dto.GetReplyResponse, error)
}
