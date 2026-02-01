package single_flight_comp

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"singleflight",
	fx.Provide(NewGroup),
)
