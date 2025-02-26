package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

// OrderRepository 定义订单相关的数据库操作接口
// 提供订单的增删改查等基本操作方法
type OrderRepository interface {
	// CreateOrder 创建新订单
	// 参数:
	// - ctx: 上下文
	// - order: 需要创建的订单实体
	// 返回:
	// - error: 创建过程中的错误信息
	CreateOrder(ctx context.Context, order *entity.Order) error

	// UpdateOrder 更新订单信息
	// 参数:
	// - ctx: 上下文
	// - order: 更新后的订单信息
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateOrder(ctx context.Context, order *entity.Order) error

	// DeleteOrder 删除订单
	// 参数:
	// - ctx: 上下文
	// - orderID: 要删除的订单ID
	// 返回:
	// - error: 删除过程中的错误信息
	DeleteOrder(ctx context.Context, orderID string) error

	// GetOrdersByUserID 获取用户的订单列表
	// 参数:
	// - ctx: 上下文
	// - userID: 用户ID
	// - page: 分页页码,从1开始
	// - limit: 每页数量限制
	// 返回:
	// - []*entity.Order: 订单列表
	// - error: 获取过程中的错误信息
	GetOrdersByUserID(ctx context.Context, userID, page, limit int) ([]*entity.Order, error)

	// UpdateOrderStatus 更新订单状态
	// 参数:
	// - ctx: 上下文
	// - orderID: 订单ID
	// - status: 新的订单状态
	// 返回:
	// - error: 更新过程中的错误信息
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}
