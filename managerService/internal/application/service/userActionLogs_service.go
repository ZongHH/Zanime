package service

import (
	"context"
	"managerService/internal/domain/repository"
	"managerService/internal/interfaces/dto/request"
	"managerService/internal/interfaces/dto/response"
)

type UserActionLogsServiceImpl struct {
	userActionLogRepository repository.UserActionLogsRepository
}

func NewUserActionLogsServiceImpl(userActionLogRepository repository.UserActionLogsRepository) *UserActionLogsServiceImpl {
	return &UserActionLogsServiceImpl{
		userActionLogRepository: userActionLogRepository,
	}
}

func (s *UserActionLogsServiceImpl) GetUserActionLogs(ctx context.Context, request *request.UserActionLogsRequest) (*response.UserActionLogsResponse, error) {
	userActionLogs, err := s.userActionLogRepository.GetUserActionLogs(ctx, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	userActionLogsResponse := make([]*response.UserActionLog, 0)
	for _, userActionLog := range userActionLogs {
		userActionLogsResponse = append(userActionLogsResponse, &response.UserActionLog{
			ID:       userActionLog.ID,
			UserName: userActionLog.UserName,
			UserType: userActionLog.UserType,
			Action:   userActionLog.Action,
			Module:   userActionLog.Module,
			Time:     userActionLog.CreatedAt,
		})
	}
	return &response.UserActionLogsResponse{
		Code:           200,
		UserActionLogs: userActionLogsResponse,
	}, nil
}
