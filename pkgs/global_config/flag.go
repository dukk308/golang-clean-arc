package global_config

import (
	"flag"
)

var (
	serviceNameVal string
	environmentVal string
	logLevelVal string
	callerEnabledVal bool
	enableTracingVal bool
)

var (
	ServiceName = &serviceNameVal
	Environment = &environmentVal
	LogLevel = &logLevelVal
	CallerEnabled = &callerEnabledVal
	EnableTracing = &enableTracingVal
)

func init() {
	if flag.Lookup("service-name") == nil {
		flag.StringVar(&serviceNameVal, "service-name", "golang-clean-arc", "Service name")
	}
	if flag.Lookup("app-env") == nil {
		flag.StringVar(&environmentVal, "app-env", "local", "Application environment (local, dev, prod)")
	}
	if flag.Lookup("log-level") == nil {
		flag.StringVar(&logLevelVal, "log-level", "debug", "Log level (debug, info, warn, error, fatal)")
	}
	if flag.Lookup("caller-enabled") == nil {
		flag.BoolVar(&callerEnabledVal, "caller-enabled", true, "Caller enabled")
	}
	if flag.Lookup("enable-tracing") == nil {
		flag.BoolVar(&enableTracingVal, "enable-tracing", true, "Enable tracing")
	}
}

func LoadGlobalConfig() *GlobalConfig {
	result := &GlobalConfig{
		ServiceName: *ServiceName,
		Environment: *Environment,
		LogLevel: *LogLevel,
		CallerEnabled: *CallerEnabled,
		EnableTracing: *EnableTracing,
	}

	return result
}
