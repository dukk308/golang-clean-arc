package global_config

import "go.uber.org/fx"

func ProvideGlobalConfig() *GlobalConfig {
	return LoadGlobalConfig()
}

var GlobalConfigFx = fx.Options(
	fx.Provide(ProvideGlobalConfig),
)
