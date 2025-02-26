package service

import (
	"context"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/interfaces/dto"
)

type CommentServiceImpl struct {
	commentRepository repository.CommentRepository
	userRepository    repository.UserRepository
}

func NewCommentServiceImpl(commentRepository repository.CommentRepository, userRepository repository.UserRepository) *CommentServiceImpl {
	return &CommentServiceImpl{
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (c *CommentServiceImpl) SubmitComment(ctx context.Context, request *dto.SubmitCommentRequest) (*dto.SubmitCommentResponse, error) {
	commentEntity := &entity.Comment{
		VideoID: request.VideoID,
		UserID:  request.UserID,
		Content: request.Content,
	}
	err := c.commentRepository.CreateFirstLevelComment(ctx, commentEntity)
	if err != nil {
		return nil, fmt.Errorf("创建主评论失败: %v", err)
	}
	return &dto.SubmitCommentResponse{Code: 200}, nil
}

func (c *CommentServiceImpl) GetCommentsByVideoID(ctx context.Context, request *dto.GetCommentsRequest) (*dto.GetCommentsResponse, error) {
	comments, err := c.commentRepository.GetCommentsByVideoID(ctx, request.VideoID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取视频评论失败: %v", err)
	}

	totalCount, err := c.commentRepository.GetTotalCount(ctx, request.VideoID)
	if err != nil {
		return nil, fmt.Errorf("获取视频评论总数失败: %v", err)
	}
	totalPage := (totalCount + request.PageSize - 1) / request.PageSize

	commentsDto := make([]*dto.Comment, 0, len(comments))
	for _, comment := range comments {
		commentsDto = append(commentsDto, &dto.Comment{
			CommentID: comment.CommentID,
			Content:   comment.Content,
			UserInfo: dto.UserInfo{
				ID:        comment.UserID,
				Username:  comment.UserName,
				AvatarURL: comment.AvatarURL,
			},
			CreatedAt: comment.CreatedAt,
			ReplyNum:  comment.ReplyNum,
		})
	}
	return &dto.GetCommentsResponse{
		Code:      200,
		TotalPage: totalPage,
		Comments:  commentsDto,
	}, nil
}

func (c *CommentServiceImpl) SubmitReply(ctx context.Context, request *dto.SubmitReplyRequest) (*dto.SubmitReplyResponse, error) {
	commentEntity := &entity.Comment{
		VideoID:  request.VideoID,
		UserID:   request.UserID,
		Content:  request.Content,
		ToUserID: &request.ToUserID,
		ParentID: &request.ParentID,
		RootID:   &request.RootID,
		Level:    2,
	}

	err := c.commentRepository.CreateSecondLevelComment(ctx, commentEntity)
	if err != nil {
		return nil, fmt.Errorf("创建回复失败: %v", err)
	}
	return &dto.SubmitReplyResponse{Code: 200}, nil
}

func (c *CommentServiceImpl) GetCommentsByRootID(ctx context.Context, request *dto.GetReplyRequest) (*dto.GetReplyResponse, error) {
	replies, err := c.commentRepository.GetCommentsByRootID(ctx, request.RootID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取回复失败: %v", err)
	}

	toUserIDs := make([]int, 0, len(replies))
	for _, reply := range replies {
		if reply.ToUserID != nil {
			toUserIDs = append(toUserIDs, *reply.ToUserID)
		}
	}

	toUsers, err := c.userRepository.GetUsersByIDs(ctx, &toUserIDs)
	if err != nil {
		return nil, fmt.Errorf("获取被回复用户名失败: %v", err)
	}

	repliesDto := make([]*dto.Reply, 0, len(replies))
	for i, reply := range replies {
		repliesDto = append(repliesDto, &dto.Reply{
			CommentID: reply.CommentID,
			RootID:    *reply.RootID,
			Content:   reply.Content,
			UserInfo: dto.UserInfo{
				ID:        reply.UserID,
				Username:  reply.UserName,
				AvatarURL: reply.AvatarURL,
			},
			CreatedAt:   reply.CreatedAt,
			RepliedName: (*toUsers)[i].Username,
		})
	}

	return &dto.GetReplyResponse{
		Code:    200,
		Replies: repliesDto,
	}, nil
}
