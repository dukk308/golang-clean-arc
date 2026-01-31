package swagger_comp

import (
	"go.uber.org/fx"
)

func ProvideSwaggerConfig() *SwaggerConfig {
	return LoadSwaggerConfig()
}

var SwaggerComponentFx = fx.Module("swagger",
	fx.Provide(ProvideSwaggerConfig),
	fx.Provide(NewSwaggerComponent),
)
