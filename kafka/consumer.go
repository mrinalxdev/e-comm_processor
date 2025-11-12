package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	config Config
	reader *kafka.Reader
	ctx context.Context
	wg *sync.WaitGroup
}

func newConsumer(config Config, ctx context.Context, wg *sync.WaitGroup) *kafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		GroupID : config.GroupID,
	})
	
	
	return &kafkaConsumer{
		config : config,
		reader: reader,
		ctx : ctx,
		wg : wg,
	}
}

func (c *kafkaConsumer) consume(handler Handler) error {
	
}