package rabbitmq_comp

import (
	"context"
	"sync"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQComponent struct {
	config     *RabbitMQConfig
	conn       *amqp.Connection
	channel    *amqp.Channel
	client     IRabbitMQClient
	mu         sync.Mutex
}

type rabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *RabbitMQConfig
}

func NewRabbitMQComponent(config *RabbitMQConfig, log logger.Logger) *RabbitMQComponent {
	c := &RabbitMQComponent{config: config}
	if err := c.connect(log); err != nil {
		log.Errorf("Failed to connect to RabbitMQ: %v", err)
		return c
	}
	c.client = &rabbitMQClient{conn: c.conn, channel: c.channel, config: config}
	return c
}

func (c *RabbitMQComponent) connect(log logger.Logger) error {
	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		return err
	}
	c.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if c.config.PrefetchCount > 0 || c.config.PrefetchSize > 0 {
		if err := ch.Qos(c.config.PrefetchCount, c.config.PrefetchSize, false); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return err
		}
	}
	c.channel = ch
	log.Info("Connected to RabbitMQ")
	return nil
}

func (c *RabbitMQComponent) GetClient() IRabbitMQClient {
	return c.client
}

func (c *RabbitMQComponent) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	if c.channel != nil {
		if e := c.channel.Close(); e != nil {
			err = e
		}
		c.channel = nil
	}
	if c.conn != nil {
		if e := c.conn.Close(); e != nil {
			if err == nil {
				err = e
			}
		}
		c.conn = nil
	}
	return err
}

func (r *rabbitMQClient) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return r.channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (r *rabbitMQClient) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (r *rabbitMQClient) Channel() (*amqp.Channel, error) {
	return r.conn.Channel()
}

func (r *rabbitMQClient) Close() error {
	return nil
}
