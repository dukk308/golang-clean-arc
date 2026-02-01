# RabbitMQ Component

FX component for RabbitMQ using `github.com/rabbitmq/amqp091-go`.

## Config

- `--rabbitmq-url`: Connection URL (default: `amqp://guest:guest@localhost:5672/`)

## Usage

Wire the module in your bootstrap:

```go
import "github.com/dukk308/beetool.dev-go-starter/pkgs/components/rabbitmq_comp"

fx.Options(
    rabbitmq_comp.RabbitMQComponentFx,
    // ...
)
```

Inject `IRabbitMQClient` to publish or consume:

```go
type MyService struct {
    mq rabbitmq_comp.IRabbitMQClient
}

func (s *MyService) Publish(ctx context.Context, body []byte) error {
    return s.mq.Publish(ctx, "my-exchange", "my-key", false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        body,
    })
}
```

For new channels (e.g. separate consumers), use `Channel()` and manage the channel lifecycle yourself. The component closes the default connection and channel on app shutdown.
