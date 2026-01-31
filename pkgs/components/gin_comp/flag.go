package gin_comp

import (
	"flag"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/global_config"
)

var (
	ginPortVal      string
	ginModeVal      string
	ginPrefixVal    string
	enableTracerVal bool
)

var (
	ginPort      = &ginPortVal
	ginMode      = &ginModeVal
	ginPrefix    = &ginPrefixVal
	enableTracer = &enableTracerVal
)

func init() {
	if flag.Lookup("gin-port") == nil {
		flag.StringVar(&ginPortVal, "gin-port", "5005", "gin server port. Default 5005")
	}
	if flag.Lookup("gin-mode") == nil {
		flag.StringVar(&ginModeVal, "gin-mode", "debug", "gin mode (debug | release). Default debug")
	}
	if flag.Lookup("gin-prefix") == nil {
		flag.StringVar(&ginPrefixVal, "gin-prefix", "", "gin prefix")
	}
	if flag.Lookup("enable-tracer") == nil {
		flag.BoolVar(&enableTracerVal, "enable-tracer", false, "enable tracer. Default false")
	}
}

func LoadGinConfig(global_config *global_config.GlobalConfig) *GinConfig {
	return &GinConfig{
		Port:         *ginPort,
		Mode:         *ginMode,
		Prefix:       *ginPrefix,
		EnableTracer: *enableTracer,
		ServiceName:  global_config.ServiceName,
	}
}
