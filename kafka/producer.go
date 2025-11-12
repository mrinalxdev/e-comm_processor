package kafka

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type kafkaProducer struct {
	config Config
	writer *kafka.Writer
	queue chan Message
	ctx context.Context
	wg *sync.WaitGroup
}

func newProducer(config Config, ctx context.Context, wg *sync.WaitGroup) *kafkaProducer {
	writer := &kafka.Writer {
		Addr : kafka.TCP(config.Brokers...),
		Balancer: &kafka.LeastBytes{},
		BatchSize : 100,
		BatchTimeout: 10 * time.Millisecond,
		Async: true,
	}
	
	
	producer := &kafkaProducer {
		config : config,
		writer : writer,
		queue : make(chan Message, config.QueueSize),
		ctx : ctx,
		wg : wg,
	}
	
	
	producer.startWorkers()
	
	return producer
}


func (p *kafkaProducer) startWorkers(){
	for i := 0; i < p.config.Workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *kafkaProducer) worker(id int){
	defer p.wg.Done()
	
	for {
		select {
			case msg := <- p.queue :
				err := p.writer.WriteMessages(p.ctx, kafka.Message{
					Topic : msg.Topic,
					Key : msg.Key,
					Value : msg.Value,
				})
				
				if err != nil {
					logrus.Errorf("Producer worker %d failed to write message : %v", id, err)
				}
				
			
			case <- p.ctx.Done():
				return
		}
	}
}

func (p *kafkaProducer) produce(msg Message) error {
	select {
	case p.queue <- msg:
		return nil
	
	case <- p.ctx.Done():
		return p.ctx.Err()
	
	default:
		return ErrQueueFull
		
	}
}

var ErrQueueFull = errors.New("producer queue is full")