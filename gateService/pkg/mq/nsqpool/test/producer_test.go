package test

import (
	"context"
	"fmt"
	"gateService/pkg/mq/nsqpool"
	"math/rand/v2"
	"sync"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	// PoolSize:      100,
	// producer.GetMetrics(): map[
	// active_producers:0
	// available_count:100
	// error_count:19198
	// error_rate:0.19198
	// message_count:100000
	// processing_time:171695174000]
	// goroutine: 100000
	// ctxWithTimeout: 500ms

	// PoolSize:      300,
	// producer.GetMetrics(): map[
	// active_producers:0
	// available_count:300
	// error_count:11664
	// error_rate:0.11664
	// message_count:100000
	// processing_time:458497338200]
	// goroutine: 100000
	// ctxWithTimeout: 1000ms

	// PoolSize:      1000,
	// producer.GetMetrics(): map[
	// active_producers:0
	// available_count:1000
	// error_count:9861
	// error_rate:0.09861
	// message_count:100000
	// processing_time:1640026643600]
	// goroutine: 100000
	// ctxWithTimeout: 500ms

	producer, err := nsqpool.NewProducerPool(&nsqpool.ProducerOptions{
		PoolSize:      100,
		EnableMetrics: true,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancelWithTimeout()

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.IntN(3000)) * time.Millisecond)
			err = producer.PublishAsync(ctxWithTimeout, "test", []byte("test"))
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
		}()
	}
	wg.Wait()

	producer.Close()

	fmt.Printf("producer.GetMetrics(): %v\n", producer.GetMetrics())
}
