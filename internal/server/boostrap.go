package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dukk308/golang-clean-arch-starter/internal/config"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gin_comp"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gorm_comp"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/swagger_comp"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/global_config"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/logger"
	"go.uber.org/fx"
)

func startHttpServer(
	lc fx.Lifecycle,
	ginComponent *gin_comp.GinEngine,
	swaggerComponent *swagger_comp.SwaggerComponent,
	config *config.Config,
	log logger.Logger,
) {
	router := ginComponent.GetRouter()
	swaggerComponent.RegisterRoutes(router)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", ginComponent.GetConfig().Port),
		Handler: router,
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
