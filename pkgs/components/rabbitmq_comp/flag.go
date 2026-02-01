package rabbitmq_comp

import (
	"flag"
)

var (
	rabbitMQURL = flag.String("rabbitmq-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ connection URL")
)

func LoadRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		URL:           *rabbitMQURL,
		PrefetchCount: 1,
		PrefetchSize:  0,
	}
}
