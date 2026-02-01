package middleware

import (
	"net/http"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/global_config"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/utils/request_id"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

var otelExcludedPaths = []string{"/health", "/ping", "/alive"}

func Tracer(globalConfig *global_config.GlobalConfig) gin.HandlerFunc {
	if globalConfig.EnableTracing {
		otelMiddleware := otelgin.Middleware(globalConfig.ServiceName,
			otelgin.WithFilter(func(r *http.Request) bool {
				path := r.URL.Path
				for _, excluded := range otelExcludedPaths {
					if path == excluded {
						return false
					}
				}
				return true
			}))

		return func(c *gin.Context) {
			otelMiddleware(c)

			span := trace.SpanFromContext(c.Request.Context())
			if span.SpanContext().IsValid() {
				traceID := span.SpanContext().TraceID().String()
				c.Header("X-Request-ID", traceID)
				c.Writer.Header().Set(constants.ContextKeyRequestID, traceID)
				c.Request = c.Request.WithContext(request_id.WithValue(c.Request.Context(), traceID))
			}
		}
	}

	provider := request_id.RequestID(func() (string, error) {
		return uuid.New().String(), nil
	})

	return func(c *gin.Context) {
		requestID, _ := provider(c.Writer, c.Request)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
