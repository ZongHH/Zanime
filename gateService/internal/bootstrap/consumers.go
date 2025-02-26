package bootstrap

import "gateService/internal/application/consumer"

type consumers struct {
	OrderConsumer   *consumer.OrderConsumer
	CommentConsumer *consumer.CommentConsumer
}

func initConsumers(bases *bases, repositories *repositories) *consumers {
	return &consumers{
		OrderConsumer:   consumer.NewOrderConsumer(repositories.OrderRepo),
		CommentConsumer: consumer.NewCommentConsumer(repositories.PostRepo, repositories.PostCommentRepo, repositories.UserRepo, bases.WebSocketManager),
	}
}

func (c *consumers) Start() {
	c.OrderConsumer.Start()
	c.CommentConsumer.Start()
}

func (c *consumers) Close() {
	c.OrderConsumer.Stop()
	c.CommentConsumer.Stop()
}
