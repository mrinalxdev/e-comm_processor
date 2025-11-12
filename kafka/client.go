package kafka

import (
	"context"
	"sync"
)




type kafkaClient struct {
	config    Config
	producer  *kafkaProducer
	consumer  *kafkaConsumer
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewClient(config Config) Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &kafkaClient{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}
