package monitor

import (
	"context"
	"crawler/pkg/mq/nsqpool"
	"fmt"
)

type log_producer struct {
	producer    *nsqpool.ProducerPool
	ServiceName string
	TopicName   string
}

func newLogProducer(config *logConfig) (*log_producer, error) {
	producer, err := nsqpool.NewProducerPool(&nsqpool.ProducerOptions{
		PoolSize: config.PoolSize,
	})
	if err != nil {
		return nil, fmt.Errorf("nsqpool.NewProducerPool %v", err)
	}
	lp := &log_producer{
		producer:    producer,
		ServiceName: config.ServiceName,
		TopicName:   config.TopicName,
	}
	return lp, nil
}

func (lp *log_producer) publish(msg []byte) error {
	err := lp.producer.Publish(context.Background(), lp.TopicName, msg)
	if err != nil {
		return fmt.Errorf("lp.producer.Publish %v", err)
	}
	return nil
}

func (lp *log_producer) close() {
	lp.producer.Close()
}
