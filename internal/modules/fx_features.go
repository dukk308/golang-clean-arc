package modules

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth"
	auth_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/presentation/http"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog"
	blog_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/presentation/http"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note"
	note_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/note/presentation/http"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user"
	user_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/user/presentation/http"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"go.uber.org/fx"
)

func SetupRoutes(
	lc fx.Lifecycle,
	ginComponent *gin_comp.GinEngine,
	userHTTP *user_http.Http,
	authHTTP *auth_http.Http,
	noteHTTP *note_http.Http,
	blogHTTP *blog_http.Http,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			userHTTP.RegisterRoutes(ginComponent.GetGroup())
			authHTTP.RegisterRoutes(ginComponent.GetGroup())
			noteHTTP.RegisterRoutes(ginComponent.GetGroup())
			blogHTTP.RegisterRoutes(ginComponent.GetGroup())
			return nil
		},
	})
}

var FeatureModuleFx = fx.Module(
	"feature_modules",
	user.Module,
	auth.Module,
	note.Module,
	blog.Module,
	fx.Invoke(SetupRoutes),
)
