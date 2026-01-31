package user

import (
	auth_domain "github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/application"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/infrastructure/persistence"
	user_http "github.com/dukk308/golang-clean-arch-starter/internal/modules/user/presentation/http"
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
		application.NewGetProfileQuery,
	),
	fx.Provide(
		func(
			getProfileQuery *application.GetProfileQuery,
			tokenService auth_domain.ITokenService,
		) *user_http.Http {
			return user_http.NewHttp(getProfileQuery, tokenService)
		},
	),
)
