package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// OrderService 定义订单相关的业务逻辑接口
// 提供订单的创建、查询、支付回调等功能
type OrderService interface {
	// CreateOrder 创建新订单
	// 参数:
	// - ctx: 上下文
	// - request: 创建订单请求,包含商品信息、用户信息等
	// 返回:
	// - *dto.CreateOrderResponse: 创建订单响应,包含订单ID等信息
	// - error: 创建过程中的错误信息
	CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)

	// GetOrdersByUserID 获取用户订单列表
	// 参数:
	// - ctx: 上下文
	// - request: 获取订单请求,包含用户ID、分页等参数
	// 返回:
	// - *dto.GetOrderResponse: 订单列表响应
	// - error: 获取过程中的错误信息
	GetOrdersByUserID(ctx context.Context, request *dto.GetOrderRequest) (*dto.GetOrderResponse, error)

	// CallbackPay 处理支付回调
	// 参数:
	// - ctx: 上下文
	// - request: 支付回调请求,包含支付结果等信息
	// 返回:
	// - *dto.PayResponse: 支付回调处理响应
	// - error: 处理过程中的错误信息
	CallbackPay(ctx context.Context, request *dto.PayRequest) (*dto.PayResponse, error)
}
