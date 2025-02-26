package collect

import (
	"log"
	"monitorService/pkg/mq/nsqpool"
)

var consumer *nsqpool.ConsumerPool
var logCache *LogCache

func Init() {
	logCache = InitCache()
	var err error
	consumer, err = nsqpool.NewConsumerPool(nsqpool.NewConsumerConfig())
	if err != nil {
		log.Fatal("nsqpool.NewConsumerPool ", err)
	}
	consumer.RegisterCallbackFunc(collect)
	if err := consumer.Start(); err != nil {
		log.Fatal("consumer.Start() ", err)
	}
}

func Close() {
	if consumer != nil {
		consumer.Close()
	}
	if logCache != nil {
		logCache.Stop()
	}
}
