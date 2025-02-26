package service

import (
	"context"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/middleware/lock"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/mq/nsqpool"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderServiceImpl struct {
	redisClient       *redis.Client
	lockOptions       *lock.LockOptions
	producerPool      *nsqpool.ProducerPool
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository, productRepository repository.ProductRepository, producerPool *nsqpool.ProducerPool, redisClient *redis.Client, lockOptions *lock.LockOptions) *OrderServiceImpl {
	return &OrderServiceImpl{orderRepository: orderRepository, productRepository: productRepository, producerPool: producerPool, redisClient: redisClient, lockOptions: lockOptions}
}

func (o *OrderServiceImpl) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	lockKey := fmt.Sprintf("product%d:size%s:color%s", request.ProductID, request.Size, request.Color)
	lock := lock.NewRedisLock(o.redisClient, lockKey, o.lockOptions)

	if err := lock.WaitLock(ctx); err != nil {
		return nil, fmt.Errorf("获取订单锁失败: %w", err)
	}
	defer lock.Unlock(ctx)

	stock, err := o.productRepository.GetProductStock(ctx, request.ProductID, request.Size, request.Color)
	if err != nil {
		return nil, fmt.Errorf("获取商品库存失败: %w", err)
	}
	if stock < 1 {
		return nil, fmt.Errorf("库存不足")
	}

	err = o.productRepository.DecreaseProductStock(ctx, request.ProductID, request.Size, request.Color, 1)
	if err != nil {
		return nil, fmt.Errorf("减少商品库存失败: %w", err)
	}

	// 根据订单信息生成唯一的orderID
	orderID := fmt.Sprintf("%s_%d_%d_%s",
		time.Now().Format("20060102150405"), // 时间戳
		request.UserID,                      // 用户ID
		request.ProductID,                   // 商品ID
		uuid.New().String()[0:8])            // UUID前8位

	order := &entity.Order{
		OrderID:       orderID,
		UserID:        request.UserID,
		ProductID:     request.ProductID,
		ProductName:   request.ProductName,
		SelectedSize:  request.Size,
		SelectedColor: request.Color,
		Price:         request.Price,
		Phone:         request.Phone,
		UserName:      request.UserName,
		Address:       request.Address,
		Status:        "待支付",
		CreateTime:    time.Now().Format("2006-01-02 15:04:05"),
	}

	err = o.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		o.productRepository.IncreaseProductStock(ctx, request.ProductID, request.Size, request.Color, 1)
		return nil, err
	}

	err = o.producerPool.DeferredPublish(ctx, "Deferredorders", 5*time.Minute, []byte(orderID))
	if err != nil {
		o.productRepository.IncreaseProductStock(ctx, request.ProductID, request.Size, request.Color, 1)
		return nil, err
	}

	return &dto.CreateOrderResponse{
		Code:        200,
		OrderID:     orderID,
		UserID:      order.UserID,
		ProductID:   order.ProductID,
		ProductName: order.ProductName,
		Size:        order.SelectedSize,
		Color:       order.SelectedColor,
		Price:       order.Price,
		Phone:       order.Phone,
		UserName:    order.UserName,
		Address:     order.Address,
		CreatedTime: order.CreateTime,
	}, nil
}

func (o *OrderServiceImpl) GetOrdersByUserID(ctx context.Context, request *dto.GetOrderRequest) (*dto.GetOrderResponse, error) {
	orders, err := o.orderRepository.GetOrdersByUserID(ctx, request.UserID, request.Page, request.PageSize)
	if err != nil {
		return &dto.GetOrderResponse{Code: 500, Orders: nil}, fmt.Errorf("获取用户订单失败: %v", err)
	}

	return &dto.GetOrderResponse{
		Code:   200,
		Orders: orders,
	}, nil
}

func (o *OrderServiceImpl) CallbackPay(ctx context.Context, request *dto.PayRequest) (*dto.PayResponse, error) {
	err := o.orderRepository.UpdateOrderStatus(ctx, request.OrderID, "已支付")
	if err != nil {
		return &dto.PayResponse{Code: 500}, fmt.Errorf("更新订单状态失败: %v", err)
	}

	return &dto.PayResponse{Code: 200}, nil
}
