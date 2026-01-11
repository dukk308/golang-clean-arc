package user

import (
	"github.com/dukk308/golang-clean-arc/internal/modules/user/application"
	"github.com/dukk308/golang-clean-arc/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arc/internal/modules/user/infrastructure/persistence"
	"github.com/dukk308/golang-clean-arc/internal/modules/user/presentation/http"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(
		application.NewCreateNewViewerCommand,
	),
	fx.Provide(
		fx.Annotate(
			persistence.NewViewerRepository,
			fx.As(new(domain.IViewerRepository)),
		),
	),
	fx.Provide(
		http.NewHttp,
	),
)
