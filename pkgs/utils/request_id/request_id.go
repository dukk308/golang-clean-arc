package request_id

import (
	"context"
	"net/http"
	"os"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/constants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Provider represents the interface for the request id provider.
type Provider func(http.ResponseWriter, *http.Request) (string, error)

// Factory represents the function to generate a request id.
type Factory func() (string, error)

// RequestID reads the request id from the current request header, and creates
// a new one if it does not exist. If OpenTelemetry tracing is enabled, it will
// prioritize using the trace ID from the span context.
func RequestID(factory Factory) Provider {
	return func(w http.ResponseWriter, r *http.Request) (string, error) {
		var id string

		if traceID := getTraceIDFromContext(r.Context()); traceID != "" {
			id = traceID
		} else {
			id = r.Header.Get(constants.ContextKeyRequestID)
		}

		if id == "" {
			var err error
			id, err = factory()
			if err != nil {
				return "", err
			}
			r.Header.Set(constants.ContextKeyRequestID, id)
		}

		w.Header().Set(constants.ContextKeyRequestID, id)
		*r = *r.WithContext(WithValue(r.Context(), id))
		return id, nil
	}
}

// WithValue populates the context with the given request id.
func WithValue(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, constants.ContextKeyRequestID, reqID)
}

// Value extracts the request id from the provided context.
func Value(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(constants.ContextKeyRequestID).(string)
	return id, ok
}

// isTracingEnabled checks if OpenTelemetry tracing is enabled
func isTracingEnabled() bool {
	enabled := os.Getenv("TRACER_ENABLED")
	if enabled == "true" || enabled == "1" {
		return true
	}
	tracerProvider := otel.GetTracerProvider()
	return tracerProvider != nil
}

// getTraceIDFromContext extracts trace ID from OpenTelemetry span context if available
func getTraceIDFromContext(ctx context.Context) string {
	if !isTracingEnabled() {
		return ""
	}
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	if !spanCtx.IsValid() {
		return ""
	}
	traceID := spanCtx.TraceID()
	if traceID.IsValid() {
		return traceID.String()
	}
	return ""
}
