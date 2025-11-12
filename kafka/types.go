package kafka

type Message struct {
	Key []byte
	Value []byte
	Topic string
}

type Ack func() error

type Handler func(msg Message, ack Ack) error

type Config struct {
	Brokers []string
	GroupID string
	AutoAck bool
	Workers int
	QueueSize int
}


type Client interface {
	Produce(topic string, key, value []byte) error
	Consume(handler Handler) error
	Close() error
}