package consumer

import (
	"context"
	"fmt"
	"gateService/internal/domain/repository"
	"gateService/pkg/mq/nsqpool"
	"log"
	"strconv"
)

type AccountConsumer struct {
	userRepository      repository.UserRepository
	accountConsumerPool *nsqpool.ConsumerPool
}

func NewAccountConsumer(userRepository repository.UserRepository) *AccountConsumer {
	return &AccountConsumer{
		userRepository: userRepository,
	}
}

func (c *AccountConsumer) DeleteTestAccount(ctx context.Context, msg []byte) error {
	userID, err := strconv.Atoi(string(msg))
	if err != nil {
		return fmt.Errorf("转换用户ID失败: %v", err)
	}

	tx, err := c.userRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer tx.Rollback()

	err = c.userRepository.DeleteUser(ctx, tx, userID)
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

func (c *AccountConsumer) Start() {
	accountConsumerPool, err := nsqpool.NewConsumerPool(
		&nsqpool.ConsumerOptions{
			Topic:    "test_account_queue",
			Channel:  "test_account_channel",
			PoolSize: 2,
		},
	)
	if err != nil {
		log.Fatalf("创建体验账号消费者池失败: %v", err)
	}
	c.accountConsumerPool = accountConsumerPool

	accountConsumerPool.RegisterCallback(c.DeleteTestAccount)

	err = accountConsumerPool.Start()
	if err != nil {
		log.Fatalf("启动体验账号消费者池失败: %v", err)
	}
}

func (c *AccountConsumer) Stop() {
	c.accountConsumerPool.Stop()
}
