package config

import "github.com/dukk308/golang-clean-arch-starter/pkgs/global_config"

type LogOptions struct {
	LogLevel      string
	CallerEnabled bool
	EnableTracing bool
}

func ProvideLogConfig(globCfg *global_config.GlobalConfig) *LogOptions {
	return &LogOptions{
		LogLevel:      globCfg.LogLevel,
		CallerEnabled: globCfg.CallerEnabled,
		EnableTracing: globCfg.EnableTracing,
	}
}
