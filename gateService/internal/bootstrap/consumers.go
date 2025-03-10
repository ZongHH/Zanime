package bootstrap

import (
	"gateService/internal/application/consumer"
	"gateService/internal/infrastructure/config"
)

type consumers struct {
	OrderConsumer   *consumer.OrderConsumer
	CommentConsumer *consumer.CommentConsumer
	AccountConsumer *consumer.AccountConsumer
}

func initConsumers(cfg *config.Config, bases *bases, repositories *repositories) *consumers {
	return &consumers{
		OrderConsumer:   consumer.NewOrderConsumer(repositories.OrderRepo),
		CommentConsumer: consumer.NewCommentConsumer(repositories.PostRepo, repositories.PostCommentRepo, repositories.UserRepo, bases.WebSocketManager),
		AccountConsumer: consumer.NewAccountConsumer(repositories.UserRepo),
	}
}

func (c *consumers) Start() {
	c.OrderConsumer.Start()
	c.CommentConsumer.Start()
	c.AccountConsumer.Start()
}

func (c *consumers) Close() {
	c.OrderConsumer.Stop()
	c.CommentConsumer.Stop()
	c.AccountConsumer.Stop()
}
