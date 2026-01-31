package config

import "go.uber.org/fx"

var ConfigModuleFx = fx.Module("config",
	fx.Provide(LoadConfig),
)
