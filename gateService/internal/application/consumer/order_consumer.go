package consumer

import (
	"context"
	"gateService/internal/domain/repository"
	"gateService/pkg/mq/nsqpool"
	"log"
)

type OrderConsumer struct {
	orderRepository repository.OrderRepository
	consumerPool    *nsqpool.ConsumerPool
}

func NewOrderConsumer(orderRepository repository.OrderRepository) *OrderConsumer {
	return &OrderConsumer{orderRepository: orderRepository}
}

func (o *OrderConsumer) updateOrderStatus(ctx context.Context, msg []byte) error {
	orderID := string(msg)
	o.orderRepository.UpdateOrderStatus(ctx, orderID, "已失效")
	return nil
}

func (o *OrderConsumer) Start() {
	consumerPool, err := nsqpool.NewConsumerPool(&nsqpool.ConsumerOptions{
		Topic:    "Deferredorders",
		Channel:  "update_order_status",
		PoolSize: 3,
	})
	if err != nil {
		log.Fatalf("创建订单状态消费者池失败: %v\n", err)
	}
	o.consumerPool = consumerPool

	consumerPool.RegisterCallback(o.updateOrderStatus)
	err = consumerPool.Start()
	if err != nil {
		log.Fatalf("启动订单状态消费者池失败: %v\n", err)
	}
}

func (o *OrderConsumer) Stop() {
	o.consumerPool.Stop()
}
