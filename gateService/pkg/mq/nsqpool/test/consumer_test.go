package test

import (
	"context"
	"fmt"
	"gateService/pkg/mq/nsqpool"
	"math/rand/v2"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	consumer, err := nsqpool.NewConsumerPool(&nsqpool.ConsumerOptions{
		Topic:         "test",
		Channel:       "test",
		PoolSize:      200,
		EnableMetrics: true,
		MaxInFlight:   100,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	consumer.RegisterCallback(func(ctx context.Context, msg []byte) error {
		time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond)
		return nil
	})

	consumer.Start()

	time.Sleep(200 * time.Second)
	fmt.Printf("consumer.GetMetrics(): %v\n", consumer.GetMetrics())

	consumer.Stop()
}
