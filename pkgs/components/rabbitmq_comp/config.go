package rabbitmq_comp

type RabbitMQConfig struct {
	URL           string
	PrefetchCount int
	PrefetchSize  int
}
