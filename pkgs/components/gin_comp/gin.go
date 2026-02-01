package gin_comp

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type GinEngine struct {
	config *GinConfig
	logger logger.Logger
	router *gin.Engine
	group  *gin.RouterGroup
}

func NewGinComp(logger logger.Logger, config *GinConfig) *GinEngine {
	engine := &GinEngine{
		logger: logger,
		config: config,
	}

	if config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = ginLoggerWriter{logger: logger}
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			logger.Debugf("Route registered: %s %s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
		}
	}

	engine.router = gin.New()
	engine.group = engine.router.Group(config.Prefix)

	if config.EnableTracer {
		engine.withInstrumentation()
	}

	return engine
}

func (gs *GinEngine) GetConfig() *GinConfig {
	return gs.config
}

func (gs *GinEngine) GetLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		gs.logger.Infof(fmt.Sprintf("%d | %s | %s | %s | %s | %s",
			param.StatusCode,
			param.ClientIP,
			param.Method,
			param.Path,
			param.Latency,
			param.ErrorMessage,
		))
		return "" // Gin's logger expects empty string to suppress default output
	})
}

func (gs *GinEngine) GetRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger := logger.FromContext(c.Request.Context())
		c.Header("Content-Type", "application/json")

		stack := debug.Stack()
		logger.Errorf("Panic recovered: %+v\nStack trace:\n%s", recovered, string(stack))

		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"code":    "internal_server_error",
					"message": "An internal server error occurred",
				},
			},
		})
	})
}

func (gs *GinEngine) GetRouter() *gin.Engine {
	return gs.router
}

func (gs *GinEngine) GetGroup() *gin.RouterGroup {
	if gs.group == nil {
		gs.group = gs.router.Group(gs.config.Prefix)
	}

	return gs.group
}

// RegisterHealthCheck registers a health check endpoint on the Gin router
func (gs *GinEngine) RegisterHealthCheck(path string) {
	if gs.router == nil {
		gs.logger.Error("Cannot register health check: router not initialized")
		return
	}

	gs.router.GET(path, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": gs.config.ServiceName,
		})
	})

	gs.logger.Infof("Health check endpoint registered at %s", path)
}

// RegisterReadinessCheck registers a readiness check endpoint on the Gin router
func (gs *GinEngine) RegisterReadinessCheck(path string) {
	if gs.router == nil {
		gs.logger.Error("Cannot register readiness check: router not initialized")
		return
	}

	gs.router.GET(path, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ready",
			"service": gs.config.ServiceName,
		})
	})

	gs.logger.Infof("Readiness check endpoint registered at %s", path)
}

func (gs *GinEngine) RegisterLivenessCheck(path string) {
	if gs.router == nil {
		gs.logger.Error("Cannot register liveness check: router not initialized")
		return
	}

	gs.router.GET(path, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "alive",
			"service": gs.config.ServiceName,
		})
	})

	gs.logger.Infof("Liveness check endpoint registered at %s", path)
}

func GetHealthCheckTracerFilter() func(*http.Request) bool {
	return func(r *http.Request) bool {
		path := r.URL.Path
		return path != "/health" && path != "/ready" && path != "/alive"
	}
}

func (gs *GinEngine) withInstrumentation() *gin.Engine {
	if gs.router == nil {
		gs.logger.Error("Cannot setup tracer: router not initialized")
		return gs.router
	}
	gs.router.Use(
		otelgin.Middleware(
			gs.config.ServiceName,
			otelgin.WithFilter(GetHealthCheckTracerFilter()),
		),
	)

	return gs.router
}
