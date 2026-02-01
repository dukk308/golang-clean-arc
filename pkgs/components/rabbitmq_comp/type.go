package rabbitmq_comp

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMQClient interface {
	Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Channel() (*amqp.Channel, error)
	Close() error
}
