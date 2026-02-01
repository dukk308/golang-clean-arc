package request_id

import (
	"context"
	"net/http"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
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
		var id = r.Header.Get(constants.ContextKeyRequestID)

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
