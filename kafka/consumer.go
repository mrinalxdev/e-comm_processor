package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
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
	for i := 0; i < c.config.Workers; i ++ {
		c.wg.Add(1)
		go c.worker(i, handler)
	}
	
	return nil
}

func (c *kafkaConsumer) worker(id int, handler Handler){
	defer c.wg.Done()
	
	for {
		select {
			case <- c.ctx.Done():
				return
			
			
			
			default:
				msg, err := c.reader.FetchMessage(c.ctx)
				if err != nil {
					if err == context.Canceled{
						return
					}
					logrus.Errorf("Consumer worker %d failed to fetch message %v", id, err)
					
					continue
				}
				
				ack := func() error {
					return c.reader.CommitMessages(c.ctx, msg)
				}
				
				kafkaMsg := Message {
					Key : msg.Key,
					Value : msg.Value,
					Topic : msg.Topic,
				}
				
				if err := handler(kafkaMsg, ack); err != nil {
					logrus.Errorf("Handler failed in worker %d : %v", id, err)
				}
				
				if c.config.AutoAck {
					if err := ack(); err != nil {
						logrus.Errorf("Failed to auto ack message in worker %d : %v", id, err)
					}
				}
		}
	}
}

func (c *kafkaConsumer) close() {
	c.reader.Close()
}