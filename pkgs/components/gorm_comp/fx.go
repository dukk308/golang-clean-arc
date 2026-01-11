package gorm_comp

import "go.uber.org/fx"

var GormComponentFx = fx.Module("gorm",
	fx.Provide(NewGormComponent),
)
