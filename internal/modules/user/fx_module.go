package user

import (
	auth_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user/application"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user/domain"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user/infrastructure/persistence"
	user_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/user/presentation/http"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(
		fx.Annotate(
			persistence.NewViewerRepository,
			fx.As(new(domain.IViewerRepository)),
		),
	),
	fx.Provide(
		application.NewViewerGetProfileQuery,
	),
	fx.Provide(
		func(
			viewerGetProfileQuery *application.ViewerGetProfileQuery,
			tokenService auth_domain.ITokenService,
		) *user_http.Http {
			return user_http.NewHttp(viewerGetProfileQuery, tokenService)
		},
	),
)
