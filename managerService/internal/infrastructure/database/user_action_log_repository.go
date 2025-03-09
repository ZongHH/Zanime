package database

import (
	"context"
	"database/sql"
	"managerService/internal/domain/entity"
)

type UserActionLogRepositoryImpl struct {
	db *sql.DB
}

func NewUserActionLogRepository(db *sql.DB) *UserActionLogRepositoryImpl {
	return &UserActionLogRepositoryImpl{
		db: db,
	}
}

func (r *UserActionLogRepositoryImpl) GetUserActionLogs(ctx context.Context, page int, pageSize int) ([]*entity.UserActionLog, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT ua.id, ua.user_id, ua.action, ua.module, ua.ip_address, ua.status, ua.message, ua.created_at, u.username, u.account_type
		FROM user_action_logs ua
		LEFT JOIN user_infos u ON ua.user_id = u.user_id
		ORDER BY ua.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]*entity.UserActionLog, 0)
	for rows.Next() {
		var log entity.UserActionLog
		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Module, &log.IPAddress, &log.Status, &log.Message, &log.CreatedAt, &log.UserName, &log.UserType)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, nil
}
