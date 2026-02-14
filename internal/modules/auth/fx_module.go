package auth

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/config"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/application"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/infrastructure"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/infrastructure/repository"
	auth_http "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/presentation/http"
	user_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/user/domain"
	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		func(cfg *config.Config) domain.ITokenService {
			return domain.NewTokenService(cfg.Auth.AccessTokenSecret, cfg.Auth.RefreshTokenSecret)
		},
	),
	fx.Provide(
		fx.Annotate(
			domain.NewInMemoryTokenStorage,
			fx.As(new(domain.ITokenStorage)),
		),
	),
	fx.Provide(
		func(userRepository user_domain.IViewerRepository) domain.IUserRepository {
			return repository.NewUserRepositoryAdapter(userRepository)
		},
	),
	fx.Provide(
		func(cfg *config.Config) domain.IGoogleOAuthService {
			return infrastructure.NewGoogleOAuthService(&cfg.Auth)
		},
	),
	fx.Provide(
		application.NewSignupCommand,
		application.NewSigninCommand,
		application.NewSignoutCommand,
		application.NewRefreshTokenCommand,
		application.NewGoogleSigninCommand,
	),
	fx.Provide(
		auth_http.NewGoogleOAuthHttp,
		auth_http.NewHttp,
	),
)
