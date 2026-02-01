package otel_comp

import (
	"flag"
)

var (
	tracerExporter  = flag.String("tracer-exporter", "otlp-grpc", "Tracer exporter type (noop, memory, stdout, otlp-grpc)")
	tracerCollector = flag.String("tracer-collector", "localhost:4320", "OTLP collector endpoint (host:port)")
)

func LoadTracerConfig(enabled bool) *TracerConfig {
	collector := *tracerCollector
	exporter := *tracerExporter

	return &TracerConfig{
		Enabled:   enabled,
		Exporter:  exporter,
		Collector: collector,
	}
}

type TracerConfig struct {
	Enabled  bool
	Exporter string
	Collector string
}
