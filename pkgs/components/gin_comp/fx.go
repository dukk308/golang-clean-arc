package gin_comp

import (
	"github.com/dukk308/golang-clean-arch-starter/pkgs/global_config"
	"go.uber.org/fx"
)

func ProvideGinConfig(global_config *global_config.GlobalConfig) *GinConfig {
	return LoadGinConfig(global_config)
}

var GinComponentFx = fx.Module("gin",
	fx.Provide(ProvideGinConfig),
	fx.Provide(NewGinComp),
)
