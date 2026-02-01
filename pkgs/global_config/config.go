package global_config

type GlobalConfig struct {
	ServiceName   string `mapstructure:"service_name"`
	Environment   string `mapstructure:"app_env"`
	LogLevel      string `mapstructure:"log_level"`
	CallerEnabled bool   `mapstructure:"caller_enabled"`
	EnableTracing bool   `mapstructure:"enable_tracing"`
	IsLogRequest  bool   `mapstructure:"is_log_request"`
	IsLogResponse bool   `mapstructure:"is_log_response"`
}
