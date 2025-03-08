package nsqpool

import (
	"fmt"
	"managerService/pkg/config"

	"sync"

	"github.com/nsqio/go-nsq"
)

// NSQConsumerPool表示NSQ消费者连接池结构体
type ConsumerPool struct {
	mu           sync.Mutex
	consumers    []*nsq.Consumer
	topic        string
	channel      string
	config       *nsq.Config
	callBackFunc func(msg []byte) error
}

type ConsumerConfig struct {
	Topic    string
	Channel  string
	PoolSize int
	config   *nsq.Config
}

func NewConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		config: nsq.NewConfig(),
	}
}

// NewNSQConsumerPool创建一个新的NSQ消费者连接池实例
func NewConsumerPool(config *ConsumerConfig) (*ConsumerPool, error) {
	checkConfig(config)
	consumers := make([]*nsq.Consumer, config.PoolSize)

	for i := 0; i < config.PoolSize; i++ {
		consumer, err := nsq.NewConsumer(config.Topic, config.Channel, config.config)
		if err != nil {
			return nil, fmt.Errorf("创建NSQ消费者实例失败: %v", err)
		}

		// 将消费者添加到切片中
		consumers[i] = consumer
	}

	return &ConsumerPool{
		consumers: consumers,
		topic:     config.Topic,
		channel:   config.Channel,
		config:    config.config,
	}, nil
}

func checkConfig(config *ConsumerConfig) {
	if config.Topic == "" {
		config.Topic = "ServiceLog"
	}
	if config.Channel == "" {
		config.Channel = "update"
	}
	if config.PoolSize == 0 {
		config.PoolSize = 2
	}
}

// 消息回调处理
func (p *ConsumerPool) HandleMessage(msg *nsq.Message) error {
	// 执行用户定义的业务处理逻辑
	if err := p.callBackFunc(msg.Body); err != nil {
		return err
	}
	// 倘若业务处理成功，则调用Finish方法，发送消息的ack
	msg.Finish()
	return nil
}

// Close关闭NSQ消费者连接池，释放所有资源
func (p *ConsumerPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, consumer := range p.consumers {
		consumer.Stop()
	}
}

// RegisterCallbackFunc用于外部注册处理函数
func (p *ConsumerPool) RegisterCallbackFunc(callback func(msg []byte) error) {
	p.callBackFunc = callback
}

func (p *ConsumerPool) Start() error {
	nsqdaddr, err := config.GetHostAndPort("nsq")
	if err != nil {
		return err
	}
	for _, consumer := range p.consumers {
		consumer.AddHandler(nsq.HandlerFunc(p.HandleMessage))

		if err := consumer.ConnectToNSQD(nsqdaddr); err != nil {
			return err
		}
	}
	return nil
}
