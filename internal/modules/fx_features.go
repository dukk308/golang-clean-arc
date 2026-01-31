package modules

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth"
	auth_http "github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/presentation/http"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user"
	user_http "github.com/dukk308/golang-clean-arch-starter/internal/modules/user/presentation/http"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gin_comp"
	"go.uber.org/fx"
)

func SetupRoutes(
	lc fx.Lifecycle,
	ginComponent *gin_comp.GinEngine,
	userHTTP *user_http.Http,
	authHTTP *auth_http.Http,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			userHTTP.RegisterRoutes(ginComponent.GetGroup())
			authHTTP.RegisterRoutes(ginComponent.GetGroup())
			return nil
		},
	})
}

var FeatureModuleFx = fx.Module(
	"feature_modules",
	user.Module,
	auth.Module,

	fx.Invoke(SetupRoutes),
)
