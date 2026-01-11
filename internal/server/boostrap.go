package server

import (
	"context"
	"net/http"

	"github.com/dukk308/golang-clean-arc/internal/config"
	"github.com/dukk308/golang-clean-arc/internal/modules"
	"github.com/dukk308/golang-clean-arc/pkgs/components/gin_comp"
	"github.com/dukk308/golang-clean-arc/pkgs/components/gorm_comp"
	"go.uber.org/fx"
)

func startHttpServer(lc fx.Lifecycle, ginComponent *gin_comp.GinComponent, config *config.Config) {
	httpServer := &http.Server{
		Addr:    config.Server.Port,
		Handler: ginComponent.Engine(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return httpServer.ListenAndServe()
		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})
}

func Bootstrap() *fx.App {
	app := fx.New(
		fx.Options(
			config.ConfigModuleFx,
			gorm_comp.GormComponentFx,
			gin_comp.GinComponentFx,
			modules.FeatureModuleFx,

			fx.Invoke(startHttpServer),
		),
	)

	return app
}
