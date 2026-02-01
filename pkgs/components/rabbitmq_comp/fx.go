package rabbitmq_comp

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"go.uber.org/fx"
)

var (
	RabbitMQComponentFx = fx.Module(
		"rabbitmq",
		rabbitMQProviders,
		rabbitMQInvokes,
	)

	rabbitMQProviders = fx.Options(
		fx.Provide(LoadRabbitMQConfig),
		fx.Provide(NewRabbitMQComponent),
		fx.Provide(ProvideRabbitMQClient),
	)

	rabbitMQInvokes = fx.Options(
		fx.Invoke(registerRabbitMQHooks),
	)
)

func ProvideRabbitMQClient(comp *RabbitMQComponent) IRabbitMQClient {
	return comp.GetClient()
}

func registerRabbitMQHooks(
	lc fx.Lifecycle,
	comp *RabbitMQComponent,
	log logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if comp.GetClient() == nil {
				return nil
			}
			log.Info("RabbitMQ connection verified")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := comp.Stop(); err != nil {
				log.Errorf("error closing RabbitMQ: %v", err)
				return err
			}
			log.Info("RabbitMQ closed gracefully")
			return nil
		},
	})
}
