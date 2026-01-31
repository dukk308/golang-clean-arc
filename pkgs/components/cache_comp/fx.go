package redis_component

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/logger"
	"go.uber.org/fx"
)

var (
	CacheComponent = fx.Module(
		"redis",
		redisProviders,
		redisInvokes,
	)

	redisProviders = fx.Options(
		fx.Provide(NewRedisComponent),
		fx.Provide(ProvideCacheService),
	)

	redisInvokes = fx.Options(
		fx.Invoke(registerHooks),
	)
)

func ProvideCacheService(redisComponent *RedisComponent) ICacheService {
	return redisComponent.GetClient()
}

func registerHooks(
	lc fx.Lifecycle,
	redisComponent *RedisComponent,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			client := redisComponent.GetClient()
			if client == nil {
				return nil
			}
			if err := client.Ping(ctx); err != nil {
				logger.Errorf("failed to ping redis: %v", err)
				return err
			}
			logger.Info("redis connection verified")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := redisComponent.Stop(); err != nil {
				logger.Errorf("error closing redis: %v", err)
				return err
			}
			logger.Info("redis closed gracefully")
			return nil
		},
	})
}
