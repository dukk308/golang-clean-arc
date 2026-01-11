package modules

import (
	"context"

	"github.com/dukk308/golang-clean-arc/internal/modules/user"
	user_http "github.com/dukk308/golang-clean-arc/internal/modules/user/presentation/http"
	"github.com/dukk308/golang-clean-arc/pkgs/components/gin_comp"
	"go.uber.org/fx"
)

func SetupRoutes(
	lc fx.Lifecycle,
	ginComponent *gin_comp.GinComponent,
	userHTTP *user_http.Http,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			userHTTP.RegisterRoutes(ginComponent.Router())
			return nil
		},
	})
}

var FeatureModuleFx = fx.Module(
	"feature_modules",
	user.Module,

	fx.Invoke(SetupRoutes),
)
