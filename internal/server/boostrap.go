package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dukk308/beetool.dev-go-starter/internal/config"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules"
	"github.com/dukk308/beetool.dev-go-starter/internal/validation"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gorm_comp"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/swagger_comp"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/global_config"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	middleware "github.com/dukk308/beetool.dev-go-starter/pkgs/middlewares/gin"
	"go.uber.org/fx"
)

func registerValidation(log logger.Logger) {
	if err := validation.RegisterValidations(); err != nil {
		log.Errorf("Failed to register validations: %v", err)
	}
}

func startHttpServer(
	lc fx.Lifecycle,
	ginComponent *gin_comp.GinEngine,
	swaggerComponent *swagger_comp.SwaggerComponent,
	config *config.Config,
	globalConfig *global_config.GlobalConfig,
	log logger.Logger,
) {
	group := ginComponent.GetGroup()
	group.Use(middleware.CORS())
	group.Use(middleware.Tracer(globalConfig))
	group.Use(middleware.CorrelateLogger(log))
	group.Use(middleware.Logger(
		globalConfig.IsLogRequest,
		globalConfig.IsLogResponse,
	))
	swaggerComponent.RegisterRoutes(ginComponent.GetRouter())

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", ginComponent.GetConfig().Port),
		Handler: ginComponent.GetRouter(),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof("HTTP server starting on %s", httpServer.Addr)
			go func() {
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Errorf("HTTP server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})
}

func Bootstrap(ctx context.Context) *fx.App {
	app := fx.New(
		global_config.GlobalConfigFx,
		logger.ZapModuleFx,
		config.ConfigModuleFx,
		fx.WithLogger(logger.ProvideFXEventLogger),
		fx.Invoke(registerValidation),
		fx.Options(
			gorm_comp.GormComponentFx,
			gin_comp.GinComponentFx,
			swagger_comp.SwaggerComponentFx,
			modules.FeatureModuleFx,
			fx.Invoke(startHttpServer),
		),
	)

	return app
}
