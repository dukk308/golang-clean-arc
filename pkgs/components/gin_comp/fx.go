package gin_comp

import (
	"go.uber.org/fx"
)

var GinComponentFx = fx.Module("gin",
	fx.Provide(NewGinComponent),
)
