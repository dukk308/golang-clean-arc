package otel_comp

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/global_config"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var TracerModule = fx.Module(
	"tracer",
	fx.Provide(
		NewDefaultTracerProviderFactory,
		NewTracerConfig,
		NewFxTracerProvider,
	),
	fx.Invoke(func(*trace.TracerProvider) {}),
)

type TracerConfigParam struct {
	fx.In
	GlobalConfig *global_config.GlobalConfig `optional:"true"`
}

func NewTracerConfig(p TracerConfigParam) *TracerConfig {
	enabled := false
	if p.GlobalConfig != nil {
		enabled = p.GlobalConfig.EnableTracing
	}
	return LoadTracerConfig(enabled)
}

type FxTracerParam struct {
	fx.In
	LifeCycle    fx.Lifecycle
	Factory      TracerProviderFactory
	ServiceName  string `name:"serviceName"`
	TracerConfig *TracerConfig
	Logger       logger.Logger
}

func NewFxTracerProvider(p FxTracerParam) (*trace.TracerProvider, error) {
	if p.TracerConfig == nil {
		p.Logger.Info("Tracer config is not provided, using noop exporter")
		tracerProvider, err := p.Factory.Create(
			WithName(p.ServiceName),
			WithExporter(Noop),
		)
		if err != nil {
			return nil, err
		}
		return tracerProvider, nil
	}

	if !p.TracerConfig.Enabled {
		p.Logger.Info("Tracer is disabled, using noop exporter")
		tracerProvider, err := p.Factory.Create(
			WithName(p.ServiceName),
			WithExporter(Noop),
		)
		if err != nil {
			return nil, err
		}
		return tracerProvider, nil
	}

	exporter := FetchExporter(p.TracerConfig.Exporter)
	if exporter == Noop {
		p.Logger.Info("Invalid exporter type, using noop exporter")
		tracerProvider, err := p.Factory.Create(
			WithName(p.ServiceName),
			WithExporter(Noop),
		)
		if err != nil {
			return nil, err
		}
		return tracerProvider, nil
	}

	p.Logger.Info("Initializing tracer",
		"service", p.ServiceName,
		"exporter", p.TracerConfig.Exporter,
		"collector", p.TracerConfig.Collector,
	)

	tracerProvider, err := p.Factory.Create(
		WithName(p.ServiceName),
		WithExporter(exporter),
		WithCollector(p.TracerConfig.Collector),
	)
	if err != nil {
		p.Logger.Error("error creating tracer provider", "error", err)
		return nil, err
	}

	p.Logger.Info("Tracer provider initialized successfully",
		"service", p.ServiceName,
		"exporter", p.TracerConfig.Exporter,
		"collector", p.TracerConfig.Collector,
	)

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if exporter == Noop || exporter == Memory {
				return nil
			}

			if err = tracerProvider.ForceFlush(ctx); err != nil {
				p.Logger.Error("error flushing tracer provider", err)
				return err
			}

			if err = tracerProvider.Shutdown(ctx); err != nil {
				p.Logger.Error("error while shutting down tracer provider", err)
				return err
			}

			return nil
		},
	})

	return tracerProvider, nil
}
